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
	listRepo  ListRepository
	boardRepo BoardRepository
	logger    Logger
}

// NewListInteractor generates new interactor for a List.
func NewListInteractor(
	listRepo ListRepository,
	boardRepo BoardRepository,
	logger Logger,
) (ListInteractor, error) {
	i := ListInteractor{
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

	if err := i.listRepo.Create(list); err != nil {
		logError(i.logger, err)
		return model.List{}, err
	}
	i.logger.Info(formatLogMsg(list.UserID, "Create list("+list.ID+")"))
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
	if err := i.listRepo.Delete(list); err != nil {
		logError(i.logger, err)
		return err
	}
	i.logger.Info(formatLogMsg(list.UserID, "Delete list("+list.ID+")"))
	return nil
}

// Update replaces a List and returns new List.
func (i *ListInteractor) Update(list model.List) (model.List, error) {
	if err := i.validateList(list); err != nil {
		logError(i.logger, err)
		return model.List{}, err
	}

	if err := i.listRepo.Update(list); err != nil {
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

	if list.Before != "" {
		beforeList, err := i.listRepo.Find(model.List{ID: list.Before, UserID: list.UserID})
		if err != nil {
			logError(i.logger, err)
			return err
		}
		list.After = beforeList.After
	}

	if err := i.listRepo.Move(list); err != nil {
		logError(i.logger, err)
		return err
	}
	i.logger.Info(formatLogMsg(list.UserID, "Move list("+list.ID+") after list("+list.Before+") in list("+list.BoardID+")"))
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
	_, err := i.boardRepo.Find(model.Board{ID: list.BoardID, UserID: list.UserID})
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
