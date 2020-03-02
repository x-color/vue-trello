package usecase

import (
	"github.com/google/uuid"
	"github.com/x-color/vue-trello/model"
)

// BoardUsecase is interface. It defines to control a Board.
type BoardUsecase interface {
	Get(board model.Board) (model.Board, error)
	GetBoards(user model.User) (model.Boards, error)
	Create(board model.Board) (model.Board, error)
	Delete(board model.Board) error
	Update(board model.Board) (model.Board, error)
}

// BoardInteractor includes repogitories and a logger.
type BoardInteractor struct {
	boardRepo BoardRepository
	listRepo  ListRepository
	itemRepo  ItemRepository
	logger    Logger
}

// NewBoardInteractor generates new interactor for a Board.
func NewBoardInteractor(
	boardRepo BoardRepository,
	listRepo ListRepository,
	itemRepo ItemRepository,
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

// Create saves new Board to a repository and returns created Board.
func (i *BoardInteractor) Create(board model.Board) (model.Board, error) {
	board.ID = uuid.New().String()
	if err := i.validateBoard(board); err != nil {
		logError(i.logger, err)
		return model.Board{}, err
	}

	if err := i.boardRepo.Create(board); err != nil {
		logError(i.logger, err)
		return model.Board{}, err
	}
	i.logger.Info(formatLogMsg(board.UserID, "Create board("+board.ID+")"))
	return board, nil
}

// Delete removes Board in repository.
func (i *BoardInteractor) Delete(board model.Board) error {
	if board.ID == "" {
		i.logger.Info(formatLogMsg(board.UserID, "Invalid board. ID is empty"))
		return model.InvalidContentError{
			UserID: board.UserID,
			Err:    nil,
			ID:     "(No-ID)",
			Act:    "validate board id",
		}
	}
	if err := i.boardRepo.Delete(board); err != nil {
		logError(i.logger, err)
		return err
	}
	i.logger.Info(formatLogMsg(board.UserID, "Delete board("+board.ID+")"))
	return nil
}

// Update replaces a Board and returns new Board.
func (i *BoardInteractor) Update(board model.Board) (model.Board, error) {
	if err := i.validateBoard(board); err != nil {
		logError(i.logger, err)
		return model.Board{}, err
	}

	if err := i.boardRepo.Update(board); err != nil {
		logError(i.logger, err)
		return model.Board{}, err
	}
	i.logger.Info(formatLogMsg(board.UserID, "Update board("+board.ID+")"))
	return board, nil
}

// Get returns Board embedded all data.
func (i *BoardInteractor) Get(board model.Board) (model.Board, error) {
	board, err := i.boardRepo.Find(board)
	if err != nil {
		logError(i.logger, err)
		return model.Board{}, err
	}

	// Get Lists in Board.
	lists, err := i.listRepo.FindLists(board)
	if err != nil {
		logError(i.logger, err)
		return model.Board{}, err
	}
	board.Lists = sortLists(lists)

	// Get Items in Lists.
	for j, list := range lists {
		items, err := i.itemRepo.FindItems(list)
		if err != nil {
			logError(i.logger, err)
			return model.Board{}, err
		}
		board.Lists[j].Items = sortItems(items)
	}

	i.logger.Info(formatLogMsg(board.UserID, "Get board("+board.ID+")"))
	return board, nil
}

// GetBoards returns User's Boards.
func (i *BoardInteractor) GetBoards(user model.User) (model.Boards, error) {
	boards, err := i.boardRepo.FindBoards(user)
	if err != nil {
		logError(i.logger, err)
		return model.Boards{}, err
	}
	i.logger.Info(formatLogMsg(user.ID, "Get boards"))
	return boards, nil
}

func (i *BoardInteractor) validateBoard(board model.Board) error {
	if board.ID == "" || board.Title == "" || board.UserID == "" {
		return model.InvalidContentError{
			UserID: board.UserID,
			Err:    nil,
			ID:     board.ID,
			Act:    "validate contents in board",
		}
	}

	// Validate color of Board.
	for _, c := range model.COLORS {
		if board.Color == c {
			return nil
		}
	}

	return model.InvalidContentError{
		UserID: board.UserID,
		Err:    nil,
		ID:     board.ID,
		Act:    "validate color of board",
	}
}
