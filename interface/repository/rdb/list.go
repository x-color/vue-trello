package rdb

import (
	"time"

	"github.com/jinzhu/gorm"
	"github.com/x-color/vue-trello/model"
)

// List is List data model for DB.
type List struct {
	ID        string `gorm:"primary_key"`
	UserID    string `gorm:"primary_key"`
	BoardID   string
	Title     string
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time
}

func (l *List) convertFrom(list model.List) {
	l.ID = list.ID
	l.UserID = list.UserID
	l.BoardID = list.BoardID
	l.Title = list.Title
}

func (l *List) convertTo() model.List {
	list := model.List{
		ID:      l.ID,
		UserID:  l.UserID,
		BoardID: l.BoardID,
		Title:   l.Title,
	}
	return list
}

// Lists is a slice of List data model.
type Lists []List

// ListDBManager is DB manager for List.
type ListDBManager struct {
	db *gorm.DB
}

// Create registers a List to DB.
func (m *ListDBManager) Create(list model.List) error {
	l := List{}
	l.convertFrom(list)

	if err := m.db.Create(&l).Error; err != nil {
		return model.ServerError{
			Err: err,
			ID:  l.ID,
			Act: "create list",
		}
	}
	return nil
}

// Update updates all fields of specific List in DB.
func (m *ListDBManager) Update(list model.List) error {
	if validatePrimaryKey(list.ID) {
		return model.NotFoundError{
			Err: nil,
			ID:  "(No ID)",
			Act: "validate list",
		}
	}

	l := List{}
	l.convertFrom(list)
	err := m.db.Model(&l).Updates(map[string]interface{}{
		"title": list.Title,
	}).Error

	if err != nil {
		if gorm.IsRecordNotFoundError(err) {
			return model.NotFoundError{
				Err: err,
				ID:  l.ID,
				Act: "update list",
			}
		}
		return model.ServerError{
			Err: err,
			ID:  l.ID,
			Act: "update list",
		}
	}
	return nil
}

// Delete removes a List from DB.
func (m *ListDBManager) Delete(list model.List) error {
	if validatePrimaryKey(list.ID) {
		return model.NotFoundError{
			Err: nil,
			ID:  "(No ID)",
			Act: "validate list",
		}
	}

	l := List{}
	l.convertFrom(list)

	tx := m.db.Begin()

	if err := tx.Delete(&l).Error; err != nil {
		tx.Rollback()
		if gorm.IsRecordNotFoundError(err) {
			return model.NotFoundError{
				Err: err,
				ID:  l.ID,
				Act: "delete list",
			}
		}
		return model.ServerError{
			Err: err,
			ID:  l.ID,
			Act: "delete list",
		}
	}

	if err := tx.Where(&Item{ListID: l.ID}).Delete(List{}).Error; err != nil {
		tx.Rollback()
		return model.ServerError{
			Err: err,
			ID:  l.ID,
			Act: "delete list",
		}
	}

	tx.Commit()
	return nil
}

// Find gets a List had specific ID from DB.
func (m *ListDBManager) Find(list model.List) (model.List, error) {
	r := List{}
	if err := m.db.Where(&List{ID: list.ID, UserID: list.UserID}).First(&r).Error; err != nil {
		if gorm.IsRecordNotFoundError(err) {
			return model.List{}, model.NotFoundError{
				Err: err,
				ID:  list.ID,
				Act: "find list",
			}
		}
		return model.List{}, model.ServerError{
			Err: err,
			ID:  list.ID,
			Act: "find list",
		}
	}
	return r.convertTo(), nil
}

// FindLists gets all Lists in a specific board from DB.
func (m *ListDBManager) FindLists(board model.Board) (model.Lists, error) {
	if validatePrimaryKey(board.ID) {
		return model.Lists{}, model.NotFoundError{
			Err: nil,
			ID:  "(No ID)",
			Act: "validate list",
		}
	}
	r := Lists{}
	if err := m.db.Where(&List{BoardID: board.ID, UserID: board.UserID}).Find(&r).Error; err != nil {
		return model.Lists{}, model.ServerError{
			Err: err,
			ID:  board.ID,
			Act: "find lists in board",
		}
	}

	lists := model.Lists{}
	for _, rl := range r {
		lists = append(lists, rl.convertTo())
	}

	return lists, nil
}
