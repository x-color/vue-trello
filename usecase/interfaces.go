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
	FindByID(tx Transaction, ID, userID string) (model.Item, error)
	Find(tx Transaction, conditions map[string]interface{}) (model.Items, error)
}

// ListRepository is interface. It defines CURD methods for List.
type ListRepository interface {
	Create(list model.List) error
	Update(list model.List) error
	Delete(list model.List) error
	Move(list model.List) error
	Find(list model.List) (model.List, error)
	FindLists(board model.Board) (model.Lists, error)
}

// BoardRepository is interface. It defines CURD methods for Board.
type BoardRepository interface {
	Create(board model.Board) error
	Update(board model.Board) error
	Delete(board model.Board) error
	Move(board model.Board) error
	Find(board model.Board) (model.Board, error)
	FindBoards(user model.User) (model.Boards, error)
}

// UserRepository is interface. It defines CR methods for User.
type UserRepository interface {
	Create(user model.User) error
	FindByName(user model.User) (model.User, error)
}

// TagRepository is interface. It defines CR methods for Tag.
type TagRepository interface {
	Create(tag model.Tag) error
	FindAll() (model.Tags, error)
}
