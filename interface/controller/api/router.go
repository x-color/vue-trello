package api

import (
	"errors"
	"net/http"

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

	echo.NotFoundHandler = func(c echo.Context) error {
		return c.Redirect(http.StatusMovedPermanently, "/?redirect="+c.Request().URL.Path)
	}

	e := echo.New()

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	e.Static("/", "web/dist")
	e.File("/", "web/dist/index.html")

	auth := e.Group("/auth")
	auth.POST("/signup", userHandler.SignUp)
	auth.POST("/signin", userHandler.SignIn)
	auth.GET("/signout", userHandler.SignOut)

	api := e.Group("/api")
	api.Use(middleware.JWTWithConfig(middleware.JWTConfig{
		Claims:      &jwt.StandardClaims{},
		SigningKey:  handler.SECRET,
		TokenLookup: "cookie:token",
	}))
	api.Use(checkContentType("application/json; charset=UTF-8"))
	api.Use(checkCSRFToken())

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

	api.PATCH("/items/:id/move", itemHandler.Move)
	api.PATCH("/lists/:id/move", listHandler.Move)
	api.PATCH("/boards/:id/move", boardHandler.Move)

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

func checkCSRFToken() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			csrf := c.Request().Header.Get("X-XSRF-TOKEN")
			if csrf == "" {
				return echo.ErrForbidden
			}
			return next(c)
		}
	}
}
