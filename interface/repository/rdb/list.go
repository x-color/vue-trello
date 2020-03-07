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
	Before    *string
	After     *string
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time
}

func (l *List) convertFrom(list model.List) {
	l.ID = list.ID
	l.UserID = list.UserID
	l.BoardID = list.BoardID
	l.Title = list.Title

	if list.Before == "" {
		l.Before = nil
	} else {
		l.Before = &list.Before
	}

	if list.After == "" {
		l.After = nil
	} else {
		l.After = &list.After
	}
}

func (l *List) convertTo() model.List {
	list := model.List{
		ID:      l.ID,
		UserID:  l.UserID,
		BoardID: l.BoardID,
		Title:   l.Title,
	}

	if l.After == nil {
		list.After = ""
	} else {
		list.After = *l.After
	}

	if l.Before == nil {
		list.Before = ""
	} else {
		list.Before = *l.Before
	}

	return list
}

// Lists is a slice of List data model.
type Lists []List

// ListDBManager is DB manager for List.
type ListDBManager struct {
	db *gorm.DB
}

func newListDBManager(db *gorm.DB) ListDBManager {
	db.AutoMigrate(&List{})
	return ListDBManager{
		db: db,
	}
}

// Create registers a List to DB.
func (m *ListDBManager) Create(list model.List) error {
	if err := validatePrimaryKeys("list", list.ID, list.UserID); err != nil {
		return err
	}

	tx := m.db.Begin()

	beforeList := new(List)

	if err := tx.Where(map[string]interface{}{"board_id": list.BoardID, "user_id": list.UserID, "after": nil}).First(beforeList).Error; err != nil {
		if gorm.IsRecordNotFoundError(err) {
			beforeList = nil
		} else {
			tx.Rollback()
			return model.ServerError{
				UserID: list.UserID,
				Err:    err,
				ID:     list.ID,
				Act:    "create list",
			}
		}
	}

	l := List{}
	l.convertFrom(list)
	l.After = nil
	if beforeList == nil {
		l.Before = nil
	} else {
		l.Before = &beforeList.ID

		err := tx.Model(beforeList).Updates(map[string]interface{}{
			"after": l.ID,
		}).Error
		if err != nil {
			tx.Rollback()
			return convertError(err, beforeList.ID, beforeList.UserID, "update List to create new list")
		}
	}

	if err := tx.Create(&l).Error; err != nil {
		tx.Rollback()
		return model.ServerError{
			UserID: l.UserID,
			Err:    err,
			ID:     l.ID,
			Act:    "create list",
		}
	}
	tx.Commit()
	return nil
}

// Update updates all fields of specific List in DB.
func (m *ListDBManager) Update(list model.List) error {
	if err := validatePrimaryKeys("list", list.ID, list.UserID); err != nil {
		return err
	}

	l := List{}
	l.convertFrom(list)
	err := m.db.Model(&l).Updates(map[string]interface{}{
		"title": list.Title,
	}).Error

	if err != nil {
		return convertError(err, l.ID, l.UserID, "update list")
	}
	return nil
}

// Move updates list's position data..
func (m *ListDBManager) Move(list model.List) error {
	if err := validatePrimaryKeys("list", list.ID, list.UserID); err != nil {
		return err
	}

	l := List{}
	l.convertFrom(list)

	tx := m.db.Begin()

	oldList := new(List)
	err := tx.Where(&List{ID: l.ID, UserID: l.UserID}).First(oldList).Error
	if err != nil {
		tx.Rollback()
		return convertError(err, l.ID, l.UserID, "find moved list")
	}

	// Update a list before moved list's old position
	if oldList.Before != nil {
		oldBeforeList := List{
			ID:     *oldList.Before,
			UserID: oldList.UserID,
			After:  oldList.After,
		}

		err = tx.Model(&oldBeforeList).Updates(map[string]interface{}{
			"after": convertData(oldBeforeList.After),
		}).Error

		if err != nil {
			tx.Rollback()
			return convertError(err, oldBeforeList.ID, oldBeforeList.UserID, "update list before moved list's old position")
		}
	}

	// Update a list after moved list's old position
	if oldList.After != nil {
		oldAfterList := List{
			ID:     *oldList.After,
			UserID: oldList.UserID,
			Before: oldList.Before,
		}

		err = tx.Model(&oldAfterList).Updates(map[string]interface{}{
			"before": convertData(oldAfterList.Before),
		}).Error

		if err != nil {
			tx.Rollback()
			return convertError(err, oldAfterList.ID, oldAfterList.UserID, "update list after moved list's old position")
		}
	}

	if l.Before == nil {
		newAfterList := new(List)
		err := tx.Where(&List{UserID: l.UserID, Before: nil}).First(newAfterList).Error
		if err != nil {
			tx.Rollback()
			return convertError(err, l.ID, l.UserID, "find list after moved list")
		}
		l.After = &newAfterList.ID
	} else {
		newBeforeList := new(List)
		err := tx.Where(&List{ID: *l.Before, UserID: l.UserID}).First(newBeforeList).Error
		if err != nil {
			tx.Rollback()
			return convertError(err, l.ID, l.UserID, "find list before moved list")
		}
		l.After = newBeforeList.After
	}

	// Update a list before moved list's new position
	if l.Before != nil {
		newBeforeList := List{
			ID:     *l.Before,
			UserID: l.UserID,
			After:  &l.ID,
		}

		err = tx.Model(&newBeforeList).Updates(map[string]interface{}{
			"after": convertData(newBeforeList.After),
		}).Error

		if err != nil {
			tx.Rollback()
			return convertError(err, newBeforeList.ID, newBeforeList.UserID, "update list before moved list's new position")
		}
	}

	// Update a list after moved list's new position
	if l.After != nil {
		newAfterList := List{
			ID:     *l.After,
			UserID: l.UserID,
			Before: &l.ID,
		}

		err = tx.Model(&newAfterList).Updates(map[string]interface{}{
			"before": convertData(newAfterList.Before),
		}).Error

		if err != nil {
			tx.Rollback()
			return convertError(err, newAfterList.ID, newAfterList.UserID, "update list after moved list's new position")
		}
	}

	err = tx.Model(&l).Updates(map[string]interface{}{
		"boardID": l.BoardID,
		"after":   convertData(l.After),
		"before":  convertData(l.Before),
	}).Error

	if err != nil {
		tx.Rollback()
		return convertError(err, l.ID, l.UserID, "move list")
	}
	tx.Commit()
	return nil
}

