package rdb

import (
	"strings"
	"time"

	"github.com/jinzhu/gorm"
	"github.com/x-color/vue-trello/model"
)

// Item is Item data model for DB.
type Item struct {
	ID        string `gorm:"primary_key"`
	UserID    string `gorm:"primary_key"`
	ListID    string
	Title     string
	Text      *string
	Tags      *string
	Before    *string
	After     *string
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time
}

func (i *Item) convertFrom(item model.Item) {
	i.ID = item.ID
	i.UserID = item.UserID
	i.ListID = item.ListID
	i.Title = item.Title

	tags := []string{}
	for _, t := range item.Tags {
		tags = append(tags, t.ID)
	}
	ts := strings.Join(tags, ",")
	if ts == "" {
		i.Tags = nil
	} else {
		i.Tags = &ts
	}

	if item.Text == "" {
		i.Text = nil
	} else {
		i.Text = &item.Text
	}

	if item.Before == "" {
		i.Before = nil
	} else {
		i.Before = &item.Before
	}

	if item.After == "" {
		i.After = nil
	} else {
		i.After = &item.After
	}
}

func (i *Item) convertTo() model.Item {
	item := model.Item{
		ID:     i.ID,
		UserID: i.UserID,
		ListID: i.ListID,
		Title:  i.Title,
		Tags:   model.Tags{},
	}

	if i.Text == nil {
		item.Text = ""
	} else {
		item.Text = *i.Text
	}

	if i.Tags != nil {
		for _, tagID := range strings.Split(*i.Tags, ",") {
			item.Tags = append(item.Tags, model.Tag{ID: tagID})
		}
	}

	if i.After == nil {
		item.After = ""
	} else {
		item.After = *i.After
	}

	if i.Before == nil {
		item.Before = ""
	} else {
		item.Before = *i.Before
	}

	return item
}

// Items is a slice of Item data model.
type Items []Item

// ItemDBManager is DB manager for Item.
type ItemDBManager struct {
	db *gorm.DB
}

func newItemDBManager(db *gorm.DB) ItemDBManager {
	db.AutoMigrate(&Item{})
	return ItemDBManager{
		db: db,
	}
}

// Create registers a Item to DB.
func (m *ItemDBManager) Create(item model.Item) error {
	if err := validatePrimaryKeys("item", item.ID, item.UserID); err != nil {
		return err
	}

	tx := m.db.Begin()

	beforeItem := new(Item)

	if err := tx.Where(map[string]interface{}{"list_id": item.ListID, "after": nil}).First(beforeItem).Error; err != nil {
		if gorm.IsRecordNotFoundError(err) {
			beforeItem = nil
		} else {
			tx.Rollback()
			return model.ServerError{
				UserID: item.UserID,
				Err:    err,
				ID:     item.ID,
				Act:    "create item",
			}
		}
	}

	i := Item{}
	i.convertFrom(item)
	i.After = nil
	if beforeItem == nil {
		i.Before = nil
	} else {
		i.Before = &beforeItem.ID

		err := tx.Model(beforeItem).Updates(map[string]interface{}{
			"after": i.ID,
		}).Error
		if err != nil {
			tx.Rollback()
			return convertError(err, beforeItem.ID, beforeItem.UserID, "update item to create new item")
		}
	}

	if err := tx.Create(&i).Error; err != nil {
		tx.Rollback()
		return model.ServerError{
			UserID: i.UserID,
			Err:    err,
			ID:     i.ID,
			Act:    "create item",
		}
	}

	tx.Commit()
	return nil
}

// Update updates all fields of specific Item in DB.
func (m *ItemDBManager) Update(item model.Item) error {
	if err := validatePrimaryKeys("item", item.ID, item.UserID); err != nil {
		return err
	}

	i := Item{}
	i.convertFrom(item)

	err := m.db.Model(&i).Updates(map[string]interface{}{
		"title": i.Title,
		"text":  convertData(i.Text),
		"tags":  convertData(i.Tags),
	}).Error

	if err != nil {
		return convertError(err, i.ID, i.UserID, "update item")
	}
	return nil
}

