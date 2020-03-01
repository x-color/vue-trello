package handler

import (
	"net/http"

	"github.com/labstack/echo"
	"github.com/x-color/vue-trello/model"
	"github.com/x-color/vue-trello/usecase"
)

// Item includes request data for Item.
type Item struct {
	ID     string   `json:"id"`
	ListID string   `json:"list_id"`
	Title  string   `json:"title"`
	Text   string   `json:"text"`
	Tags   []string `json:"tags"`
	Before string   `json:"before"`
	After  string   `json:"after"`
}

func (i *Item) convertTo() model.Item {
	tags := model.Tags{}
	for _, tagID := range i.Tags {
		tags = append(tags, model.Tag{ID: tagID})
	}

	item := model.Item{
		ID:     i.ID,
		ListID: i.ListID,
		Title:  i.Title,
		Text:   i.Text,
		Tags:   tags,
		Before: i.Before,
		After:  i.After,
	}

	return item
}

func (i *Item) convertFrom(item model.Item) {
	tags := []string{}
	for _, tag := range item.Tags {
		tags = append(tags, tag.ID)
	}

	i.ID = item.ID
	i.ListID = item.ListID
	i.Title = item.Title
	i.Text = item.Text
	i.Tags = tags
	i.Before = item.Before
	i.After = item.After
}

// ItemHandler includes a interactor for Item usecase.
type ItemHandler struct {
	intractor usecase.ItemUsecase
}

// NewItemHandler returns a new ItemHandler.
func NewItemHandler(i usecase.ItemUsecase) ItemHandler {
	return ItemHandler{
		intractor: i,
	}
}

// Create is http handler to create a item process.
func (h *ItemHandler) Create(c echo.Context) error {
	reqItem := new(Item)
	if err := c.Bind(reqItem); err != nil {
		return err
	}

	item := reqItem.convertTo()
	item.UserID = getUserIDFromToken(c)

	i, err := h.intractor.Create(item)
	if err != nil {
		return convertToHTTPError(c, err)
	}

	resItem := Item{}
	resItem.convertFrom(i)

	return c.JSON(http.StatusCreated, resItem)
}

// Update is http handler to update a item process.
func (h *ItemHandler) Update(c echo.Context) error {
	reqItem := new(Item)
	if err := c.Bind(reqItem); err != nil {
		return err
	}
	reqItem.ID = c.Param("id")

	item := reqItem.convertTo()
	item.UserID = getUserIDFromToken(c)

	i, err := h.intractor.Update(item)
	if err != nil {
		return convertToHTTPError(c, err)
	}

	resItem := Item{}
	resItem.convertFrom(i)

	return c.JSON(http.StatusOK, resItem)
}

// Move is http handler to move a item process.
func (h *ItemHandler) Move(c echo.Context) error {
	reqItem := new(Item)
	if err := c.Bind(reqItem); err != nil {
		return err
	}
	reqItem.ID = c.Param("id")

	item := reqItem.convertTo()
	item.UserID = getUserIDFromToken(c)

	err := h.intractor.Move(item)
	if err != nil {
		return convertToHTTPError(c, err)
	}

	return c.NoContent(http.StatusNoContent)
}

// Delete is http handler to delete a item process.
func (h *ItemHandler) Delete(c echo.Context) error {
	reqItem := new(Item)
	if err := c.Bind(reqItem); err != nil {
		return err
	}
	reqItem.ID = c.Param("id")

	item := reqItem.convertTo()
	item.UserID = getUserIDFromToken(c)

	err := h.intractor.Delete(item)
	if err != nil {
		return convertToHTTPError(c, err)
	}

	return c.NoContent(http.StatusNoContent)
}
