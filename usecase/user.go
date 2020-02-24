package usecase

import (
	"errors"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/google/uuid"
	"github.com/x-color/vue-trello/model"
)

//SECRET uses to encode token for JWT.
var SECRET = []byte("secret")

// UserUsecase is interface. It defines to control a User authentication.
type UserUsecase interface {
	SignUp(user model.User) (model.User, error)
	Login(user model.User) (model.User, error)
}

// UserInteractor includes repogitories and a logger.
type UserInteractor struct {
	userRepo UserRepository
	logger   Logger
}

// NewUserInteractor generates new interactor for a User.
func NewUserInteractor(
	userRepo UserRepository,
	logger Logger,
) (UserInteractor, error) {
	i := UserInteractor{
		userRepo: userRepo,
		logger:   logger,
	}
	return i, nil
}

// SignUp registers new User to repository.
func (i *UserInteractor) SignUp(user model.User) (model.User, error) {
	if user.Name == "" || user.Password == "" {
		return model.User{}, model.InvalidContentError{
			Err: nil,
			ID:  "(No ID)",
			Act: "sigh up",
		}
	}

	u, err := i.userRepo.FindByName(model.User{Name: user.Name})
	if err != nil && !errors.Is(err, model.NotFoundError{}) {
		return model.User{}, err
	}
	if user.Name == u.Name {
		return model.User{}, model.ConflictError{
			Err: nil,
			ID:  user.Name,
			Act: "sign up",
		}
	}

	user.ID = uuid.New().String()
	if err := i.userRepo.Create(user); err != nil {
		return model.User{}, err
	}
	return user, nil
}

// SignIn returns JWT if a user succeds at authentication.
func (i *UserInteractor) SignIn(user model.User) (string, error) {
	u, err := i.userRepo.FindByName(model.User{Name: user.Name})
	if err != nil {
		return "", err
	}
	if user.Password != u.Password {
		return "", model.NotFoundError{
			Err: nil,
			ID:  u.ID,
			Act: "signin user",
		}
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.StandardClaims{
		Subject:   u.ID,
		ExpiresAt: time.Now().Add(time.Hour * 72).Unix(),
	})
	t, err := token.SignedString(SECRET)
	if err != nil {
		return "", model.ServerError{
			Err: err,
			ID:  u.ID,
			Act: "issue token",
		}
	}
	return t, nil
}
