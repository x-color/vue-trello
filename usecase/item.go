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
}

// ItemInteractor includes repogitories and a logger.
type ItemInteractor struct {
	itemRepo ItemRepogitory
	listRepo ListRepogitory
	tagRepo  TagRepogitory
	logger   Logger
}

// NewItemInteractor generates new interactor for a Item.
func NewItemInteractor(
	ItemRepo ItemRepogitory,
	ListRepo ListRepogitory,
	TagRepo TagRepogitory,
	logger Logger,
) (ItemInteractor, error) {
	i := ItemInteractor{
		itemRepo: ItemRepo,
		listRepo: ListRepo,
		tagRepo:  TagRepo,
		logger:   logger,
	}
	return i, nil
}

// Create saves new Item to a repogitory and returns created Item.
func (i *ItemInteractor) Create(item model.Item) (model.Item, error) {
	item.ID = uuid.New().String()
	if err := i.validateItem(item); err != nil {
		return model.Item{}, err
	}

	if err := i.itemRepo.Create(item); err != nil {
		return model.Item{}, err
	}
	return item, nil
}

// Delete removes item in repogitory.
func (i *ItemInteractor) Delete(item model.Item) error {
	if item.ID == "" {
		return new(ErrGen).InvalidContentError(nil, item.ID, "delete item")
	}
	if err := i.itemRepo.Delete(item); err != nil {
		return err
	}
	return nil
}

// Update replaces item and returns new item.
func (i *ItemInteractor) Update(item model.Item) (model.Item, error) {
	if err := i.validateItem(item); err != nil {
		return model.Item{}, nil
	}

	if err := i.itemRepo.Update(item); err != nil {
		return model.Item{}, err
	}
	return item, nil
}

func (i *ItemInteractor) validateItem(item model.Item) error {
	if item.ID == "" || item.Title == "" || item.ListID == "" || item.UserID == "" {
		return new(ErrGen).InvalidContentError(nil, item.ID, "validate contents")
	}
	_, err := i.listRepo.Find(model.List{ID: item.ListID})
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
			return new(ErrGen).InvalidContentError(nil, item.ID, "validate tags")
		}
	}
	return nil
}
