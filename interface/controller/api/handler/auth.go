package handler

import (
	"errors"
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo"
	"github.com/x-color/vue-trello/model"
	"github.com/x-color/vue-trello/usecase"
)

//SECRET uses to encode token for JWT.
var SECRET = []byte("secret")

// User includes request data for authentication.
type User struct {
	Name     string `json:"name"`
	Password string `json:"password"`
}

func (u *User) convertTo() model.User {
	user := model.User{
		Name:     u.Name,
		Password: u.Password,
	}

	return user
}

func (u *User) convertFrom(user model.User) {
	u.Name = user.Name
	u.Password = ""
}

// UserHandler includes a interactor for user usecase.
type UserHandler struct {
	interactor usecase.UserUsecase
}

// NewUserHandler returns a new UserHandler.
func NewUserHandler(u usecase.UserUsecase) UserHandler {
	return UserHandler{
		interactor: u,
	}
}

// SignUp is http handler to signup process.
func (h *UserHandler) SignUp(c echo.Context) error {
	user := User{}
	if err := c.Bind(&user); err != nil {
		return err
	}

	if user.Name == "" || user.Password == "" {
		return echo.ErrBadRequest
	}

	u, err := h.interactor.SignUp(user.convertTo())
	if err != nil {
		return convertToHTTPError(c, err)
	}

	r := User{}
	r.convertFrom(u)
	return c.JSON(http.StatusCreated, r)
}

// SignIn is http handler to sign in process.
func (h *UserHandler) SignIn(c echo.Context) error {
	user := User{}
	if err := c.Bind(&user); err != nil {
		return err
	}

	u, err := h.interactor.SignIn(user.convertTo())
	if err != nil {
		if errors.Is(err, model.NotFoundError{}) {
			return echo.ErrUnauthorized
		}
		return convertToHTTPError(c, err)
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.StandardClaims{
		Subject:   u.ID,
		ExpiresAt: time.Now().Add(time.Hour * 72).Unix(),
	})

	t, err := token.SignedString(SECRET)
	if err != nil {
		return echo.ErrInternalServerError
	}

	cookie := new(http.Cookie)
	cookie.Name = "token"
	cookie.Value = t
	cookie.HttpOnly = true
	// NOTE: It should activate Secure attribute of Cookie.
	//		 But this code is sample, it does not set this attribute and TLS.
	// cookie.Secure = true
	cookie.SameSite = http.SameSiteStrictMode
	cookie.Expires = time.Now().Add(72 * time.Hour)
	cookie.Path = "/"
	c.SetCookie(cookie)

	return c.JSON(http.StatusOK, map[string]string{
		"message": "Sign in",
	})
}

// SignOut is http handler to sign out process.
func (h *UserHandler) SignOut(c echo.Context) error {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.StandardClaims{
		ExpiresAt: time.Now().Unix(),
	})

	t, err := token.SignedString(SECRET)
	if err != nil {
		return echo.ErrInternalServerError
	}

	cookie := new(http.Cookie)
	cookie.Name = "token"
	cookie.Value = t
	cookie.HttpOnly = true
	// NOTE: It should activate Secure attribute of Cookie.
	//		 But this code is sample, it does not set this attribute and TLS.
	// cookie.Secure = true
	cookie.SameSite = http.SameSiteStrictMode
	cookie.Path = "/"
	c.SetCookie(cookie)

	return c.JSON(http.StatusOK, map[string]string{
		"message": "Sign out",
	})
}

func getUserIDFromToken(c echo.Context) string {
	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(*jwt.StandardClaims)
	uid := claims.Subject
	return uid
}
