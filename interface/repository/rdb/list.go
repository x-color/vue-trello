package rdb

import (
	"time"

	"github.com/jinzhu/gorm"
	"github.com/x-color/vue-trello/model"
	"github.com/x-color/vue-trello/usecase"
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
type ListDBManager struct{}

func newListDBManager(db *gorm.DB) ListDBManager {
	db.AutoMigrate(&List{})
	return ListDBManager{}
}

// Create registers a List to DB.
func (*ListDBManager) Create(tx usecase.Transaction, list model.List) error {
	if err := validatePrimaryKeys("list", list.ID, list.UserID); err != nil {
		return err
	}

	l := List{}
	l.convertFrom(list)

	if err := tx.DB().(*gorm.DB).Create(&l).Error; err != nil {
		tx.Rollback()
		return model.ServerError{
			UserID: l.UserID,
			Err:    err,
			ID:     l.ID,
			Act:    "create list",
		}
	}

	return nil
}

// Update updates all fields of specific List in DB.
func (*ListDBManager) Update(tx usecase.Transaction, list model.List, updates map[string]interface{}) error {
	if err := validatePrimaryKeys("list", list.ID, list.UserID); err != nil {
		return err
	}

	l := List{}
	l.convertFrom(list)
	err := tx.DB().(*gorm.DB).Model(&l).Updates(queryForList(updates)).Error
	if err != nil {
		return convertError(err, l.ID, l.UserID, "update list")
	}
	return nil
}

// Delete removes a List from DB.
func (*ListDBManager) Delete(tx usecase.Transaction, list model.List) error {
	if err := validatePrimaryKeys("list", list.ID, list.UserID); err != nil {
		return err
	}

	l := List{}
	l.convertFrom(list)

	if err := tx.DB().(*gorm.DB).Delete(&l).Error; err != nil {
		return convertError(err, l.ID, l.UserID, "delete list")
	}

	return nil
}

// FindByID gets a List had specific ID from DB.
func (*ListDBManager) FindByID(tx usecase.Transaction, id, userID string) (model.List, error) {
	if err := validatePrimaryKeys("list", id, userID); err != nil {
		return model.List{}, err
	}

	r := List{}
	if err := tx.DB().(*gorm.DB).Where(&List{ID: id, UserID: userID}).First(&r).Error; err != nil {
		return model.List{}, convertError(err, id, userID, "find list")
	}
	return r.convertTo(), nil
}

// Find gets Lists.
func (*ListDBManager) Find(tx usecase.Transaction, conditions map[string]interface{}) (model.Lists, error) {
	r := Lists{}
	if err := tx.DB().(*gorm.DB).Where(queryForList(conditions)).Find(&r).Error; err != nil {
		userID := "(No-ID)"
		if v, ok := conditions["user_id"]; ok {
			userID = v.(string)
		}
		id := "(No-ID)"
		if v, ok := conditions["id"]; ok {
			id = v.(string)
		}
		return model.Lists{}, model.ServerError{
			UserID: userID,
			Err:    err,
			ID:     id,
			Act:    "find lists",
		}
	}

	lists := model.Lists{}
	for _, rl := range r {
		lists = append(lists, rl.convertTo())
	}

	return lists, nil
}

func queryForList(data map[string]interface{}) map[string]interface{} {
	query := make(map[string]interface{})
	if v, ok := data["ID"]; ok {
		query["id"] = v
	}
	if v, ok := data["BoardID"]; ok {
		query["board_id"] = v
	}
	if v, ok := data["UserID"]; ok {
		query["user_id"] = v
	}
	if v, ok := data["Title"]; ok {
		query["title"] = v
	}
	if v, ok := data["Before"]; ok {
		if v.(string) == "" {
			query["before"] = nil
		} else {
			query["before"] = v
		}
	}
	if v, ok := data["After"]; ok {
		if v.(string) == "" {
			query["after"] = nil
		} else {
			query["after"] = v
		}
	}
	return query
}
