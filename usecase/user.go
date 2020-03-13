package usecase

import (
	"errors"

	"github.com/google/uuid"
	"github.com/x-color/vue-trello/model"
	"golang.org/x/crypto/bcrypt"
)

// UserUsecase is interface. It defines to control a User authentication.
type UserUsecase interface {
	SignUp(user model.User) (model.User, error)
	SignIn(user model.User) (model.User, error)
}

// UserInteractor includes repogitories and a logger.
type UserInteractor struct {
	txRepo   TransactionRepository
	userRepo UserRepository
	logger   Logger
}

// NewUserInteractor generates new interactor for a User.
func NewUserInteractor(
	txRepo TransactionRepository,
	userRepo UserRepository,
	logger Logger,
) (UserInteractor, error) {
	i := UserInteractor{
		txRepo:   txRepo,
		userRepo: userRepo,
		logger:   logger,
	}
	return i, nil
}

// SignUp registers new User to repository.
func (i *UserInteractor) SignUp(user model.User) (model.User, error) {
	tx := i.txRepo.BeginTransaction(false)

	u, err := i.userRepo.Find(tx, map[string]interface{}{
		"Name": user.Name,
	})
	if err != nil && !errors.Is(err, model.NotFoundError{}) {
		logError(i.logger, err)
		return model.User{}, err
	}
	if user.Name == u.Name {
		i.logger.Info(formatLogMsg(user.ID, "New user name conflicts. '"+user.Name+"' already exists"))
		return model.User{}, model.ConflictError{
			UserID: "(No-ID)",
			Err:    nil,
			ID:     user.Name,
			Act:    "validate name",
		}
	}

	user.ID = uuid.New().String()
	p, err := hashPassword(user.Password)
	if err != nil {
		i.logger.Info(formatLogMsg(user.ID, err.Error()))
		return model.User{}, model.InvalidContentError{
			UserID: "(No-ID)",
			Err:    err,
			ID:     user.Name,
			Act:    "hash password",
		}
	}
	user.Password = p
	if err := i.userRepo.Create(tx, user); err != nil {
		logError(i.logger, err)
		return model.User{}, err
	}
	i.logger.Info(formatLogMsg(user.ID, "Create user("+user.ID+")"))
	return user, nil
}

// SignIn returns User data if a user succeds at authentication.
func (i *UserInteractor) SignIn(user model.User) (model.User, error) {
	tx := i.txRepo.BeginTransaction(false)
	u, err := i.userRepo.Find(tx, map[string]interface{}{
		"Name": user.Name,
	})
	if err != nil {
		logError(i.logger, err)
		return model.User{}, err
	}

	if err := comparePassword(user.Password, u.Password); err != nil {
		i.logger.Info(formatLogMsg(u.ID, "Invalid password. '"+u.ID+"' Fails to sign in"))
		return model.User{}, model.NotFoundError{
			UserID: u.ID,
			Err:    err,
			ID:     u.ID,
			Act:    "validate password",
		}
	}

	i.logger.Info(formatLogMsg(u.ID, "Sign in user("+u.ID+")"))
	return u, nil
}

func hashPassword(password string) (string, error) {
	if len(password) > 72 {
		return "", errors.New("password length is too long")
	}
	h, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(h), err
}

func comparePassword(password string, hash string) error {
	if len(password) > 72 {
		return errors.New("password length is too long")
	}
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
}
