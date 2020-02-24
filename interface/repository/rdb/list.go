package rdb

import (
	"github.com/jinzhu/gorm"
	"github.com/x-color/vue-trello/model"
)

// ListDBManager is DB manager for List.
type ListDBManager struct {
	db *gorm.DB
}

// Create registers a List to DB.
func (m *ListDBManager) Create(list model.List) error {
	if err := m.db.Create(&list).Error; err != nil {
		return model.ServerError{
			Err: err,
			ID:  list.ID,
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

	err := m.db.Model(&list).Where(&model.List{ID: list.ID, UserID: list.UserID}).Updates(map[string]interface{}{
		"title": list.Title,
	}).Error

	if err != nil {
		if gorm.IsRecordNotFoundError(err) {
			return model.NotFoundError{
				Err: err,
				ID:  list.ID,
				Act: "update list",
			}
		}
		return model.ServerError{
			Err: err,
			ID:  list.ID,
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

	tx := m.db.Begin()

	if err := tx.Where(&model.List{ID: list.ID, UserID: list.UserID}).Delete(&list).Error; err != nil {
		tx.Rollback()
		if gorm.IsRecordNotFoundError(err) {
			return model.NotFoundError{
				Err: err,
				ID:  list.ID,
				Act: "delete list",
			}
		}
		return model.ServerError{
			Err: err,
			ID:  list.ID,
			Act: "delete list",
		}
	}

	if err := tx.Where(&model.Item{ListID: list.ID}).Delete(model.List{}).Error; err != nil {
		tx.Rollback()
		return model.ServerError{
			Err: err,
			ID:  list.ID,
			Act: "delete list",
		}
	}

	tx.Commit()
	return nil
}

// Find gets a List had specific ID from DB.
func (m *ListDBManager) Find(list model.List) (model.List, error) {
	r := model.List{}
	if err := m.db.Where(&model.List{ID: list.ID, UserID: list.UserID}).First(&r).Error; err != nil {
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
	return r, nil
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
	r := model.Lists{}
	if err := m.db.Where(&model.List{BoardID: board.ID, UserID: board.UserID}).Find(r).Error; err != nil {
		return model.Lists{}, model.ServerError{
			Err: err,
			ID:  board.ID,
			Act: "find lists in board",
		}
	}
	return r, nil
}
