package rdb

import (
	"time"

	"github.com/jinzhu/gorm"
	"github.com/x-color/vue-trello/model"
)

// Board is Board data model for DB.
type Board struct {
	ID        string `gorm:"primary_key"`
	UserID    string `gorm:"primary_key"`
	Title     string
	Text      *string
	Color     string
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
	return board
}

// Boards is a slice of Board data model.
type Boards []Board

// BoardDBManager is DB manager for Board.
type BoardDBManager struct {
	db *gorm.DB
}

func newBoardDBManager(db *gorm.DB) BoardDBManager {
	db.AutoMigrate(&Board{})
	return BoardDBManager{
		db: db,
	}
}

// Create registers a Board to DB.
func (m *BoardDBManager) Create(board model.Board) error {
	if err := validatePrimaryKeys("board", board.ID, board.UserID); err != nil {
		return err
	}

	b := Board{}
	b.convertFrom(board)
	if err := m.db.Create(&b).Error; err != nil {
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
func (m *BoardDBManager) Update(board model.Board) error {
	if err := validatePrimaryKeys("board", board.ID, board.UserID); err != nil {
		return err
	}

	b := Board{}
	b.convertFrom(board)
	err := m.db.Model(&b).Updates(map[string]interface{}{
		"title": b.Title,
		"text":  convertData(b.Text),
		"color": b.Color,
	}).Error

	if err != nil {
		if gorm.IsRecordNotFoundError(err) {
			return model.NotFoundError{
				UserID: board.UserID,
				Err:    err,
				ID:     b.ID,
				Act:    "update board",
			}
		}
		return model.ServerError{
			UserID: board.UserID,
			Err:    err,
			ID:     b.ID,
			Act:    "update board",
		}
	}
	return nil
}

// Delete removes a Board from DB.
func (m *BoardDBManager) Delete(board model.Board) error {
	if err := validatePrimaryKeys("board", board.ID, board.UserID); err != nil {
		return err
	}

	b := Board{}
	b.convertFrom(board)

	tx := m.db.Begin()

	if err := tx.Delete(&b).Error; err != nil {
		tx.Rollback()
		if gorm.IsRecordNotFoundError(err) {
			return model.NotFoundError{
				UserID: board.UserID,
				Err:    err,
				ID:     board.ID,
				Act:    "delete board",
			}
		}
		return model.ServerError{
			UserID: board.UserID,
			Err:    err,
			ID:     board.ID,
			Act:    "delete board",
		}
	}

	// Remove Lists in removed Board
	lists := model.Lists{}
	if err := tx.Where(&List{BoardID: board.ID}).Delete(List{}).Find(&lists).Error; err != nil {
		tx.Rollback()
		return model.ServerError{
			UserID: board.UserID,
			Err:    err,
			ID:     board.ID,
			Act:    "delete lists in deleted board",
		}
	}

	// Remove Items in removed Lists
	for _, list := range lists {
		if err := tx.Where(&Item{ListID: list.ID}).Delete(Item{}).Error; err != nil {
			tx.Rollback()
			return model.ServerError{
				UserID: board.UserID,
				Err:    err,
				ID:     board.ID,
				Act:    "delete items in deleted board",
			}
		}
	}

	tx.Commit()
	return nil
}

// Find gets a Board had specific ID from DB.
func (m *BoardDBManager) Find(board model.Board) (model.Board, error) {
	if err := validatePrimaryKeys("board", board.ID, board.UserID); err != nil {
		return model.Board{}, err
	}

	r := Board{}
	if err := m.db.Where(&Board{ID: board.ID, UserID: board.UserID}).First(&r).Error; err != nil {
		if gorm.IsRecordNotFoundError(err) {
			return model.Board{}, model.NotFoundError{
				UserID: board.UserID,
				Err:    err,
				ID:     board.ID,
				Act:    "find board",
			}
		}
		return model.Board{}, model.ServerError{
			UserID: board.UserID,
			Err:    err,
			ID:     board.ID,
			Act:    "find board",
		}
	}
	return r.convertTo(), nil
}

// FindBoards gets User's all Boards from DB.
func (m *BoardDBManager) FindBoards(user model.User) (model.Boards, error) {
	if err := validatePrimaryKeys("user", user.ID); err != nil {
		return model.Boards{}, err
	}

	r := Boards{}
	if err := m.db.Where(&Board{UserID: user.ID}).Find(&r).Error; err != nil {
		return model.Boards{}, model.ServerError{
			UserID: user.ID,
			Err:    err,
			ID:     user.ID,
			Act:    "find boards",
		}
	}

	boards := model.Boards{}
	for _, rb := range r {
		boards = append(boards, rb.convertTo())
	}
	return boards, nil
}
