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
	txRepo   TransactionRepository
	itemRepo ItemRepository
	listRepo ListRepository
	tagRepo  TagRepository
	logger   Logger
}

// NewItemInteractor generates new interactor for a Item.
func NewItemInteractor(
	txRepo TransactionRepository,
	itemRepo ItemRepository,
	listRepo ListRepository,
	tagRepo TagRepository,
	logger Logger,
) (ItemInteractor, error) {
	i := ItemInteractor{
		txRepo:   txRepo,
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

	tx := i.txRepo.BeginTransaction(true)
	i.logger.Info(formatLogMsg(item.UserID, "Start transaction"))

	// Get last item in list
	items, err := i.itemRepo.Find(tx, map[string]interface{}{
		"ListID": item.ListID,
		"UserID": item.UserID,
		"After":  "",
	})
	if err != nil {
		tx.Rollback()
		i.logger.Info(formatLogMsg(item.UserID, "Rollback transaction"))
		logError(i.logger, err)
		return model.Item{}, err
	}

	if len(items) > 0 {
		lastItem := items[0]
		i.logger.Info(formatLogMsg(item.UserID, "Find last item("+lastItem.ID+") in list("+lastItem.ListID+")"))

		item.Before = lastItem.ID

		query := map[string]interface{}{
			"After": item.ID,
		}
		if err := i.itemRepo.Update(tx, lastItem, query); err != nil {
			tx.Rollback()
			i.logger.Info(formatLogMsg(item.UserID, "Rollback transaction"))
			logError(i.logger, err)
			return model.Item{}, err
		}
		i.logger.Info(formatLogMsg(item.UserID, "Update item("+lastItem.ID+")"))
	} else {
		i.logger.Info(formatLogMsg(item.UserID, "Find no item in list("+item.ListID+")"))
	}

	if err := i.itemRepo.Create(tx, item); err != nil {
		tx.Rollback()
		i.logger.Info(formatLogMsg(item.UserID, "Rollback transaction"))
		logError(i.logger, err)
		return model.Item{}, err
	}
	i.logger.Info(formatLogMsg(item.UserID, "Create item("+item.ID+")"))

	tx.Commit()
	i.logger.Info(formatLogMsg(item.UserID, "Commit transaction"))

	return item, nil
}

// Delete removes Item in repository.
func (i *ItemInteractor) Delete(item model.Item) error {
	if item.ID == "" || item.UserID == "" {
		i.logger.Info(formatLogMsg(item.UserID, "Invalid item. ID is empty"))
		return model.InvalidContentError{
			UserID: item.UserID,
			Err:    nil,
			ID:     "(No-ID)",
			Act:    "validate item id",
		}
	}

	tx := i.txRepo.BeginTransaction(true)
	i.logger.Info(formatLogMsg(item.UserID, "Start transaction"))

	// Get item's info (e.g. item.Before, item.After...) and rewrite 'item'.
	item, err := i.itemRepo.FindByID(tx, item.ID, item.UserID)
	if err != nil {
		tx.Rollback()
		i.logger.Info(formatLogMsg(item.UserID, "Rollback transaction"))
		logError(i.logger, err)
		return err
	}
	i.logger.Info(formatLogMsg(item.UserID, "Find item("+item.ID+")"))

	if item.Before != "" {
		// Update item before deleting item
		before := model.Item{
			ID:     item.Before,
			UserID: item.UserID,
		}
		query := map[string]interface{}{
			"After": item.After,
		}
		if err := i.itemRepo.Update(tx, before, query); err != nil {
			tx.Rollback()
			i.logger.Info(formatLogMsg(item.UserID, "Rollback transaction"))
			logError(i.logger, err)
			return err
		}
		i.logger.Info(formatLogMsg(item.UserID, "Update item("+before.ID+") before deleted item("+item.ID+")"))
	}

	if item.After != "" {
		// Update item after deleting item
		after := model.Item{
			ID:     item.After,
			UserID: item.UserID,
		}
		query := map[string]interface{}{
			"Before": item.Before,
		}
		if err := i.itemRepo.Update(tx, after, query); err != nil {
			tx.Rollback()
			i.logger.Info(formatLogMsg(item.UserID, "Rollback transaction"))
			logError(i.logger, err)
			return err
		}
		i.logger.Info(formatLogMsg(item.UserID, "Update item("+after.ID+") after deleted item("+item.ID+")"))
	}

	if err := i.itemRepo.Delete(tx, item); err != nil {
		tx.Rollback()
		i.logger.Info(formatLogMsg(item.UserID, "Rollback transaction"))
		logError(i.logger, err)
		return err
	}
	i.logger.Info(formatLogMsg(item.UserID, "Delete item("+item.ID+")"))

	tx.Commit()
	i.logger.Info(formatLogMsg(item.UserID, "Commit transaction"))

	return nil
}

// Update replaces a Item and returns new Item.
func (i *ItemInteractor) Update(item model.Item) (model.Item, error) {
	if err := i.validateItem(item); err != nil {
		logError(i.logger, err)
		return model.Item{}, err
	}

	tags := []string{}
	for _, t := range item.Tags {
		tags = append(tags, t.ID)
	}

	query := map[string]interface{}{
		"Title": item.Title,
		"Text":  item.Text,
		"Tags":  tags,
	}

	tx := i.txRepo.BeginTransaction(false)
	if err := i.itemRepo.Update(tx, item, query); err != nil {
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

	tx := i.txRepo.BeginTransaction(true)
	i.logger.Info(formatLogMsg(item.UserID, "Start transaction"))

	// Get a item to move
	if item.Before != "" {
		before, err := i.itemRepo.FindByID(tx, item.Before, item.UserID)
		if err != nil {
			tx.Rollback()
			i.logger.Info(formatLogMsg(item.UserID, "Rollback transaction"))
			logError(i.logger, err)
			return err
		}
		i.logger.Info(formatLogMsg(item.UserID, "Find a item("+before.ID+") before item to move"))
		if item.ListID != before.ListID {
			tx.Rollback()
			i.logger.Info(formatLogMsg(item.UserID, "Rollback transaction"))
			i.logger.Info(formatLogMsg(item.UserID, "Invalid list id("+item.ListID+"). Does not equal before item's list id("+before.ListID+")"))
			return err
		}
	}

	// Get a item to move
	old, err := i.itemRepo.FindByID(tx, item.ID, item.UserID)
	if err != nil {
		tx.Rollback()
		i.logger.Info(formatLogMsg(item.UserID, "Rollback transaction"))
		logError(i.logger, err)
		return err
	}
	i.logger.Info(formatLogMsg(item.UserID, "Find item("+item.ID+") to move"))

	// Get a item before item to move
	if old.Before != "" {
		beforeOld := model.Item{
			ID:     old.Before,
			UserID: old.UserID,
		}
		query := map[string]interface{}{
			"After": old.After,
		}
		if err := i.itemRepo.Update(tx, beforeOld, query); err != nil {
			tx.Rollback()
			i.logger.Info(formatLogMsg(item.UserID, "Rollback transaction"))
			logError(i.logger, err)
			return err
		}
		i.logger.Info(formatLogMsg(item.UserID, "Find a item("+beforeOld.ID+") before item("+item.ID+") to move"))
	}

	// Get a item after item to move
	if old.After != "" {
		afterOld := model.Item{
			ID:     old.After,
			UserID: old.UserID,
		}
		query := map[string]interface{}{
			"Before": old.Before,
		}
		if err := i.itemRepo.Update(tx, afterOld, query); err != nil {
			tx.Rollback()
			i.logger.Info(formatLogMsg(item.UserID, "Rollback transaction"))
			logError(i.logger, err)
			return err
		}
		i.logger.Info(formatLogMsg(item.UserID, "Find a item("+afterOld.ID+") after item("+item.ID+") to move"))
	}

	// Get a item after moved item
	if item.Before == "" {
		conditions := map[string]interface{}{
			"ListID": item.ListID,
			"UserID": item.UserID,
			"Before": "",
		}
		l, err := i.itemRepo.Find(tx, conditions)
		if err != nil {
			tx.Rollback()
			i.logger.Info(formatLogMsg(item.UserID, "Rollback transaction"))
			logError(i.logger, err)
			return err
		}
		if len(l) == 1 {
			item.After = l[0].ID
		} else if len(l) > 1 {
			tx.Rollback()
			i.logger.Info(formatLogMsg(item.UserID, "Rollback transaction"))
			err = model.ServerError{
				ID:     item.ID,
				UserID: item.UserID,
				Err:    nil,
				Act:    "find a item before moved item",
			}
			logError(i.logger, err)
			return err
		}
	} else {
		before, err := i.itemRepo.FindByID(tx, item.Before, item.UserID)
		if err != nil {
			tx.Rollback()
			i.logger.Info(formatLogMsg(item.UserID, "Rollback transaction"))
			logError(i.logger, err)
			return err
		}
		item.After = before.After
	}
	i.logger.Info(formatLogMsg(item.UserID, "Find a item("+item.After+") after moved item("+item.ID+")"))

	// Update a item before moved item
	if item.Before != "" {
		before := model.Item{
			ID:     item.Before,
			UserID: item.UserID,
		}
		query := map[string]interface{}{
			"After": item.ID,
		}
		if err := i.itemRepo.Update(tx, before, query); err != nil {
			tx.Rollback()
			i.logger.Info(formatLogMsg(item.UserID, "Rollback transaction"))
			logError(i.logger, err)
			return err
		}
		i.logger.Info(formatLogMsg(item.UserID, "Update item("+before.ID+") before item("+item.ID+") in list("+item.ListID+")"))
	}

	// Update a item after moved item
	if item.After != "" {
		after := model.Item{
			ID:     item.After,
			UserID: item.UserID,
		}
		query := map[string]interface{}{
			"Before": item.ID,
		}
		if err := i.itemRepo.Update(tx, after, query); err != nil {
			tx.Rollback()
			i.logger.Info(formatLogMsg(item.UserID, "Rollback transaction"))
			logError(i.logger, err)
			return err
		}
		i.logger.Info(formatLogMsg(item.UserID, "Update item("+after.ID+") after item("+item.ID+") in list("+item.ListID+")"))
	}

	// Move item
	query := map[string]interface{}{
		"ListID": item.ListID,
		"Before": item.Before,
		"After":  item.After,
	}
	if err := i.itemRepo.Update(tx, item, query); err != nil {
		tx.Rollback()
		i.logger.Info(formatLogMsg(item.UserID, "Rollback transaction"))
		logError(i.logger, err)
		return err
	}
	i.logger.Info(formatLogMsg(item.UserID, "Move item("+item.ID+") after item("+item.Before+") in list("+item.ListID+")"))

	tx.Commit()
	i.logger.Info(formatLogMsg(item.UserID, "Commit transaction"))

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

	tx := i.txRepo.BeginTransaction(false)
	_, err := i.listRepo.FindByID(tx, item.ListID, item.UserID)
	if err != nil {
		return err
	}

	allTags, err := i.tagRepo.Find(tx, map[string]interface{}{})
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
