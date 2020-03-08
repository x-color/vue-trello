package rdb

import (
	"time"

	"github.com/jinzhu/gorm"
	"github.com/x-color/vue-trello/model"
	"github.com/x-color/vue-trello/usecase"
)

// Board is Board data model for DB.
type Board struct {
	ID        string `gorm:"primary_key"`
	UserID    string `gorm:"primary_key"`
	Title     string
	Text      *string
	Color     string
	Before    *string
	After     *string
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time
}

func (b *Board) convertFrom(board model.Board) {
	b.ID = board.ID
	b.UserID = board.UserID
	b.Title = board.Title
	b.Color = string(board.Color)

	if board.Text == "" {
		b.Text = nil
	} else {
		b.Text = &board.Text
	}

	if board.Before == "" {
		b.Before = nil
	} else {
		b.Before = &board.Before
	}

	if board.After == "" {
		b.After = nil
	} else {
		b.After = &board.After
	}
}

func (b *Board) convertTo() model.Board {
	board := model.Board{
		ID:     b.ID,
		UserID: b.UserID,
		Title:  b.Title,
		Color:  model.Color(b.Color),
	}

	if b.Text == nil {
		board.Text = ""
	} else {
		board.Text = *b.Text
	}

	if b.After == nil {
		board.After = ""
	} else {
		board.After = *b.After
	}

	if b.Before == nil {
		board.Before = ""
	} else {
		board.Before = *b.Before
	}

	return board
}

// Boards is a slice of Board data model.
type Boards []Board

// BoardDBManager is DB manager for Board.
type BoardDBManager struct{}

func newBoardDBManager(db *gorm.DB) BoardDBManager {
	db.AutoMigrate(&Board{})
	return BoardDBManager{}
}

// Create registers a Board to DB.
func (*BoardDBManager) Create(tx usecase.Transaction, board model.Board) error {
	if err := validatePrimaryKeys("board", board.ID, board.UserID); err != nil {
		return err
	}

	b := Board{}
	b.convertFrom(board)

	if err := tx.DB().(*gorm.DB).Create(&b).Error; err != nil {
		return model.ServerError{
			UserID: board.UserID,
			Err:    err,
			ID:     board.ID,
			Act:    "create board",
		}
	}

	return nil
}

// Update updates all fields of specific Board in DB.
func (*BoardDBManager) Update(tx usecase.Transaction, board model.Board, updates map[string]interface{}) error {
	if err := validatePrimaryKeys("board", board.ID, board.UserID); err != nil {
		return err
	}

	b := Board{}
	b.convertFrom(board)
	err := tx.DB().(*gorm.DB).Model(&b).Updates(queryForBoard(updates)).Error

	if err != nil {
		return convertError(err, b.ID, b.UserID, "update board")
	}
	return nil
}

// Delete removes a Board from DB.
func (*BoardDBManager) Delete(tx usecase.Transaction, board model.Board) error {
	if err := validatePrimaryKeys("board", board.ID, board.UserID); err != nil {
		return err
	}

	b := Board{}
	b.convertFrom(board)

	if err := tx.DB().(*gorm.DB).Delete(&b).Error; err != nil {
		return convertError(err, b.ID, b.UserID, "delete board")
	}
	return nil
}

// FindByID gets a Board had specific ID from DB.
func (*BoardDBManager) FindByID(tx usecase.Transaction, id, userID string) (model.Board, error) {
	if err := validatePrimaryKeys("board", id, userID); err != nil {
		return model.Board{}, err
	}

	r := Board{}
	if err := tx.DB().(*gorm.DB).Where(&Board{ID: id, UserID: userID}).First(&r).Error; err != nil {
		return model.Board{}, convertError(err, id, userID, "find board")
	}
	return r.convertTo(), nil
}

// Find gets Boards.
func (*BoardDBManager) Find(tx usecase.Transaction, conditions map[string]interface{}) (model.Boards, error) {
	r := Boards{}
	if err := tx.DB().(*gorm.DB).Where(queryForBoard(conditions)).Find(&r).Error; err != nil {
		userID := "(No-ID)"
		if v, ok := conditions["user_id"]; ok {
			userID = v.(string)
		}
		id := "(No-ID)"
		if v, ok := conditions["id"]; ok {
			id = v.(string)
		}
		return model.Boards{}, model.ServerError{
			UserID: userID,
			Err:    err,
			ID:     id,
			Act:    "find items",
		}
	}

	boards := model.Boards{}
	for _, rb := range r {
		boards = append(boards, rb.convertTo())
	}

	return boards, nil
}

func queryForBoard(data map[string]interface{}) map[string]interface{} {
	query := make(map[string]interface{})
	if v, ok := data["ID"]; ok {
		query["id"] = v
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
	if v, ok := data["Color"]; ok {
		query["color"] = v.(string)
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
