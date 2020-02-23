package usecase

import (
	"github.com/google/uuid"
	"github.com/x-color/vue-trello/model"
)

// BoardUsecase is interface. It defines to control a Board.
type BoardUsecase interface {
	Get(board model.Board) (model.Board, error)
	Create(board model.Board) (model.Board, error)
	Delete(board model.Board) error
	Update(board model.Board) (model.Board, error)
}

// BoardInteractor includes repogitories and a logger.
type BoardInteractor struct {
	boardRepo BoardRepogitory
	listRepo  ListRepogitory
	itemRepo  ItemRepogitory
	logger    Logger
}

// NewBoardInteractor generates new interactor for a Board.
func NewBoardInteractor(
	boardRepo BoardRepogitory,
	listRepo ListRepogitory,
	itemRepo ItemRepogitory,
	logger Logger,
) (BoardInteractor, error) {
	i := BoardInteractor{
		boardRepo: boardRepo,
		listRepo:  listRepo,
		itemRepo:  itemRepo,
		logger:    logger,
	}
	return i, nil
}

// Create saves new Board to a repogitory and returns created Board.
func (i *BoardInteractor) Create(board model.Board) (model.Board, error) {
	board.ID = uuid.New().String()
	if err := i.validateBoard(board); err != nil {
		return model.Board{}, err
	}

	if err := i.boardRepo.Create(board); err != nil {
		return model.Board{}, err
	}
	return board, nil
}

// Delete removes Board in repogitory.
func (i *BoardInteractor) Delete(board model.Board) error {
	if board.ID == "" {
		return model.InvalidContentError{
			Err: nil,
			ID:  board.ID,
			Act: "validate contents",
		}
	}
	if err := i.boardRepo.Delete(board); err != nil {
		return err
	}
	return nil
}

// Update replaces a Board and returns new Board.
func (i *BoardInteractor) Update(board model.Board) (model.Board, error) {
	if err := i.validateBoard(board); err != nil {
		return model.Board{}, err
	}

	if err := i.boardRepo.Update(board); err != nil {
		return model.Board{}, err
	}
	return board, nil
}

// Get returns Board embedded all data.
func (i *BoardInteractor) Get(board model.Board) (model.Board, error) {
	board, err := i.boardRepo.Find(board)
	if err != nil {
		return model.Board{}, err
	}

	// Get Lists in Board.
	lists, err := i.listRepo.FindLists(board)
	if err != nil {
		return model.Board{}, err
	}
	board.Lists = lists

	// Get Items in Lists.
	for j, list := range lists {
		items, err := i.itemRepo.FindItems(list)
		if err != nil {
			return model.Board{}, err
		}
		board.Lists[j].Items = items
	}

	return board, nil
}

// GetBoards returns User's Boards.
func (i *BoardInteractor) GetBoards(user model.User) (model.Boards, error) {
	boards, err := i.boardRepo.FindBoards(user)
	if err != nil {
		return model.Boards{}, err
	}
	return boards, nil
}

func (i *BoardInteractor) validateBoard(board model.Board) error {
	if board.ID == "" || board.Title == "" || board.UserID == "" {
		return model.InvalidContentError{
			Err: nil,
			ID:  board.ID,
			Act: "validate contents",
		}
	}

	// Validate color of Board.
	for _, c := range model.COLORS {
		if board.Color == c {
			return nil
		}
	}

	return model.InvalidContentError{
		Err: nil,
		ID:  board.ID,
		Act: "validate color",
	}
}
