package rdb

import (
	"github.com/jinzhu/gorm"
	"github.com/x-color/vue-trello/model"
)

// ItemDBManager is DB manager for Item.
type ItemDBManager struct {
	db *gorm.DB
}

// Create registers a Item to DB.
func (m *ItemDBManager) Create(item model.Item) error {
	if err := m.db.Create(&item).Error; err != nil {
		return model.ServerError{
			Err: err,
			ID:  item.ID,
			Act: "create item",
		}
	}
	return nil
}

// Update updates all fields of specific Item in DB.
func (m *ItemDBManager) Update(item model.Item) error {
	if validatePrimaryKey(item.ID) {
		return model.NotFoundError{
			Err: nil,
			ID:  "(No ID)",
			Act: "validate item",
		}
	}
	if err := m.db.Save(&item).Error; err != nil {
		if gorm.IsRecordNotFoundError(err) {
			return model.NotFoundError{
				Err: err,
				ID:  item.ID,
				Act: "update item",
			}
		}
		return model.ServerError{
			Err: err,
			ID:  item.ID,
			Act: "update item",
		}
	}
	return nil
}

// Delete removes a Item from DB.
func (m *ItemDBManager) Delete(item model.Item) error {
	if validatePrimaryKey(item.ID) {
		return model.NotFoundError{
			Err: nil,
			ID:  "(No ID)",
			Act: "validate item",
		}
	}
	if err := m.db.Delete(&model.Item{ID: item.ID}).Error; err != nil {
		if gorm.IsRecordNotFoundError(err) {
			return model.NotFoundError{
				Err: err,
				ID:  item.ID,
				Act: "delete item",
			}
		}
		return model.ServerError{
			Err: err,
			ID:  item.ID,
			Act: "delete item",
		}
	}
	return nil
}

// Find gets a Item had specific ID from DB.
func (m *ItemDBManager) Find(item model.Item) (model.Item, error) {
	r := model.Item{}
	if err := m.db.Where(&model.Item{ID: item.ID}).First(&r).Error; err != nil {
		if gorm.IsRecordNotFoundError(err) {
			return model.Item{}, model.NotFoundError{
				Err: err,
				ID:  item.ID,
				Act: "find item",
			}
		}
		return model.Item{}, model.ServerError{
			Err: err,
			ID:  item.ID,
			Act: "find item",
		}
	}
	return r, nil
}

// FindItems gets all Item in a specific List from DB.
func (m *ItemDBManager) FindItems(list model.List) (model.Items, error) {
	if validatePrimaryKey(list.ID) {
		return model.Items{}, model.NotFoundError{
			Err: nil,
			ID:  "(No ID)",
			Act: "validate list",
		}
	}
	r := model.Items{}
	if err := m.db.Where(&model.Item{ListID: list.ID}).Find(r).Error; err != nil {
		return model.Items{}, model.ServerError{
			Err: err,
			ID:  list.ID,
			Act: "find items in list",
		}
	}
	return r, nil
}

func validatePrimaryKey(key string) bool {
	return key != ""
}
