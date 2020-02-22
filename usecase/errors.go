package usecase

import "fmt"

// ErrorGenerator is interface. It includes error generators.
type ErrorGenerator interface {
	NotFoundError(err error, id string, act string) NotFoundError
	InvalidContentError(err error, id string, act string) InvalidContentError
	ServerError(err error, id string, act string) ServerError
}

// ErrGen defines error generators.
type ErrGen struct {
}

// NotFoundError generates NotFoundError object.
func (*ErrGen) NotFoundError(err error, id string, act string) NotFoundError {
	return NotFoundError{
		ID:  id,
		Act: act,
		Err: err,
	}
}

// InvalidContentError generates InvalidContentError object.
func (*ErrGen) InvalidContentError(err error, id string, act string) InvalidContentError {
	return InvalidContentError{
		ID:  id,
		Act: act,
		Err: err,
	}
}

// ServerError generates ServerError object.
func (*ErrGen) ServerError(err error, id string, act string) ServerError {
	return ServerError{
		ID:  id,
		Act: act,
		Err: err,
	}
}

// NotFoundError is occured if a content is not in a repogitory.
type NotFoundError struct {
	ID  string
	Act string
	Err error
}

func (e NotFoundError) Error() string {
	return fmt.Sprintf("%s: %s does not exist", e.Act, e.ID)
}

// Unwrap returns a error wrapped by NotFoundError.
func (e NotFoundError) Unwrap() error {
	return e.Err
}

// InvalidContentError is occured if a content is invalid.
type InvalidContentError struct {
	ID  string
	Act string
	Err error
}

func (e InvalidContentError) Error() string {
	return fmt.Sprintf("%s: %s has invalid content", e.Act, e.ID)
}

// Unwrap returns a error wrapped by InvalidContentError.
func (e InvalidContentError) Unwrap() error {
	return e.Err
}

// ServerError is occured if a error is occured in server and not bad request.
type ServerError struct {
	ID  string
	Act string
	Err error
}

func (e ServerError) Error() string {
	return fmt.Sprintf("%s: %s. internal server error is occured", e.Act, e.ID)
}

// Unwrap returns a error wrapped by ServerError.
func (e ServerError) Unwrap() error {
	return e.Err
}
