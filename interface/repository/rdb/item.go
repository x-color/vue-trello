package rdb

import (
	"strings"
	"time"

	"github.com/jinzhu/gorm"
	"github.com/x-color/vue-trello/model"
	"github.com/x-color/vue-trello/usecase"
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
type ItemDBManager struct{}

func newItemDBManager(db *gorm.DB) ItemDBManager {
	db.AutoMigrate(&Item{})
	return ItemDBManager{}
}

// Create registers a Item to DB.
func (*ItemDBManager) Create(tx usecase.Transaction, item model.Item) error {
	if err := validatePrimaryKeys("item", item.ID, item.UserID); err != nil {
		return err
	}

	i := Item{}
	i.convertFrom(item)

	if err := tx.DB().(*gorm.DB).Create(&i).Error; err != nil {
		return model.ServerError{
			UserID: i.UserID,
			Err:    err,
			ID:     i.ID,
			Act:    "create item",
		}
	}

	return nil
}

// Update updates all fields of specific Item in DB.
func (*ItemDBManager) Update(tx usecase.Transaction, item model.Item, updates map[string]interface{}) error {
	if err := validatePrimaryKeys("item", item.ID, item.UserID); err != nil {
		return err
	}

	i := Item{}
	i.convertFrom(item)

	err := tx.DB().(*gorm.DB).Model(&i).Updates(queryForItem(updates)).Error

	if err != nil {
		return convertError(err, i.ID, i.UserID, "update item")
	}
	return nil
}

// Delete removes a Item from DB.
func (*ItemDBManager) Delete(tx usecase.Transaction, item model.Item) error {
	if err := validatePrimaryKeys("item", item.ID, item.UserID); err != nil {
		return err
	}

	i := Item{}
	i.convertFrom(item)

	if err := tx.DB().(*gorm.DB).Delete(&i).Error; err != nil {
		return convertError(err, i.ID, i.UserID, "delete item")
	}
	return nil
}

// FindByID gets a Item had specific ID from DB.
func (*ItemDBManager) FindByID(tx usecase.Transaction, id, userID string) (model.Item, error) {
	if err := validatePrimaryKeys("item", id, userID); err != nil {
		return model.Item{}, err
	}

	r := Item{}
	if err := tx.DB().(*gorm.DB).Where(&Item{ID: id, UserID: userID}).First(&r).Error; err != nil {
		return model.Item{}, convertError(err, id, userID, "find item")
	}
	return r.convertTo(), nil
}

// Find gets Items.
func (*ItemDBManager) Find(tx usecase.Transaction, conditions map[string]interface{}) (model.Items, error) {
	r := Items{}
	if err := tx.DB().(*gorm.DB).Where(queryForItem(conditions)).Find(&r).Error; err != nil {
		userID := "(No-ID)"
		if v, ok := conditions["user_id"]; ok {
			userID = v.(string)
		}
		id := "(No-ID)"
		if v, ok := conditions["id"]; ok {
			id = v.(string)
		}
		return model.Items{}, model.ServerError{
			UserID: userID,
			Err:    err,
			ID:     id,
			Act:    "find items",
		}
	}

	items := model.Items{}
	for _, ri := range r {
		items = append(items, ri.convertTo())
	}

	return items, nil
}

func queryForItem(data map[string]interface{}) map[string]interface{} {
	query := make(map[string]interface{})
	if v, ok := data["ID"]; ok {
		query["id"] = v
	}
	if v, ok := data["ListID"]; ok {
		query["list_id"] = v
	}
	if v, ok := data["UserID"]; ok {
		query["user_id"] = v
	}
	if v, ok := data["Title"]; ok {
		query["title"] = v
	}
	if v, ok := data["Text"]; ok {
		if v.(string) == "" {
			query["text"] = nil
		} else {
			query["text"] = v
		}
	}
	if v, ok := data["Tags"]; ok {
		tags := v.([]string)
		if len(tags) == 0 {
			query["tags"] = nil
		} else {
			query["tags"] = strings.Join(tags, ",")
		}
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
