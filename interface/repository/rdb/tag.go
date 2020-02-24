package rdb

import (
	"time"

	"github.com/jinzhu/gorm"
	"github.com/x-color/vue-trello/model"
)

// Tag is Tag data model for DB.
type Tag struct {
	ID        string `gorm:"primary_key"`
	Name      string
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time
}

func (t *Tag) convertFrom(tag model.Tag) {
	t.ID = tag.ID
	t.Name = tag.Name
}

func (t *Tag) convertTo() model.Tag {
	tag := model.Tag{
		ID:   t.ID,
		Name: t.Name,
	}
	return tag
}

// Tags is a slice of Tag data model.
type Tags []Tag

// TagDBManager is DB manager for Tag.
type TagDBManager struct {
	db *gorm.DB
}

// Create registers a Tag to DB.
func (m *TagDBManager) Create(tag model.Tag) error {
	t := Tag{}
	t.convertFrom(tag)

	if err := m.db.Create(&t).Error; err != nil {
		return model.ServerError{
			Err: err,
			ID:  t.ID,
			Act: "create tag",
		}
	}
	return nil
}

// FindAll gets all Tags from DB.
func (m *TagDBManager) FindAll() (model.Tags, error) {
	r := Tags{}
	if err := m.db.Find(&r).Error; err != nil {
		return model.Tags{}, model.ServerError{
			Err: err,
			ID:  "(No ID)",
			Act: "find all tags",
		}
	}

	tags := model.Tags{}
	for _, rt := range r {
		tags = append(tags, rt.convertTo())
	}

	return tags, nil
}
