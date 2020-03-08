package usecase

import (
	"github.com/x-color/vue-trello/model"
)

// Logger is interface. It defines logging methods.
type Logger interface {
	Debug(msg string)
	Info(msg string)
	Error(msg string)
}

// Transaction is interface. It defined transaction methods.
type Transaction interface {
	Commit()
	Rollback()
	DB() interface{}
}

// TransactionRepository is interface. It defines Transaction getter.
type TransactionRepository interface {
	BeginTransaction(on bool) Transaction
}

// ItemRepository is interface. It defines CURD methods for Item.
type ItemRepository interface {
	Create(tx Transaction, item model.Item) error
	Update(tx Transaction, item model.Item, updates map[string]interface{}) error
	Delete(tx Transaction, item model.Item) error
	FindByID(tx Transaction, id, userID string) (model.Item, error)
	Find(tx Transaction, conditions map[string]interface{}) (model.Items, error)
}

// ListRepository is interface. It defines CURD methods for List.
type ListRepository interface {
	Create(tx Transaction, list model.List) error
	Update(tx Transaction, list model.List, updates map[string]interface{}) error
	Delete(tx Transaction, list model.List) error
	FindByID(tx Transaction, id, UserID string) (model.List, error)
	Find(tx Transaction, conditions map[string]interface{}) (model.Lists, error)
}

// BoardRepository is interface. It defines CURD methods for Board.
type BoardRepository interface {
	Create(tx Transaction, board model.Board) error
	Update(tx Transaction, board model.Board, updates map[string]interface{}) error
	Delete(tx Transaction, board model.Board) error
	FindByID(tx Transaction, id, userID string) (model.Board, error)
	Find(tx Transaction, condititons map[string]interface{}) (model.Boards, error)
}

// UserRepository is interface. It defines CR methods for User.
type UserRepository interface {
	Create(tx Transaction, user model.User) error
	Find(tx Transaction, conditions map[string]interface{}) (model.User, error)
}

// TagRepository is interface. It defines CR methods for Tag.
type TagRepository interface {
	Create(tx Transaction, tag model.Tag) error
	Find(tx Transaction, conditions map[string]interface{}) (model.Tags, error)
}
