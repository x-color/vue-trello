package usecase

import (
	"github.com/google/uuid"
	"github.com/x-color/vue-trello/model"
)

// ItemUsecase is interface. It defines to control a Item.
type ItemUsecase interface {
	Create(item model.Item) (model.Item, error)
	Delete(item model.Item) error
	Update(item model.Item) (model.Item, error)
	Move(item model.Item) error
}

// ItemInteractor includes repogitories and a logger.
type ItemInteractor struct {
	itemRepo ItemRepository
	listRepo ListRepository
	tagRepo  TagRepository
	logger   Logger
}

// NewItemInteractor generates new interactor for a Item.
func NewItemInteractor(
	itemRepo ItemRepository,
	listRepo ListRepository,
	tagRepo TagRepository,
	logger Logger,
) (ItemInteractor, error) {
	i := ItemInteractor{
		itemRepo: itemRepo,
		listRepo: listRepo,
		tagRepo:  tagRepo,
		logger:   logger,
	}
	return i, nil
}

// Create saves new Item to a repository and returns created Item.
func (i *ItemInteractor) Create(item model.Item) (model.Item, error) {
	item.ID = uuid.New().String()
	if err := i.validateItem(item); err != nil {
		logError(i.logger, err)
		return model.Item{}, err
	}

	if err := i.itemRepo.Create(item); err != nil {
		logError(i.logger, err)
		return model.Item{}, err
	}
	i.logger.Info(formatLogMsg(item.UserID, "Create item("+item.ID+")"))
	return item, nil
}

// Delete removes Item in repository.
func (i *ItemInteractor) Delete(item model.Item) error {
	if item.ID == "" {
		i.logger.Info(formatLogMsg(item.UserID, "Invalid item. ID is empty"))
		return model.InvalidContentError{
			UserID: item.UserID,
			Err:    nil,
			ID:     "(No-ID)",
			Act:    "validate item id",
		}
	}
	if err := i.itemRepo.Delete(item); err != nil {
		logError(i.logger, err)
		return err
	}
	i.logger.Info(formatLogMsg(item.UserID, "Delete item("+item.ID+")"))
	return nil
}

// Update replaces a Item and returns new Item.
func (i *ItemInteractor) Update(item model.Item) (model.Item, error) {
	if err := i.validateItem(item); err != nil {
		logError(i.logger, err)
		return model.Item{}, err
	}

	if err := i.itemRepo.Update(item); err != nil {
		logError(i.logger, err)
		return model.Item{}, err
	}
	i.logger.Info(formatLogMsg(item.UserID, "Update item("+item.ID+")"))
	return item, nil
}

// Move moves Items.
func (i *ItemInteractor) Move(item model.Item) error {
	item.Title = "dummy title"
	if err := i.validateItem(item); err != nil {
		logError(i.logger, err)
		return err
	}
	item.Title = ""

	if item.Before != "" {
		beforeItem, err := i.itemRepo.Find(model.Item{ID: item.Before, UserID: item.UserID})
		if err != nil {
			logError(i.logger, err)
			return err
		}
		item.After = beforeItem.After
	}

	if err := i.itemRepo.Move(item); err != nil {
		logError(i.logger, err)
		return err
	}
	i.logger.Info(formatLogMsg(item.UserID, "Move item("+item.ID+") after item("+item.Before+") in list("+item.ListID+")"))
	return nil
}

func (i *ItemInteractor) validateItem(item model.Item) error {
	if item.ID == "" || item.Title == "" || item.ListID == "" || item.UserID == "" {
		return model.InvalidContentError{
			UserID: item.UserID,
			Err:    nil,
			ID:     item.ID,
			Act:    "validate contents in item",
		}
	}
	_, err := i.listRepo.Find(model.List{ID: item.ListID, UserID: item.UserID})
	if err != nil {
		return err
	}

	allTags, err := i.tagRepo.FindAll()
	if err != nil {
		return err
	}

	// Validate tags attached to item
	for _, tag := range item.Tags {
		isValid := false
		for _, t := range allTags {
			if t.ID == tag.ID {
				isValid = true
				break
			}
		}
		if !isValid {
			return model.InvalidContentError{
				UserID: item.UserID,
				Err:    nil,
				ID:     item.ID,
				Act:    "validate tags of item",
			}
		}
	}
	return nil
}

func sortItems(items model.Items) model.Items {
	l := map[string]model.Item{}
	for _, i := range items {
		l[i.Before] = i
	}

	item := l[""]
	r := model.Items{}
	for range items {
		r = append(r, item)
		item = l[item.ID]
	}

	return r
}
