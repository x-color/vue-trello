package usecase

import (
	"github.com/google/uuid"
	"github.com/x-color/vue-trello/model"
)

// ListUsecase is interface. It defines to control a List.
type ListUsecase interface {
	Create(list model.List) (model.List, error)
	Delete(list model.List) error
	Update(list model.List) (model.List, error)
	Move(list model.List) error
}

// ListInteractor includes repogitories and a logger.
type ListInteractor struct {
	txRepo    TransactionRepository
	itemRepo  ItemRepository
	listRepo  ListRepository
	boardRepo BoardRepository
	logger    Logger
}

// NewListInteractor generates new interactor for a List.
func NewListInteractor(
	txRepo TransactionRepository,
	itemRepo ItemRepository,
	listRepo ListRepository,
	boardRepo BoardRepository,
	logger Logger,
) (ListInteractor, error) {
	i := ListInteractor{
		txRepo:    txRepo,
		itemRepo:  itemRepo,
		listRepo:  listRepo,
		boardRepo: boardRepo,
		logger:    logger,
	}
	return i, nil
}

// Create saves new List to a repository and returns created List.
func (i *ListInteractor) Create(list model.List) (model.List, error) {
	list.ID = uuid.New().String()
	if err := i.validateList(list); err != nil {
		logError(i.logger, err)
		return model.List{}, err
	}

	tx := i.txRepo.BeginTransaction(true)
	i.logger.Info(formatLogMsg(list.UserID, "Start transaction"))

	// Get last list in board
	lists, err := i.listRepo.Find(tx, map[string]interface{}{
		"BoardID": list.BoardID,
		"UserID":  list.UserID,
		"After":   "",
	})
	if err != nil {
		tx.Rollback()
		i.logger.Info(formatLogMsg(list.UserID, "Rollback transaction"))
		logError(i.logger, err)
		return model.List{}, err
	}

	if len(lists) > 0 {
		lastList := lists[0]
		i.logger.Info(formatLogMsg(list.UserID, "Find last list("+lastList.ID+") in board("+lastList.BoardID+")"))

		list.Before = lastList.ID

		query := map[string]interface{}{
			"After": list.ID,
		}
		if err := i.listRepo.Update(tx, lastList, query); err != nil {
			tx.Rollback()
			i.logger.Info(formatLogMsg(list.UserID, "Rollback transaction"))
			logError(i.logger, err)
			return model.List{}, err
		}
		i.logger.Info(formatLogMsg(list.UserID, "Update list("+lastList.ID+")"))
	} else {
		i.logger.Info(formatLogMsg(list.UserID, "Find no list in board("+list.BoardID+")"))
	}

	if err := i.listRepo.Create(tx, list); err != nil {
		tx.Rollback()
		i.logger.Info(formatLogMsg(list.UserID, "Rollback transaction"))
		logError(i.logger, err)
		return model.List{}, err
	}
	i.logger.Info(formatLogMsg(list.UserID, "Create list("+list.ID+")"))

	tx.Commit()
	i.logger.Info(formatLogMsg(list.UserID, "Commit transaction"))

	return list, nil
}

