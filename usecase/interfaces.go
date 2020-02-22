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

// ItemRepogitory is interface. It defines CURD methods for Item.
type ItemRepogitory interface {
	Create(item model.Item) error
	Update(item model.Item) error
	Delete(item model.Item) error
	Find(item model.Item) (model.Item, error)
	FindItems(list model.List) (model.Items, error)
}

// ListRepogitory is interface. It defines CURD methods for List.
type ListRepogitory interface {
	Create(list model.List) error
	Update(list model.List) error
	Delete(list model.List) error
	Find(list model.List) (model.List, error)
	FindLists(board model.Board) (model.Lists, error)
}

// BoardRepogitory is interface. It defines CURD methods for Board.
type BoardRepogitory interface {
	Create(board model.Board) error
	Update(board model.Board) error
	Delete(board model.Board) error
	Find(board model.Board) (model.Board, error)
	FindBoards(user model.User) (model.Boards, error)
}

// UserRepogitory is interface. It defines CR methods for User.
type UserRepogitory interface {
	Create(user model.User) error
	Find(user model.User) (model.User, error)
}

// TagRepogitory is interface. It defines CR methods for Tag.
type TagRepogitory interface {
	Create(tag model.Tag) error
	FindAll() (model.Tags, error)
}
