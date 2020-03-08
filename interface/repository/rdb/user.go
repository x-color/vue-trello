package rdb

import (
	"time"

	"github.com/jinzhu/gorm"
	"github.com/x-color/vue-trello/model"
	"github.com/x-color/vue-trello/usecase"
)

// User is User data model for DB.
type User struct {
	ID   string `gorm:"primary_key"`
	Name string
	// Password is raw. It should be hash.
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
type UserDBManager struct{}

func newUserDBManager(db *gorm.DB) UserDBManager {
	db.AutoMigrate(&User{})
	return UserDBManager{}
}

// Create registers a User to DB.
func (*UserDBManager) Create(tx usecase.Transaction, user model.User) error {
	if err := validatePrimaryKeys("user", user.ID); err != nil {
		return err
	}

	u := User{}
	u.convertFrom(user)

	if err := tx.DB().(*gorm.DB).Create(&u).Error; err != nil {
		return model.ServerError{
			UserID: u.ID,
			Err:    err,
			ID:     u.ID,
			Act:    "create user",
		}
	}
	return nil
}

// Find gets a User.
func (*UserDBManager) Find(tx usecase.Transaction, conditions map[string]interface{}) (model.User, error) {
	r := User{}
	if err := tx.DB().(*gorm.DB).Where(queryForUser(conditions)).First(&r).Error; err != nil {
		userID := "(No-ID)"
		if v, ok := conditions["user_id"]; ok {
			userID = v.(string)
		}
		return model.User{}, convertError(err, userID, userID, "find user")
	}
	return r.convertTo(), nil
}

func queryForUser(data map[string]interface{}) map[string]interface{} {
	query := make(map[string]interface{})
	if v, ok := data["ID"]; ok {
		query["id"] = v
	}
	if v, ok := data["Name"]; ok {
		query["name"] = v
	}
	if v, ok := data["Password"]; ok {
		query["password"] = v
	}
	return query
}