// Delete removes List in repository.
func (i *ListInteractor) Delete(list model.List) error {
	if list.ID == "" {
		i.logger.Info(formatLogMsg(list.UserID, "Invalid list. ID is empty"))
		return model.InvalidContentError{
			UserID: list.UserID,
			Err:    nil,
			ID:     "(No-ID)",
			Act:    "validate list id",
		}
	}

	tx := i.txRepo.BeginTransaction(true)
	i.logger.Info(formatLogMsg(list.UserID, "Start transaction"))

	// Get list's info (e.g. list.Before, list.After...) and rewrite 'list'.
	list, err := i.listRepo.FindByID(tx, list.ID, list.UserID)
	if err != nil {
		tx.Rollback()
		i.logger.Info(formatLogMsg(list.UserID, "Rollback transaction"))
		logError(i.logger, err)
		return err
	}
	i.logger.Info(formatLogMsg(list.UserID, "Find list("+list.ID+")"))

	if list.Before != "" {
		// Update list before deleting list
		before := model.List{
			ID:     list.Before,
			UserID: list.UserID,
		}
		query := map[string]interface{}{
			"After": list.After,
		}
		if err := i.listRepo.Update(tx, before, query); err != nil {
			tx.Rollback()
			i.logger.Info(formatLogMsg(list.UserID, "Rollback transaction"))
			logError(i.logger, err)
			return err
		}
		i.logger.Info(formatLogMsg(list.UserID, "Update list("+before.ID+") before deleted list("+list.ID+")"))
	}

	if list.After != "" {
		// Update list after deleting list
		after := model.List{
			ID:     list.After,
			UserID: list.UserID,
		}
		query := map[string]interface{}{
			"Before": list.Before,
		}
		if err := i.listRepo.Update(tx, after, query); err != nil {
			tx.Rollback()
			i.logger.Info(formatLogMsg(list.UserID, "Rollback transaction"))
			logError(i.logger, err)
			return err
		}
		i.logger.Info(formatLogMsg(list.UserID, "Update list("+after.ID+") after deleted list("+list.ID+")"))
	}

	if err := i.listRepo.Delete(tx, list); err != nil {
		tx.Rollback()
		i.logger.Info(formatLogMsg(list.UserID, "Rollback transaction"))
		logError(i.logger, err)
		return err
	}
	i.logger.Info(formatLogMsg(list.UserID, "Delete list("+list.ID+")"))

	items, err := i.itemRepo.Find(tx, map[string]interface{}{
		"UserID": list.UserID,
		"ListID": list.ID,
	})
	if err != nil {
		tx.Rollback()
		i.logger.Info(formatLogMsg(list.UserID, "Rollback transaction"))
		logError(i.logger, err)
		return err
	}
	i.logger.Info(formatLogMsg(list.UserID, "Find items in deleted list("+list.ID+")"))

	for _, item := range items {
		if err := i.itemRepo.Delete(tx, item); err != nil {
			tx.Rollback()
			i.logger.Info(formatLogMsg(item.UserID, "Rollback transaction"))
			logError(i.logger, err)
			return err
		}
	}
	i.logger.Info(formatLogMsg(list.UserID, "Delete items in deleted list("+list.ID+")"))

	tx.Commit()
	i.logger.Info(formatLogMsg(list.UserID, "Commit transaction"))

	return nil
}

// Update replaces a List and returns new List.
func (i *ListInteractor) Update(list model.List) (model.List, error) {
	if err := i.validateList(list); err != nil {
		logError(i.logger, err)
		return model.List{}, err
	}

	query := map[string]interface{}{
		"Title": list.Title,
	}

	tx := i.txRepo.BeginTransaction(false)
	if err := i.listRepo.Update(tx, list, query); err != nil {
		logError(i.logger, err)
		return model.List{}, err
	}
	i.logger.Info(formatLogMsg(list.UserID, "Update list("+list.ID+")"))

	return list, nil
}

