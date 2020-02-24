package rdb

import (
	"errors"
	"os"

	"github.com/jinzhu/gorm"

	// SQLite driver
	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

// DBManager includes DB managers for all data model.
type DBManager struct {
	ItemDBManager  ItemDBManager
	ListDBManager  ListDBManager
	BoardDBManager BoardDBManager
	UserDBManager  UserDBManager
	TagDBManager   TagDBManager
}

// NewDBManager generates new DB manager.
func NewDBManager() (DBManager, error) {
	db, err := gorm.Open("sqlite3", os.Getenv("DB_PATH"))
	if err != nil {
		return DBManager{}, errors.New("failed to connect database")
	}

	dbm := DBManager{
		ItemDBManager:  newItemDBManager(db),
		ListDBManager:  newListDBManager(db),
		BoardDBManager: newBoardDBManager(db),
		UserDBManager:  newUserDBManager(db),
		TagDBManager:   newTagDBManager(db),
	}
	return dbm, nil
}
