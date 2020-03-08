package rdb

import (
	"errors"
	"os"

	"github.com/jinzhu/gorm"
	"github.com/x-color/vue-trello/model"

	// SQLite driver
	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

// DBManager includes DB managers for all data model.
type DBManager struct {
	TransactionManager TransactionManager
	ItemDBManager      ItemDBManager
	ListDBManager      ListDBManager
	BoardDBManager     BoardDBManager
	UserDBManager      UserDBManager
	TagDBManager       TagDBManager
}

// NewDBManager generates new DB manager.
func NewDBManager() (DBManager, error) {
	db, err := gorm.Open("sqlite3", os.Getenv("DB_PATH"))
	if err != nil {
		return DBManager{}, errors.New("failed to connect database")
	}

	dbm := DBManager{
		TransactionManager: newTransactionManager(db),
		ItemDBManager:      newItemDBManager(db),
		ListDBManager:      newListDBManager(db),
		BoardDBManager:     newBoardDBManager(db),
		UserDBManager:      newUserDBManager(db),
		TagDBManager:       newTagDBManager(db),
	}
	return dbm, nil
}

func convertData(data interface{}) interface{} {
	if data == nil {
		return gorm.Expr("NULL")
	}
	return data
}

func validatePrimaryKeys(targetName string, keys ...string) error {
	for _, k := range keys {
		if k == "" {
			return model.NotFoundError{
				UserID: "(No-ID)",
				Err:    nil,
				ID:     "(No-ID)",
				Act:    "validate " + targetName + "'s primary key",
			}
		}
	}
	return nil
}

func convertError(err error, id, userID, act string) error {
	if gorm.IsRecordNotFoundError(err) {
		return model.NotFoundError{
			UserID: userID,
			Err:    err,
			ID:     id,
			Act:    act,
		}
	}
	return model.ServerError{
		UserID: userID,
		Err:    err,
		ID:     id,
		Act:    act,
	}
}