// Move updates item's position data..
func (m *ItemDBManager) Move(item model.Item) error {
	if err := validatePrimaryKeys("item", item.ID, item.UserID); err != nil {
		return err
	}

	i := Item{}
	i.convertFrom(item)

	tx := m.db.Begin()

	oldItem := new(Item)
	err := tx.Where(&Item{ID: i.ID}).First(oldItem).Error
	if err != nil {
		tx.Rollback()
		return convertError(err, i.ID, i.UserID, "find moved item")
	}

	// Update a item before moved item's old position
	if oldItem.Before != nil {
		oldBeforeItem := Item{
			ID:     *oldItem.Before,
			UserID: oldItem.UserID,
			After:  oldItem.After,
		}

		err = tx.Model(&oldBeforeItem).Updates(map[string]interface{}{
			"after": convertData(oldBeforeItem.After),
		}).Error

		if err != nil {
			tx.Rollback()
			return convertError(err, oldBeforeItem.ID, oldBeforeItem.UserID, "update item before moved item's old position")
		}
	}

	// Update a item after moved item's old position
	if oldItem.After != nil {
		oldAfterItem := Item{
			ID:     *oldItem.After,
			UserID: oldItem.UserID,
			Before: oldItem.Before,
		}

		err = tx.Model(&oldAfterItem).Updates(map[string]interface{}{
			"before": convertData(oldAfterItem.Before),
		}).Error

		if err != nil {
			tx.Rollback()
			return convertError(err, oldAfterItem.ID, oldAfterItem.UserID, "update item after moved item's old position")
		}
	}

	if i.Before == nil {
		newAfterItem := new(Item)
		err := tx.Where(&Item{Before: nil}).First(newAfterItem).Error
		if err != nil {
			tx.Rollback()
			return convertError(err, i.ID, i.UserID, "find item after moved item")
		}
		i.After = &newAfterItem.ID
	} else {
		newBeforeItem := new(Item)
		err := tx.Where(&Item{ID: *i.Before}).First(newBeforeItem).Error
		if err != nil {
			tx.Rollback()
			return convertError(err, i.ID, i.UserID, "find item before moved item")
		}
		i.After = newBeforeItem.After
	}

	// Update a item before moved item's new position
	if i.Before != nil {
		newBeforeItem := Item{
			ID:     *i.Before,
			UserID: i.UserID,
			After:  &i.ID,
		}

		err = tx.Model(&newBeforeItem).Updates(map[string]interface{}{
			"after": convertData(newBeforeItem.After),
		}).Error

		if err != nil {
			tx.Rollback()
			return convertError(err, newBeforeItem.ID, newBeforeItem.UserID, "update item before moved item's new position")
		}
	}

	// Update a item after moved item's new position
	if i.After != nil {
		newAfterItem := Item{
			ID:     *i.After,
			UserID: i.UserID,
			Before: &i.ID,
		}

		err = tx.Model(&newAfterItem).Updates(map[string]interface{}{
			"before": convertData(newAfterItem.Before),
		}).Error

		if err != nil {
			tx.Rollback()
			return convertError(err, newAfterItem.ID, newAfterItem.UserID, "update item after moved item's new position")
		}
	}

	err = tx.Model(&i).Updates(map[string]interface{}{
		"listID": i.ListID,
		"after":  convertData(i.After),
		"before": convertData(i.Before),
	}).Error

	if err != nil {
		tx.Rollback()
		return convertError(err, i.ID, i.UserID, "move item")
	}
	tx.Commit()
	return nil
}

// Delete removes a Item from DB.
func (m *ItemDBManager) Delete(item model.Item) error {
	if err := validatePrimaryKeys("item", item.ID, item.UserID); err != nil {
		return err
	}

	tx := m.db.Begin()

	i := Item{}
	i.convertFrom(item)

	deletedItem := new(Item)
	if err := tx.Delete(&i).First(deletedItem).Error; err != nil {
		return convertError(err, i.ID, i.UserID, "delete item")
	}

	// Update item before deleted item
	if deletedItem.Before != nil {
		err := tx.Model(&Item{ID: *deletedItem.Before, UserID: deletedItem.UserID}).Updates(map[string]interface{}{
			"after": convertData(deletedItem.After),
		}).Error
		if err != nil {
			tx.Rollback()
			return model.ServerError{
				UserID: deletedItem.UserID,
				Err:    err,
				ID:     *deletedItem.Before,
				Act:    "update item before deleted item",
			}
		}
	}

	// Update item after deleted item
	if deletedItem.After != nil {
		err := tx.Model(&Item{ID: *deletedItem.After, UserID: deletedItem.UserID}).Updates(map[string]interface{}{
			"before": convertData(deletedItem.Before),
		}).Error
		if err != nil {
			tx.Rollback()
			return model.ServerError{
				UserID: item.UserID,
				Err:    err,
				ID:     item.ID,
				Act:    "update item after deleted item",
			}
		}
	}

	tx.Commit()
	return nil
}

// Find gets a Item had specific ID from DB.
func (m *ItemDBManager) Find(item model.Item) (model.Item, error) {
	if err := validatePrimaryKeys("item", item.ID, item.UserID); err != nil {
		return model.Item{}, err
	}

	r := Item{}
	if err := m.db.Where(&Item{ID: item.ID, UserID: item.UserID}).First(&r).Error; err != nil {
		return model.Item{}, convertError(err, item.ID, item.UserID, "find item")
	}
	return r.convertTo(), nil
}

// FindItems gets all Items in a specific List from DB.
func (m *ItemDBManager) FindItems(list model.List) (model.Items, error) {
	if err := validatePrimaryKeys("list", list.ID, list.UserID); err != nil {
		return model.Items{}, err
	}

	r := Items{}
	if err := m.db.Where(&Item{ListID: list.ID, UserID: list.UserID}).Find(&r).Error; err != nil {
		return model.Items{}, model.ServerError{
			UserID: list.UserID,
			Err:    err,
			ID:     list.ID,
			Act:    "find items in list",
		}
	}

	items := model.Items{}
	for _, ri := range r {
		items = append(items, ri.convertTo())
	}

	return items, nil
}
