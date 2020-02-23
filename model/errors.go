package model

import "fmt"

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

// ConflictError is occured if a content already exists.
type ConflictError struct {
	ID  string
	Act string
	Err error
}

func (e ConflictError) Error() string {
	return fmt.Sprintf("%s: %s alredy exists", e.Act, e.ID)
}

// Unwrap returns a error wrapped by ConflictError.
func (e ConflictError) Unwrap() error {
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
