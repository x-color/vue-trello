package rdb

import (
	"github.com/jinzhu/gorm"
	"github.com/x-color/vue-trello/model"
)

// BoardDBManager is DB manager for Board.
type BoardDBManager struct {
	db *gorm.DB
}

// Create registers a Board to DB.
func (m *BoardDBManager) Create(board model.Board) error {
	if err := m.db.Create(&board).Error; err != nil {
		return model.ServerError{
			Err: err,
			ID:  board.ID,
			Act: "create board",
		}
	}
	return nil
}

// Update updates all fields of specific Board in DB.
func (m *BoardDBManager) Update(board model.Board) error {
	if validatePrimaryKey(board.ID) {
		return model.NotFoundError{
			Err: nil,
			ID:  "(No ID)",
			Act: "validate board",
		}
	}

	err := m.db.Model(&board).Where(&model.Board{ID: board.ID, UserID: board.UserID}).Updates(map[string]interface{}{
		"title": board.Title,
		"text":  convertData(board.Text),
		"color": board.Color,
	}).Error

	if err != nil {
		if gorm.IsRecordNotFoundError(err) {
			return model.NotFoundError{
				Err: err,
				ID:  board.ID,
				Act: "update board",
			}
		}
		return model.ServerError{
			Err: err,
			ID:  board.ID,
			Act: "update board",
		}
	}
	return nil
}

// Delete removes a Board from DB.
func (m *BoardDBManager) Delete(board model.Board) error {
	if validatePrimaryKey(board.ID) {
		return model.NotFoundError{
			Err: nil,
			ID:  "(No ID)",
			Act: "validate board",
		}
	}

	tx := m.db.Begin()

	if err := tx.Where(&model.Board{ID: board.ID, UserID: board.UserID}).Delete(&board).Error; err != nil {
		tx.Rollback()
		if gorm.IsRecordNotFoundError(err) {
			return model.NotFoundError{
				Err: err,
				ID:  board.ID,
				Act: "delete board",
			}
		}
		return model.ServerError{
			Err: err,
			ID:  board.ID,
			Act: "delete board",
		}
	}

	// Remove Lists in removed Board
	lists := model.Lists{}
	if err := tx.Where(&model.List{BoardID: board.ID}).Delete(model.List{}).Find(&lists).Error; err != nil {
		tx.Rollback()
		return model.ServerError{
			Err: err,
			ID:  board.ID,
			Act: "delete board",
		}
	}

	// Remove Items in removed Lists
	for _, list := range lists {
		if err := tx.Where(&model.Item{ListID: list.ID}).Delete(model.Item{}).Error; err != nil {
			tx.Rollback()
			return model.ServerError{
				Err: err,
				ID:  board.ID,
				Act: "delete board",
			}
		}
	}

	tx.Commit()
	return nil
}

// Find gets a Board had specific ID from DB.
func (m *BoardDBManager) Find(board model.Board) (model.Board, error) {
	r := model.Board{}
	if err := m.db.Where(&model.Board{ID: board.ID, UserID: board.UserID}).First(&r).Error; err != nil {
		if gorm.IsRecordNotFoundError(err) {
			return model.Board{}, model.NotFoundError{
				Err: err,
				ID:  board.ID,
				Act: "find board",
			}
		}
		return model.Board{}, model.ServerError{
			Err: err,
			ID:  board.ID,
			Act: "find board",
		}
	}
	return r, nil
}

// FindBoards gets User's all Boards from DB.
func (m *BoardDBManager) FindBoards(user model.User) (model.Boards, error) {
	if validatePrimaryKey(user.ID) {
		return model.Boards{}, model.NotFoundError{
			Err: nil,
			ID:  "(No ID)",
			Act: "validate board",
		}
	}
	r := model.Boards{}
	if err := m.db.Where(&model.Board{UserID: user.ID}).Find(r).Error; err != nil {
		return model.Boards{}, model.ServerError{
			Err: err,
			ID:  user.ID,
			Act: "find boards",
		}
	}
	return r, nil
}

func convertData(data string) interface{} {
	if data == "" {
		return gorm.Expr("NULL")
	}
	return data
}
