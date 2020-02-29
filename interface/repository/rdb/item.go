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
	i.Tags = &ts

	if item.Text == "" {
		i.Text = nil
	} else {
		i.Text = &item.Text
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

	i := Item{}
	i.convertFrom(item)

	if err := m.db.Create(&i).Error; err != nil {
		return model.ServerError{
			UserID: item.UserID,
			Err:    err,
			ID:     i.ID,
			Act:    "create item",
		}
	}
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
		"tags":  i.Tags,
	}).Error

	if err != nil {
		if gorm.IsRecordNotFoundError(err) {
			return model.NotFoundError{
				UserID: item.UserID,
				Err:    err,
				ID:     i.ID,
				Act:    "update item",
			}
		}
		return model.ServerError{
			UserID: item.UserID,
			Err:    err,
			ID:     i.ID,
			Act:    "update item",
		}
	}
	return nil
}

// Delete removes a Item from DB.
func (m *ItemDBManager) Delete(item model.Item) error {
	if err := validatePrimaryKeys("item", item.ID, item.UserID); err != nil {
		return err
	}

	i := Item{}
	i.convertFrom(item)

	if err := m.db.Delete(&i).Error; err != nil {
		if gorm.IsRecordNotFoundError(err) {
			return model.NotFoundError{
				UserID: item.UserID,
				Err:    err,
				ID:     i.ID,
				Act:    "delete item",
			}
		}
		return model.ServerError{
			UserID: item.UserID,
			Err:    err,
			ID:     i.ID,
			Act:    "delete item",
		}
	}
	return nil
}

// Find gets a Item had specific ID from DB.
func (m *ItemDBManager) Find(item model.Item) (model.Item, error) {
	if err := validatePrimaryKeys("item", item.ID, item.UserID); err != nil {
		return model.Item{}, err
	}

	r := Item{}
	if err := m.db.Where(&Item{ID: item.ID, UserID: item.UserID}).First(&r).Error; err != nil {
		if gorm.IsRecordNotFoundError(err) {
			return model.Item{}, model.NotFoundError{
				UserID: item.UserID,
				Err:    err,
				ID:     item.ID,
				Act:    "find item",
			}
		}
		return model.Item{}, model.ServerError{
			UserID: item.UserID,
			Err:    err,
			ID:     item.ID,
			Act:    "find item",
		}
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