// Move moves Items.
func (i *ListInteractor) Move(list model.List) error {
	list.Title = "dummy title"
	if err := i.validateList(list); err != nil {
		logError(i.logger, err)
		return err
	}
	list.Title = ""

	tx := i.txRepo.BeginTransaction(true)
	i.logger.Info(formatLogMsg(list.UserID, "Start transaction"))

	// Get a list to move
	if list.Before != "" {
		before, err := i.listRepo.FindByID(tx, list.Before, list.UserID)
		if err != nil {
			tx.Rollback()
			i.logger.Info(formatLogMsg(list.UserID, "Rollback transaction"))
			logError(i.logger, err)
			return err
		}
		i.logger.Info(formatLogMsg(list.UserID, "Find a list("+before.ID+") before list to move"))
		if list.BoardID != before.BoardID {
			tx.Rollback()
			i.logger.Info(formatLogMsg(list.UserID, "Rollback transaction"))
			i.logger.Info(formatLogMsg(list.UserID, "Invalid list id("+list.BoardID+"). Does not equal before list's list id("+before.BoardID+")"))
			return err
		}
	}

	// Get a list to move
	old, err := i.listRepo.FindByID(tx, list.ID, list.UserID)
	if err != nil {
		tx.Rollback()
		i.logger.Info(formatLogMsg(list.UserID, "Rollback transaction"))
		logError(i.logger, err)
		return err
	}
	i.logger.Info(formatLogMsg(list.UserID, "Find list("+list.ID+") to move"))

	// Get a list before list to move
	if old.Before != "" {
		beforeOld := model.List{
			ID:     old.Before,
			UserID: old.UserID,
		}
		query := map[string]interface{}{
			"After": old.After,
		}
		if err := i.listRepo.Update(tx, beforeOld, query); err != nil {
			tx.Rollback()
			i.logger.Info(formatLogMsg(list.UserID, "Rollback transaction"))
			logError(i.logger, err)
			return err
		}
		i.logger.Info(formatLogMsg(list.UserID, "Find a list("+beforeOld.ID+") before list("+list.ID+") to move"))
	}

	// Get a list after list to move
	if old.After != "" {
		afterOld := model.List{
			ID:     old.After,
			UserID: old.UserID,
		}
		query := map[string]interface{}{
			"Before": old.Before,
		}
		if err := i.listRepo.Update(tx, afterOld, query); err != nil {
			tx.Rollback()
			i.logger.Info(formatLogMsg(list.UserID, "Rollback transaction"))
			logError(i.logger, err)
			return err
		}
		i.logger.Info(formatLogMsg(list.UserID, "Find a list("+afterOld.ID+") after list("+list.ID+") to move"))
	}

	// Get a list after moved list
	if list.Before == "" {
		conditions := map[string]interface{}{
			"BoardID": list.BoardID,
			"UserID":  list.UserID,
			"Before":  "",
		}
		l, err := i.listRepo.Find(tx, conditions)
		if err != nil {
			tx.Rollback()
			i.logger.Info(formatLogMsg(list.UserID, "Rollback transaction"))
			logError(i.logger, err)
			return err
		}
		if len(l) == 1 {
			list.After = l[0].ID
		} else if len(l) > 1 {
			tx.Rollback()
			i.logger.Info(formatLogMsg(list.UserID, "Rollback transaction"))
			err = model.ServerError{
				ID:     list.ID,
				UserID: list.UserID,
				Err:    nil,
				Act:    "find a list before moved list",
			}
			logError(i.logger, err)
			return err
		}
	} else {
		before, err := i.listRepo.FindByID(tx, list.Before, list.UserID)
		if err != nil {
			tx.Rollback()
			i.logger.Info(formatLogMsg(list.UserID, "Rollback transaction"))
			logError(i.logger, err)
			return err
		}
		list.After = before.After
	}
	i.logger.Info(formatLogMsg(list.UserID, "Find a list("+list.After+") after moved list("+list.ID+")"))

	// Update a list before moved list
	if list.Before != "" {
		before := model.List{
			ID:     list.Before,
			UserID: list.UserID,
		}
		query := map[string]interface{}{
			"After": list.ID,
		}
		if err := i.listRepo.Update(tx, before, query); err != nil {
			tx.Rollback()
			i.logger.Info(formatLogMsg(list.UserID, "Rollback transaction"))
			logError(i.logger, err)
			return err
		}
		i.logger.Info(formatLogMsg(list.UserID, "Update list("+before.ID+") before list("+list.ID+") in board("+list.BoardID+")"))
	}

	// Update a list after moved list
	if list.After != "" {
		after := model.List{
			ID:     list.After,
			UserID: list.UserID,
		}
		query := map[string]interface{}{
			"Before": list.ID,
		}
		if err := i.listRepo.Update(tx, after, query); err != nil {
			tx.Rollback()
			i.logger.Info(formatLogMsg(list.UserID, "Rollback transaction"))
			logError(i.logger, err)
			return err
		}
		i.logger.Info(formatLogMsg(list.UserID, "Update list("+after.ID+") after list("+list.ID+") in board("+list.BoardID+")"))
	}

	// Move list
	query := map[string]interface{}{
		"BoardID": list.BoardID,
		"Before":  list.Before,
		"After":   list.After,
	}
	if err := i.listRepo.Update(tx, list, query); err != nil {
		tx.Rollback()
		i.logger.Info(formatLogMsg(list.UserID, "Rollback transaction"))
		logError(i.logger, err)
		return err
	}
	i.logger.Info(formatLogMsg(list.UserID, "Move list("+list.ID+") after list("+list.Before+") in board("+list.BoardID+")"))

	tx.Commit()
	i.logger.Info(formatLogMsg(list.UserID, "Commit transaction"))

	return nil
}

func (i *ListInteractor) validateList(list model.List) error {
	if list.ID == "" || list.Title == "" || list.BoardID == "" || list.UserID == "" {
		return model.InvalidContentError{
			UserID: list.UserID,
			Err:    nil,
			ID:     list.ID,
			Act:    "validate contents in list",
		}
	}

	tx := i.txRepo.BeginTransaction(false)

	_, err := i.boardRepo.FindByID(tx, list.BoardID, list.UserID)
	if err != nil {
		return err
	}
	return nil
}

func sortLists(lists model.Lists) model.Lists {
	ls := map[string]model.List{}
	for _, l := range lists {
		ls[l.Before] = l
	}

	list := ls[""]
	r := model.Lists{}
	for range lists {
		r = append(r, list)
		list = ls[list.ID]
	}

	return r
}
