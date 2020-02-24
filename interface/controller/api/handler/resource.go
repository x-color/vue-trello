package handler

import (
	"net/http"

	"github.com/labstack/echo"
	"github.com/x-color/vue-trello/model"
	"github.com/x-color/vue-trello/usecase"
)

// Resources includes request data for Resources.
type Resources struct {
	Colors []string `json:"colors"`
	Tags   []string `json:"tags"`
}

func (r *Resources) convertFrom(tags model.Tags, colors model.Colors) {
	ts := []string{}
	for _, tag := range tags {
		ts = append(ts, tag.ID)
	}
	r.Tags = ts

	cs := []string{}
	for _, color := range colors {
		cs = append(cs, string(color))
	}

	r.Colors = cs
}

// ResourceHandler includes a interactor for Resources usecase.
type ResourceHandler struct {
	intractor usecase.ResourceUsecase
}

// NewResourceHandler returns a new ResourceHandler.
func NewResourceHandler(r usecase.ResourceUsecase) ResourceHandler {
	return ResourceHandler{
		intractor: r,
	}
}

// Get is http handler to get resources process.
func (h *ResourceHandler) Get(c echo.Context) error {
	tags, colors, err := h.intractor.GetAllTagsandColors()
	if err != nil {
		return convertToHTTPError(c, err)
	}

	resResources := Resources{}
	resResources.convertFrom(tags, colors)

	return c.JSON(http.StatusOK, resResources)
}
