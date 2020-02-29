package rdb

import (
	"time"

	"github.com/jinzhu/gorm"
	"github.com/x-color/vue-trello/model"
)

// User is User data model for DB.
type User struct {
	ID   string `gorm:"primary_key"`
	Name string
	// Password is row. It will be hash.
	Password  string
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time
}

func (u *User) convertFrom(user model.User) {
	u.ID = user.ID
	u.Name = user.Name
	u.Password = user.Password
}

func (u *User) convertTo() model.User {
	user := model.User{
		ID:       u.ID,
		Name:     u.Name,
		Password: u.Password,
	}
	return user
}

// Users is a slice of User data model.
type Users []User

// UserDBManager is DB manager for User.
type UserDBManager struct {
	db *gorm.DB
}

func newUserDBManager(db *gorm.DB) UserDBManager {
	db.AutoMigrate(&User{})
	return UserDBManager{
		db: db,
	}
}

// Create registers a User to DB.
func (m *UserDBManager) Create(user model.User) error {
	if err := validatePrimaryKeys("user", user.ID); err != nil {
		return err
	}

	u := User{}
	u.convertFrom(user)

	if err := m.db.Create(&u).Error; err != nil {
		return model.ServerError{
			UserID: u.ID,
			Err:    err,
			ID:     u.ID,
			Act:    "create user",
		}
	}
	return nil
}

// FindByName gets a User had specific name from DB.
func (m *UserDBManager) FindByName(user model.User) (model.User, error) {
	r := User{}
	if err := m.db.Where(&User{Name: user.Name}).First(&r).Error; err != nil {
		if gorm.IsRecordNotFoundError(err) {
			return model.User{}, model.NotFoundError{
				UserID: "(No-ID)",
				Err:    err,
				ID:     "(No-ID)",
				Act:    "find user",
			}
		}
		return model.User{}, model.ServerError{
			UserID: user.ID,
			Err:    err,
			ID:     user.ID,
			Act:    "find user",
		}
	}
	return r.convertTo(), nil
}
