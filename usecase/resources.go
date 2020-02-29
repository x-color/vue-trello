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
	tagRepo TagRepository
	logger  Logger
}

// NewResourceInteractor generates new interactor for resources.
func NewResourceInteractor(
	tagRepo TagRepository,
	logger Logger,
) (ResourceInteractor, error) {
	i := ResourceInteractor{
		tagRepo: tagRepo,
		logger:  logger,
	}
	return i, nil
}

// GetAllTagsandColors returns all Tags and Colors.
func (i *ResourceInteractor) GetAllTagsandColors() (model.Tags, model.Colors, error) {
	tags, err := i.tagRepo.FindAll()
	if err != nil {
		logError(i.logger, err)
		return model.Tags{}, model.Colors{}, err
	}
	i.logger.Info(formatLogMsg("(No-ID)", "Get resources"))
	return tags, model.COLORS, nil
}
