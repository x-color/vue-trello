package usecase

import (
	"github.com/x-color/vue-trello/model"
)

// ResourceUsecase is interface. It defines getter for tags and colors.
type ResourceUsecase interface {
	GetAllTagsandColors() (model.Tags, model.Colors, error)
}

// ResourceInteractor includes repogitories and a logger.
type ResourceInteractor struct {
	txRepo  TransactionRepository
	tagRepo TagRepository
	logger  Logger
}

// NewResourceInteractor generates new interactor for resources.
func NewResourceInteractor(
	txRepo TransactionRepository,
	tagRepo TagRepository,
	logger Logger,
) (ResourceInteractor, error) {
	i := ResourceInteractor{
		txRepo:  txRepo,
		tagRepo: tagRepo,
		logger:  logger,
	}
	return i, nil
}

// GetAllTagsandColors returns all Tags and Colors.
func (i *ResourceInteractor) GetAllTagsandColors() (model.Tags, model.Colors, error) {
	tx := i.txRepo.BeginTransaction(false)
	tags, err := i.tagRepo.Find(tx, map[string]interface{}{})
	if err != nil {
		logError(i.logger, err)
		return model.Tags{}, model.Colors{}, err
	}
	i.logger.Info(formatLogMsg("(No-ID)", "Get resources"))
	return tags, model.COLORS, nil
}
