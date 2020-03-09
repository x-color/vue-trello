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
	Move(board model.Board) error
}

// BoardInteractor includes repogitories and a logger.
type BoardInteractor struct {
	txRepo    TransactionRepository
	boardRepo BoardRepository
	listRepo  ListRepository
	itemRepo  ItemRepository
	logger    Logger
}

// NewBoardInteractor generates new interactor for a Board.
func NewBoardInteractor(
	txRepo TransactionRepository,
	boardRepo BoardRepository,
	listRepo ListRepository,
	itemRepo ItemRepository,
	logger Logger,
) (BoardInteractor, error) {
	i := BoardInteractor{
		txRepo:    txRepo,
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

	tx := i.txRepo.BeginTransaction(true)
	i.logger.Info(formatLogMsg(board.UserID, "Start transaction"))

	// Get last board
	boards, err := i.boardRepo.Find(tx, map[string]interface{}{
		"UserID": board.UserID,
		"After":  "",
	})
	if err != nil {
		tx.Rollback()
		i.logger.Info(formatLogMsg(board.UserID, "Rollback transaction"))
		logError(i.logger, err)
		return model.Board{}, err
	}

	if len(boards) > 0 {
		lastBoard := boards[0]
		i.logger.Info(formatLogMsg(board.UserID, "Find last board("+lastBoard.ID+")"))

		board.Before = lastBoard.ID

		query := map[string]interface{}{
			"After": board.ID,
		}
		if err := i.boardRepo.Update(tx, lastBoard, query); err != nil {
			tx.Rollback()
			i.logger.Info(formatLogMsg(board.UserID, "Rollback transaction"))
			logError(i.logger, err)
			return model.Board{}, err
		}
		i.logger.Info(formatLogMsg(board.UserID, "Update board("+lastBoard.ID+")"))
	} else {
		i.logger.Info(formatLogMsg(board.UserID, "Find no board"))
	}

	if err := i.boardRepo.Create(tx, board); err != nil {
		tx.Rollback()
		i.logger.Info(formatLogMsg(board.UserID, "Rollback transaction"))
		logError(i.logger, err)
		return model.Board{}, err
	}
	i.logger.Info(formatLogMsg(board.UserID, "Create board("+board.ID+")"))

	tx.Commit()
	i.logger.Info(formatLogMsg(board.UserID, "Commit transaction"))

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

	tx := i.txRepo.BeginTransaction(true)
	i.logger.Info(formatLogMsg(board.UserID, "Start transaction"))

	// Get board's info (e.g. board.Before, board.After...) and rewrite 'board'.
	board, err := i.boardRepo.FindByID(tx, board.ID, board.UserID)
	if err != nil {
		tx.Rollback()
		i.logger.Info(formatLogMsg(board.UserID, "Rollback transaction"))
		logError(i.logger, err)
		return err
	}
	i.logger.Info(formatLogMsg(board.UserID, "Find board("+board.ID+")"))

	if board.Before != "" {
		// Update board before deleting board
		before := model.Board{
			ID:     board.Before,
			UserID: board.UserID,
		}
		query := map[string]interface{}{
			"After": board.After,
		}
		if err := i.boardRepo.Update(tx, before, query); err != nil {
			tx.Rollback()
			i.logger.Info(formatLogMsg(board.UserID, "Rollback transaction"))
			logError(i.logger, err)
			return err
		}
		i.logger.Info(formatLogMsg(board.UserID, "Update board("+before.ID+") before deleted board("+board.ID+")"))
	}

	if board.After != "" {
		// Update board after deleting board
		after := model.Board{
			ID:     board.After,
			UserID: board.UserID,
		}
		query := map[string]interface{}{
			"Before": board.Before,
		}
		if err := i.boardRepo.Update(tx, after, query); err != nil {
			tx.Rollback()
			i.logger.Info(formatLogMsg(board.UserID, "Rollback transaction"))
			logError(i.logger, err)
			return err
		}
		i.logger.Info(formatLogMsg(board.UserID, "Update board("+after.ID+") after deleted board("+board.ID+")"))
	}

	if err := i.boardRepo.Delete(tx, board); err != nil {
		tx.Rollback()
		i.logger.Info(formatLogMsg(board.UserID, "Rollback transaction"))
		logError(i.logger, err)
		return err
	}
	i.logger.Info(formatLogMsg(board.UserID, "Delete board("+board.ID+")"))

	// Get lists in deleted board
	lists, err := i.listRepo.Find(tx, map[string]interface{}{
		"UserID":  board.UserID,
		"BoardID": board.ID,
	})
	if err != nil {
		tx.Rollback()
		i.logger.Info(formatLogMsg(board.UserID, "Rollback transaction"))
		logError(i.logger, err)
		return err
	}
	i.logger.Info(formatLogMsg(board.UserID, "Find lists in deleted board("+board.ID+")"))

	for _, list := range lists {
		// Delete list
		if err := i.listRepo.Delete(tx, list); err != nil {
			tx.Rollback()
			i.logger.Info(formatLogMsg(list.UserID, "Rollback transaction"))
			logError(i.logger, err)
			return err
		}
		i.logger.Info(formatLogMsg(board.UserID, "Delete lists in deleted board("+board.ID+")"))

		items, err := i.itemRepo.Find(tx, map[string]interface{}{
			"UserID": list.UserID,
			"ListID": list.ID,
		})
		if err != nil {
			tx.Rollback()
			i.logger.Info(formatLogMsg(board.UserID, "Rollback transaction"))
			logError(i.logger, err)
			return err
		}
		i.logger.Info(formatLogMsg(board.UserID, "Find items in deleted list("+list.ID+")"))

		for _, item := range items {
			if err := i.itemRepo.Delete(tx, item); err != nil {
				tx.Rollback()
				i.logger.Info(formatLogMsg(item.UserID, "Rollback transaction"))
				logError(i.logger, err)
				return err
			}
		}
		i.logger.Info(formatLogMsg(board.UserID, "Delete items in deleted list("+list.ID+")"))
	}

	tx.Commit()
	i.logger.Info(formatLogMsg(board.UserID, "Commit transaction"))

	return nil
}

// Update replaces a Board and returns new Board.
func (i *BoardInteractor) Update(board model.Board) (model.Board, error) {
	if err := i.validateBoard(board); err != nil {
		logError(i.logger, err)
		return model.Board{}, err
	}

	query := map[string]interface{}{
		"Title": board.Title,
		"Text":  board.Text,
		"Color": string(board.Color),
	}

	tx := i.txRepo.BeginTransaction(false)
	if err := i.boardRepo.Update(tx, board, query); err != nil {
		logError(i.logger, err)
		return model.Board{}, err
	}
	i.logger.Info(formatLogMsg(board.UserID, "Update board("+board.ID+")"))

	return board, nil
}

// Move moves Boards.
func (i *BoardInteractor) Move(board model.Board) error {
	board.Title = "dummy title"
	board.Color = model.RED
	if err := i.validateBoard(board); err != nil {
		logError(i.logger, err)
		return err
	}
	board.Title = ""

	tx := i.txRepo.BeginTransaction(true)
	i.logger.Info(formatLogMsg(board.UserID, "Start transaction"))

	// Get a board to move
	old, err := i.boardRepo.FindByID(tx, board.ID, board.UserID)
	if err != nil {
		tx.Rollback()
		i.logger.Info(formatLogMsg(board.UserID, "Rollback transaction"))
		logError(i.logger, err)
		return err
	}
	i.logger.Info(formatLogMsg(board.UserID, "Find board("+board.ID+") to move"))

	// Get a board before board to move
	if old.Before != "" {
		beforeOld := model.Board{
			ID:     old.Before,
			UserID: old.UserID,
		}
		query := map[string]interface{}{
			"After": old.After,
		}
		if err := i.boardRepo.Update(tx, beforeOld, query); err != nil {
			tx.Rollback()
			i.logger.Info(formatLogMsg(board.UserID, "Rollback transaction"))
			logError(i.logger, err)
			return err
		}
		i.logger.Info(formatLogMsg(board.UserID, "Find a board("+beforeOld.ID+") before board("+board.ID+") to move"))
	}

	// Get a board after board to move
	if old.After != "" {
		afterOld := model.Board{
			ID:     old.After,
			UserID: old.UserID,
		}
		query := map[string]interface{}{
			"Before": old.Before,
		}
		if err := i.boardRepo.Update(tx, afterOld, query); err != nil {
			tx.Rollback()
			i.logger.Info(formatLogMsg(board.UserID, "Rollback transaction"))
			logError(i.logger, err)
			return err
		}
		i.logger.Info(formatLogMsg(board.UserID, "Find a board("+afterOld.ID+") after board("+board.ID+") to move"))
	}

	// Get a board after moved board
	if board.Before == "" {
		conditions := map[string]interface{}{
			"UserID": board.UserID,
			"Before": "",
		}
		l, err := i.boardRepo.Find(tx, conditions)
		if err != nil {
			tx.Rollback()
			i.logger.Info(formatLogMsg(board.UserID, "Rollback transaction"))
			logError(i.logger, err)
			return err
		}
		if len(l) != 1 {
			tx.Rollback()
			i.logger.Info(formatLogMsg(board.UserID, "Rollback transaction"))
			err = model.ServerError{
				ID:     board.ID,
				UserID: board.UserID,
				Err:    nil,
				Act:    "find a board before moved board",
			}
			logError(i.logger, err)
			return err
		}
		board.After = l[0].ID
	} else {
		before, err := i.boardRepo.FindByID(tx, board.Before, board.UserID)
		if err != nil {
			tx.Rollback()
			i.logger.Info(formatLogMsg(board.UserID, "Rollback transaction"))
			logError(i.logger, err)
			return err
		}
		board.After = before.After
	}
	i.logger.Info(formatLogMsg(board.UserID, "Find a board("+board.After+") after moved board("+board.ID+")"))

	// Update a board before moved board
	if board.Before != "" {
		before := model.Board{
			ID:     board.Before,
			UserID: board.UserID,
		}
		query := map[string]interface{}{
			"After": board.ID,
		}
		if err := i.boardRepo.Update(tx, before, query); err != nil {
			tx.Rollback()
			i.logger.Info(formatLogMsg(board.UserID, "Rollback transaction"))
			logError(i.logger, err)
			return err
		}
		i.logger.Info(formatLogMsg(board.UserID, "Update board("+before.ID+") before board("+board.ID+")"))
	}

	// Update a board after moved board
	if board.After != "" {
		after := model.Board{
			ID:     board.After,
			UserID: board.UserID,
		}
		query := map[string]interface{}{
			"Before": board.ID,
		}
		if err := i.boardRepo.Update(tx, after, query); err != nil {
			tx.Rollback()
			i.logger.Info(formatLogMsg(board.UserID, "Rollback transaction"))
			logError(i.logger, err)
			return err
		}
		i.logger.Info(formatLogMsg(board.UserID, "Update board("+after.ID+") after board("+board.ID+")"))
	}

	// Move board
	query := map[string]interface{}{
		"Before": board.Before,
		"After":  board.After,
	}
	if err := i.boardRepo.Update(tx, board, query); err != nil {
		tx.Rollback()
		i.logger.Info(formatLogMsg(board.UserID, "Rollback transaction"))
		logError(i.logger, err)
		return err
	}
	i.logger.Info(formatLogMsg(board.UserID, "Move board("+board.ID+") after board("+board.Before+")"))

	tx.Commit()
	i.logger.Info(formatLogMsg(board.UserID, "Commit transaction"))

	return nil
}

// Get returns Board embedded all data.
func (i *BoardInteractor) Get(board model.Board) (model.Board, error) {
	tx := i.txRepo.BeginTransaction(false)

	board, err := i.boardRepo.FindByID(tx, board.ID, board.UserID)
	if err != nil {
		logError(i.logger, err)
		return model.Board{}, err
	}
	i.logger.Info(formatLogMsg(board.UserID, "Find board("+board.ID+")"))

	// Get Lists in Board.
	lists, err := i.listRepo.Find(tx, map[string]interface{}{
		"BoardID": board.ID,
		"UserID":  board.UserID,
	})
	if err != nil {
		logError(i.logger, err)
		return model.Board{}, err
	}
	board.Lists = sortLists(lists)
	i.logger.Info(formatLogMsg(board.UserID, "Find lists in board("+board.ID+")"))

	// Get Items in Lists.
	for j, list := range board.Lists {
		items, err := i.itemRepo.Find(tx, map[string]interface{}{
			"ListID": list.ID,
			"UserID": list.UserID,
		})
		if err != nil {
			logError(i.logger, err)
			return model.Board{}, err
		}
		board.Lists[j].Items = sortItems(items)
	}
	i.logger.Info(formatLogMsg(board.UserID, "Find items in board("+board.ID+")"))

	i.logger.Info(formatLogMsg(board.UserID, "Get board("+board.ID+")"))
	return board, nil
}

// GetBoards returns User's Boards.
func (i *BoardInteractor) GetBoards(user model.User) (model.Boards, error) {
	tx := i.txRepo.BeginTransaction(false)

	boards, err := i.boardRepo.Find(tx, map[string]interface{}{
		"UserID": user.ID,
	})
	if err != nil {
		logError(i.logger, err)
		return model.Boards{}, err
	}
	i.logger.Info(formatLogMsg(user.ID, "Get boards"))
	return sortBoards(boards), nil
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

func sortBoards(boards model.Boards) model.Boards {
	l := map[string]model.Board{}
	for _, b := range boards {
		l[b.Before] = b
	}

	board := l[""]
	r := model.Boards{}
	for range boards {
		r = append(r, board)
		board = l[board.ID]
	}

	return r
}
