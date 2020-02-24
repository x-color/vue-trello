package handler

import (
	"errors"
	"net/http"

	"github.com/labstack/echo"
	"github.com/x-color/vue-trello/model"
)

func convertToHTTPError(c echo.Context, err error) error {
	switch {
	case errors.Is(err, model.ConflictError{}):
		return echo.NewHTTPError(http.StatusConflict, "resource already exists")
	case errors.Is(err, model.InvalidContentError{}):
		return echo.ErrBadRequest
	case errors.Is(err, model.NotFoundError{}):
		return echo.ErrBadRequest
	case errors.Is(err, model.ServerError{}):
		return echo.ErrInternalServerError
	default:
		return echo.ErrInternalServerError
	}
}
