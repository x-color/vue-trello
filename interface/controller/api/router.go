package api

import (
	"errors"

	"github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"github.com/x-color/vue-trello/interface/controller/api/handler"
	"github.com/x-color/vue-trello/usecase"
)

// InteraBox includes all usecases interactors.
type InteraBox struct {
	item     usecase.ItemUsecase
	list     usecase.ListUsecase
	board    usecase.BoardUsecase
	user     usecase.UserUsecase
	resource usecase.ResourceUsecase
}

// NewInteraBox retruns new InteraBox.
func NewInteraBox(
	itemIntera usecase.ItemUsecase,
	listIntera usecase.ListUsecase,
	boardIntera usecase.BoardUsecase,
	userIntera usecase.UserUsecase,
	resourceIntera usecase.ResourceUsecase,
) (InteraBox, error) {
	if itemIntera == nil || listIntera == nil || boardIntera == nil || userIntera == nil || resourceIntera == nil {
		return InteraBox{}, errors.New("interactors are nil at least one")
	}
	b := InteraBox{
		item:     itemIntera,
		list:     listIntera,
		board:    boardIntera,
		user:     userIntera,
		resource: resourceIntera,
	}
	return b, nil
}

// NewRouter defines routing and returns echo instance.
func NewRouter(b InteraBox) *echo.Echo {
	userHandler := handler.NewUserHandler(b.user)
	itemHandler := handler.NewItemHandler(b.item)
	listHandler := handler.NewListHandler(b.list)
	boardHandler := handler.NewBoardHandler(b.board)
	resourceHandler := handler.NewResourceHandler(b.resource)

	e := echo.New()

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	e.Static("/", "web/public/")
	e.File("/", "web/public/index.html")

	e.POST("/signup", userHandler.SignUp)
	e.POST("/signin", userHandler.SignIn)

	api := e.Group("/api")
	api.Use(middleware.JWTWithConfig(middleware.JWTConfig{
		Claims:     &jwt.StandardClaims{},
		SigningKey: handler.SECRET,
	}))
	api.Use(middleware.CSRFWithConfig(middleware.CSRFConfig{
		TokenLookup: "header:X-XSRF-TOKEN",
	}))
	api.Use(checkContentType("application/json"))

	api.GET("/boards", boardHandler.GetBoards)
	api.GET("/boards/:id", boardHandler.Get)
	api.GET("/resources", resourceHandler.Get)

	api.DELETE("/items/:id", itemHandler.Delete)
	api.DELETE("/lists/:id", listHandler.Delete)
	api.DELETE("/boards/:id", boardHandler.Delete)

	api.POST("/items", itemHandler.Create)
	api.POST("/lists", listHandler.Create)
	api.POST("/boards", boardHandler.Create)

	api.PATCH("/items/:id", itemHandler.Update)
	api.PATCH("/lists/:id", listHandler.Update)
	api.PATCH("/boards/:id", boardHandler.Update)

	return e
}

func checkContentType(typ string) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			contentType := c.Request().Header.Get("Content-Type")
			if contentType != typ {
				return echo.ErrBadRequest
			}
			return next(c)
		}
	}
}