// Delete removes a List from DB.
func (m *ListDBManager) Delete(list model.List) error {
	if err := validatePrimaryKeys("list", list.ID, list.UserID); err != nil {
		return err
	}

	tx := m.db.Begin()

	l := List{}
	l.convertFrom(list)

	deletedList := new(List)
	if err := tx.Where(&l).First(deletedList).Error; err != nil {
		tx.Rollback()
		return convertError(err, l.ID, l.UserID, "find deleting list")
	}

	if err := tx.Delete(&l).Error; err != nil {
		tx.Rollback()
		return convertError(err, l.ID, l.UserID, "delete list")
	}

	// Update list before deleted list
	if deletedList.Before != nil {
		err := tx.Model(&List{ID: *deletedList.Before, UserID: deletedList.UserID}).Updates(map[string]interface{}{
			"after": convertData(deletedList.After),
		}).Error
		if err != nil {
			tx.Rollback()
			return model.ServerError{
				UserID: deletedList.UserID,
				Err:    err,
				ID:     *deletedList.Before,
				Act:    "update list before deleted list",
			}
		}
	}

	// Update list after deleted list
	if deletedList.After != nil {
		err := tx.Model(&List{ID: *deletedList.After, UserID: deletedList.UserID}).Updates(map[string]interface{}{
			"before": convertData(deletedList.Before),
		}).Error
		if err != nil {
			tx.Rollback()
			return model.ServerError{
				UserID: deletedList.UserID,
				Err:    err,
				ID:     *deletedList.After,
				Act:    "update list after deleted list",
			}
		}
	}

	if err := tx.Where(&Item{ListID: l.ID, UserID: l.UserID}).Delete(Item{}).Error; err != nil {
		tx.Rollback()
		return model.ServerError{
			UserID: l.UserID,
			Err:    err,
			ID:     l.ID,
			Act:    "delete items in deleted list",
		}
	}

	tx.Commit()
	return nil
}

// Find gets a List had specific ID from DB.
func (m *ListDBManager) Find(list model.List) (model.List, error) {
	if err := validatePrimaryKeys("list", list.ID, list.UserID); err != nil {
		return model.List{}, err
	}

	r := List{}
	if err := m.db.Where(&List{ID: list.ID, UserID: list.UserID}).First(&r).Error; err != nil {
		return model.List{}, convertError(err, list.ID, list.UserID, "find list")
	}
	return r.convertTo(), nil
}

// FindLists gets all Lists in a specific board from DB.
func (m *ListDBManager) FindLists(board model.Board) (model.Lists, error) {
	if err := validatePrimaryKeys("board", board.ID, board.UserID); err != nil {
		return model.Lists{}, err
	}

	r := Lists{}
	if err := m.db.Where(&List{BoardID: board.ID, UserID: board.UserID}).Find(&r).Error; err != nil {
		return model.Lists{}, model.ServerError{
			UserID: board.UserID,
			Err:    err,
			ID:     board.ID,
			Act:    "find lists in board",
		}
	}

	lists := model.Lists{}
	for _, rl := range r {
		lists = append(lists, rl.convertTo())
	}

	return lists, nil
}
