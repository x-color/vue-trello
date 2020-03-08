package rdb

import (
	"time"

	"github.com/jinzhu/gorm"
	"github.com/x-color/vue-trello/model"
	"github.com/x-color/vue-trello/usecase"
)

// Tag is Tag data model for DB.
type Tag struct {
	ID        string `gorm:"primary_key"`
	Name      string
	Color     string
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time
}

func (t *Tag) convertFrom(tag model.Tag) {
	t.ID = tag.ID
	t.Name = tag.Name
	t.Color = string(tag.Color)
}

func (t *Tag) convertTo() model.Tag {
	tag := model.Tag{
		ID:    t.ID,
		Name:  t.Name,
		Color: model.Color(t.Color),
	}
	return tag
}

// Tags is a slice of Tag data model.
type Tags []Tag

// TagDBManager is DB manager for Tag.
type TagDBManager struct{}

func newTagDBManager(db *gorm.DB) TagDBManager {
	db.AutoMigrate(&Tag{})
	return TagDBManager{}
}

// Create registers a Tag to DB.
func (*TagDBManager) Create(tx usecase.Transaction, tag model.Tag) error {
	if err := validatePrimaryKeys("tag", tag.ID); err != nil {
		return err
	}

	t := Tag{}
	t.convertFrom(tag)

	if err := tx.DB().(*gorm.DB).Create(&t).Error; err != nil {
		return model.ServerError{
			UserID: "(No-ID)",
			Err:    err,
			ID:     t.ID,
			Act:    "create tag",
		}
	}
	return nil
}

// Find get Tags.
func (*TagDBManager) Find(tx usecase.Transaction, conditions map[string]interface{}) (model.Tags, error) {
	r := Tags{}
	if err := tx.DB().(*gorm.DB).Where(queryForTag(conditions)).Find(&r).Error; err != nil {
		return model.Tags{}, model.ServerError{
			UserID: "(No-ID)",
			Err:    err,
			ID:     "(No-ID)",
			Act:    "find all tags",
		}
	}

	tags := model.Tags{}
	for _, rt := range r {
		tags = append(tags, rt.convertTo())
	}

	return tags, nil
}

func queryForTag(data map[string]interface{}) map[string]interface{} {
	query := make(map[string]interface{})
	if v, ok := data["ID"]; ok {
		query["id"] = v
	}
	if v, ok := data["Name"]; ok {
		query["name"] = v
	}
	if v, ok := data["Color"]; ok {
		query["color"] = v
	}
	return query
}
