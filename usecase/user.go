package usecase

import (
	"errors"

	"github.com/google/uuid"
	"github.com/x-color/vue-trello/model"
)

// UserUsecase is interface. It defines to control a User authentication.
type UserUsecase interface {
	SignUp(user model.User) (model.User, error)
	SignIn(user model.User) (model.User, error)
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

// SignIn returns User data if a user succeds at authentication.
func (i *UserInteractor) SignIn(user model.User) (model.User, error) {
	u, err := i.userRepo.FindByName(model.User{Name: user.Name})
	if err != nil {
		return model.User{}, err
	}
	if user.Password != u.Password {
		return model.User{}, model.NotFoundError{
			Err: nil,
			ID:  u.ID,
			Act: "signin user",
		}
	}

	return u, nil
}
