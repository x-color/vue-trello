package handler

import (
	"net/http"

	"github.com/labstack/echo"
	"github.com/x-color/vue-trello/model"
	"github.com/x-color/vue-trello/usecase"
)

// Board includes request data for Board.
type Board struct {
	ID     string `json:"id"`
	Title  string `json:"title"`
	Text   string `json:"text"`
	Lists  []List `json:"lists"`
	Color  string `json:"color"`
	Before string `json:"before"`
	After  string `json:"after"`
}

func (b *Board) convertTo() model.Board {
	board := model.Board{
		ID:     b.ID,
		Title:  b.Title,
		Text:   b.Text,
		Color:  model.Color(b.Color),
		Before: b.Before,
		After:  b.After,
	}

	return board
}

func (b *Board) convertFrom(board model.Board) {
	b.ID = board.ID
	b.Title = board.Title
	b.Text = board.Text
	b.Color = string(board.Color)
	lists := []List{}
	for _, i := range board.Lists {
		list := List{}
		list.convertFrom(i)
		lists = append(lists, list)
	}
	b.Lists = lists
	b.Before = board.Before
	b.After = board.After
}

// BoardHandler includes a interactor for Board usecase.
type BoardHandler struct {
	intractor usecase.BoardUsecase
}

// NewBoardHandler returns a new BoardHandler.
func NewBoardHandler(i usecase.BoardUsecase) *BoardHandler {
	return &BoardHandler{
		intractor: i,
	}
}

// Create is http handler to create a board process.
func (h *BoardHandler) Create(c echo.Context) error {
	reqBoard := new(Board)
	if err := c.Bind(reqBoard); err != nil {
		return err
	}

	board := reqBoard.convertTo()
	board.UserID = getUserIDFromToken(c)

	b, err := h.intractor.Create(board)
	if err != nil {
		return convertToHTTPError(c, err)
	}

	resBoard := Board{}
	resBoard.convertFrom(b)

	return c.JSON(http.StatusCreated, resBoard)
}

// Update is http handler to update a board process.
func (h *BoardHandler) Update(c echo.Context) error {
	reqBoard := new(Board)
	if err := c.Bind(reqBoard); err != nil {
		return err
	}
	reqBoard.ID = c.Param("id")

	board := reqBoard.convertTo()
	board.UserID = getUserIDFromToken(c)

	b, err := h.intractor.Update(board)
	if err != nil {
		return convertToHTTPError(c, err)
	}

	resBoard := Board{}
	resBoard.convertFrom(b)

	return c.JSON(http.StatusOK, resBoard)
}

// Move is http handler to move a board process.
func (h *BoardHandler) Move(c echo.Context) error {
	reqBoard := new(Board)
	if err := c.Bind(reqBoard); err != nil {
		return err
	}
	reqBoard.ID = c.Param("id")

	board := reqBoard.convertTo()
	board.UserID = getUserIDFromToken(c)

	err := h.intractor.Move(board)
	if err != nil {
		return convertToHTTPError(c, err)
	}

	return c.NoContent(http.StatusNoContent)
}

// Delete is http handler to delete a board process.
func (h *BoardHandler) Delete(c echo.Context) error {
	reqBoard := new(Board)
	if err := c.Bind(reqBoard); err != nil {
		return err
	}
	reqBoard.ID = c.Param("id")

	board := reqBoard.convertTo()
	board.UserID = getUserIDFromToken(c)

	err := h.intractor.Delete(board)
	if err != nil {
		return convertToHTTPError(c, err)
	}

	return c.NoContent(http.StatusNoContent)
}

// Get is http handler to get user's board process.
func (h *BoardHandler) Get(c echo.Context) error {
	reqBoard := new(Board)
	reqBoard.ID = c.Param("id")

	board := reqBoard.convertTo()
	board.UserID = getUserIDFromToken(c)

	b, err := h.intractor.Get(board)
	if err != nil {
		return convertToHTTPError(c, err)
	}

	resBoard := Board{}
	resBoard.convertFrom(b)

	return c.JSON(http.StatusOK, resBoard)
}

// GetBoards is http handler to get user's boards process.
func (h *BoardHandler) GetBoards(c echo.Context) error {
	boards, err := h.intractor.GetBoards(model.User{ID: getUserIDFromToken(c)})
	if err != nil {
		return convertToHTTPError(c, err)
	}

	resBoards := []Board{}
	b := Board{}
	for _, board := range boards {
		b.convertFrom(board)
		resBoards = append(resBoards, b)
	}

	return c.JSON(http.StatusOK, map[string][]Board{
		"boards": resBoards},
	)
}
