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

	tx := m.db.Begin()

	beforeBoard := new(Board)

	if err := tx.Where(map[string]interface{}{"user_id": board.UserID, "after": nil}).First(beforeBoard).Error; err != nil {
		if gorm.IsRecordNotFoundError(err) {
			beforeBoard = nil
		} else {
			tx.Rollback()
			return model.ServerError{
				UserID: board.UserID,
				Err:    err,
				ID:     board.ID,
				Act:    "find board before last board",
			}
		}
	}

	b := Board{}
	b.convertFrom(board)
	b.After = nil
	if beforeBoard == nil {
		b.Before = nil
	} else {
		b.Before = &beforeBoard.ID

		err := tx.Model(beforeBoard).Updates(map[string]interface{}{
			"after": b.ID,
		}).Error
		if err != nil {
			tx.Rollback()
			return convertError(err, beforeBoard.ID, beforeBoard.UserID, "update board to create new board")
		}
	}

	if err := tx.Create(&b).Error; err != nil {
		tx.Rollback()
		return model.ServerError{
			UserID: board.UserID,
			Err:    err,
			ID:     board.ID,
			Act:    "create board",
		}
	}
	tx.Commit()
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
		return convertError(err, b.ID, b.UserID, "update board")
	}
	return nil
}

// Move updates board's position data..
func (m *BoardDBManager) Move(board model.Board) error {
	if err := validatePrimaryKeys("board", board.ID, board.UserID); err != nil {
		return err
	}

	b := Board{}
	b.convertFrom(board)

	tx := m.db.Begin()

	oldBoard := new(Board)
	err := tx.Where(&Board{ID: b.ID}).First(oldBoard).Error
	if err != nil {
		tx.Rollback()
		return convertError(err, b.ID, b.UserID, "find moved board")
	}

	// Update a board before moved board's old position
	if oldBoard.Before != nil {
		oldBeforeBoard := Board{
			ID:     *oldBoard.Before,
			UserID: oldBoard.UserID,
			After:  oldBoard.After,
		}

		err = tx.Model(&oldBeforeBoard).Updates(map[string]interface{}{
			"after": convertData(oldBeforeBoard.After),
		}).Error

		if err != nil {
			tx.Rollback()
			return convertError(err, oldBeforeBoard.ID, oldBeforeBoard.UserID, "update board before moved board's old position")
		}
	}

	// Update a board after moved board's old position
	if oldBoard.After != nil {
		oldAfterBoard := Board{
			ID:     *oldBoard.After,
			UserID: oldBoard.UserID,
			Before: oldBoard.Before,
		}

		err = tx.Model(&oldAfterBoard).Updates(map[string]interface{}{
			"before": convertData(oldAfterBoard.Before),
		}).Error

		if err != nil {
			tx.Rollback()
			return convertError(err, oldAfterBoard.ID, oldAfterBoard.UserID, "update board after moved board's old position")
		}
	}

	if b.Before == nil {
		newAfterBoard := new(Board)
		err := tx.Where(&Board{Before: nil}).First(newAfterBoard).Error
		if err != nil {
			tx.Rollback()
			return convertError(err, b.ID, b.UserID, "find board after moved board")
		}
		b.After = &newAfterBoard.ID
	} else {
		newBeforeBoard := new(Board)
		err := tx.Where(&Board{ID: *b.Before}).First(newBeforeBoard).Error
		if err != nil {
			tx.Rollback()
			return convertError(err, b.ID, b.UserID, "find board before moved board")
		}
		b.After = newBeforeBoard.After
	}

	// Update a board before moved board's new position
	if b.Before != nil {
		newBeforeBoard := Board{
			ID:     *b.Before,
			UserID: b.UserID,
			After:  &b.ID,
		}

		err = tx.Model(&newBeforeBoard).Updates(map[string]interface{}{
			"after": convertData(newBeforeBoard.After),
		}).Error

		if err != nil {
			tx.Rollback()
			return convertError(err, newBeforeBoard.ID, newBeforeBoard.UserID, "update board before moved board's new position")
		}
	}

	// Update a board after moved board's new position
	if b.After != nil {
		newAfterBoard := Board{
			ID:     *b.After,
			UserID: b.UserID,
			Before: &b.ID,
		}

		err = tx.Model(&newAfterBoard).Updates(map[string]interface{}{
			"before": convertData(newAfterBoard.Before),
		}).Error

		if err != nil {
			tx.Rollback()
			return convertError(err, newAfterBoard.ID, newAfterBoard.UserID, "update board after moved board's new position")
		}
	}

	err = tx.Model(&b).Updates(map[string]interface{}{
		"after":  convertData(b.After),
		"before": convertData(b.Before),
	}).Error

	if err != nil {
		tx.Rollback()
		return convertError(err, b.ID, b.UserID, "move board")
	}
	tx.Commit()
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

	deletedBoard := new(Board)
	if err := tx.Where(&b).First(deletedBoard).Error; err != nil {
		tx.Rollback()
		return convertError(err, b.ID, b.UserID, "find deleting board")
	}

	if err := tx.Delete(&b).Error; err != nil {
		tx.Rollback()
		return convertError(err, b.ID, b.UserID, "delete board")
	}

	// Update board before deleted board
	if deletedBoard.Before != nil {
		err := tx.Model(&Board{ID: *deletedBoard.Before, UserID: deletedBoard.UserID}).Updates(map[string]interface{}{
			"after": convertData(deletedBoard.After),
		}).Error
		if err != nil {
			tx.Rollback()
			return model.ServerError{
				UserID: deletedBoard.UserID,
				Err:    err,
				ID:     *deletedBoard.Before,
				Act:    "update board before deleted board",
			}
		}
	}

	// Update board after deleted board
	if deletedBoard.After != nil {
		err := tx.Model(&Board{ID: *deletedBoard.After, UserID: deletedBoard.UserID}).Updates(map[string]interface{}{
			"before": convertData(deletedBoard.Before),
		}).Error
		if err != nil {
			tx.Rollback()
			return model.ServerError{
				UserID: deletedBoard.UserID,
				Err:    err,
				ID:     *deletedBoard.After,
				Act:    "update board after deleted board",
			}
		}
	}

	// Remove Lists in removed Board
	lists := model.Lists{}
	if err := tx.Where(&List{BoardID: b.ID, UserID: b.UserID}).Find(&lists).Error; err != nil {
		tx.Rollback()
		return convertError(err, b.ID, b.UserID, "find deleting lists in deleted board")
	}

	if err := tx.Where(&List{BoardID: b.ID, UserID: b.UserID}).Delete(List{}).Error; err != nil {
		tx.Rollback()
		return model.ServerError{
			UserID: b.UserID,
			Err:    err,
			ID:     b.ID,
			Act:    "delete lists in deleted board",
		}
	}

	// Remove Items in removed Lists
	for _, list := range lists {
		if err := tx.Where(&Item{ListID: list.ID, UserID: list.UserID}).Delete(Item{}).Error; err != nil {
			tx.Rollback()
			return model.ServerError{
				UserID: b.UserID,
				Err:    err,
				ID:     b.ID,
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
		return model.Board{}, convertError(err, board.ID, board.UserID, "find board")
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
