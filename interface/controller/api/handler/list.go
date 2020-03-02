package handler

import (
	"net/http"

	"github.com/labstack/echo"
	"github.com/x-color/vue-trello/model"
	"github.com/x-color/vue-trello/usecase"
)

// List includes request data for List.
type List struct {
	ID      string `json:"id"`
	BoardID string `json:"board_id"`
	Title   string `json:"title"`
	Items   []Item `json:"items"`
	Before  string `json:"before"`
	After   string `json:"after"`
}

func (l *List) convertTo() model.List {
	list := model.List{
		ID:      l.ID,
		BoardID: l.BoardID,
		Title:   l.Title,
		Before:  l.Before,
		After:   l.After,
	}

	return list
}

func (l *List) convertFrom(list model.List) {
	l.ID = list.ID
	l.BoardID = list.BoardID
	l.Title = list.Title
	l.Before = list.Before
	l.After = list.After

	items := []Item{}
	for _, i := range list.Items {
		item := Item{}
		item.convertFrom(i)
		items = append(items, item)
	}
	l.Items = items
}

// ListHandler includes a interactor for List usecase.
type ListHandler struct {
	intractor usecase.ListUsecase
}

// NewListHandler returns a new ListHandler.
func NewListHandler(l usecase.ListUsecase) ListHandler {
	return ListHandler{
		intractor: l,
	}
}

// Create is http handler to create a list process.
func (h *ListHandler) Create(c echo.Context) error {
	reqList := new(List)
	if err := c.Bind(reqList); err != nil {
		return err
	}
	reqList.ID = c.Param("id")

	list := reqList.convertTo()
	list.UserID = getUserIDFromToken(c)

	l, err := h.intractor.Create(list)
	if err != nil {
		return convertToHTTPError(c, err)
	}

	resList := List{}
	resList.convertFrom(l)

	return c.JSON(http.StatusCreated, resList)
}

// Update is http handler to update a list process.
func (h *ListHandler) Update(c echo.Context) error {
	reqList := new(List)
	if err := c.Bind(reqList); err != nil {
		return err
	}
	reqList.ID = c.Param("id")

	list := reqList.convertTo()
	list.UserID = getUserIDFromToken(c)

	l, err := h.intractor.Update(list)
	if err != nil {
		return convertToHTTPError(c, err)
	}

	resList := List{}
	resList.convertFrom(l)

	return c.JSON(http.StatusOK, resList)
}

// Move is http handler to move a list process.
func (h *ListHandler) Move(c echo.Context) error {
	reqList := new(List)
	if err := c.Bind(reqList); err != nil {
		return err
	}
	reqList.ID = c.Param("id")

	list := reqList.convertTo()
	list.UserID = getUserIDFromToken(c)

	err := h.intractor.Move(list)
	if err != nil {
		return convertToHTTPError(c, err)
	}

	return c.NoContent(http.StatusNoContent)
}

// Delete is http handler to delete a list process.
func (h *ListHandler) Delete(c echo.Context) error {
	reqList := new(List)
	if err := c.Bind(reqList); err != nil {
		return err
	}
	reqList.ID = c.Param("id")

	list := reqList.convertTo()
	list.UserID = getUserIDFromToken(c)

	err := h.intractor.Delete(list)
	if err != nil {
		return convertToHTTPError(c, err)
	}

	return c.NoContent(http.StatusNoContent)
}
