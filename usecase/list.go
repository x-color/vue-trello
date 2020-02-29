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
			Err: nil,
			ID:  list.ID,
			Act: "validate list",
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

func (i *ListInteractor) validateList(list model.List) error {
	if list.ID == "" || list.Title == "" || list.BoardID == "" || list.UserID == "" {
		return model.InvalidContentError{
			Err: nil,
			ID:  list.ID,
			Act: "validate contents",
		}
	}
	_, err := i.boardRepo.Find(model.Board{ID: list.BoardID, UserID: list.UserID})
	if err != nil {
		return err
	}
	return nil
}
