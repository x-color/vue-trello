package model

import "fmt"

// NotFoundError is occured if a content is not in a repository.
type NotFoundError struct {
	UserID string
	ID     string
	Act    string
	Err    error
}

func (e NotFoundError) Error() string {
	return fmt.Sprintf("%s NotFoundError: %s. %s does not exist", e.UserID, e.Act, e.ID)
}

// Unwrap returns a error wrapped by NotFoundError.
func (e NotFoundError) Unwrap() error {
	return e.Err
}

// Is checks target is NotFoundError.
func (e NotFoundError) Is(target error) bool {
	_, ok := target.(NotFoundError)
	return ok
}

// InvalidContentError is occured if a content is invalid.
type InvalidContentError struct {
	UserID string
	ID     string
	Act    string
	Err    error
}

func (e InvalidContentError) Error() string {
	return fmt.Sprintf("%s InvalidContentError: %s. %s has invalid content", e.UserID, e.Act, e.ID)
}

// Unwrap returns a error wrapped by InvalidContentError.
func (e InvalidContentError) Unwrap() error {
	return e.Err
}

// Is checks target is InvalidContentError.
func (e InvalidContentError) Is(target error) bool {
	_, ok := target.(InvalidContentError)
	return ok
}

// ConflictError is occured if a content already exists.
type ConflictError struct {
	UserID string
	ID     string
	Act    string
	Err    error
}

func (e ConflictError) Error() string {
	return fmt.Sprintf("%s ConflictError: %s. %s alredy exists", e.UserID, e.Act, e.ID)
}

// Unwrap returns a error wrapped by ConflictError.
func (e ConflictError) Unwrap() error {
	return e.Err
}

// Is checks target is ConflictError.
func (e ConflictError) Is(target error) bool {
	_, ok := target.(ConflictError)
	return ok
}

// ServerError is occured if a error is occured in server and not bad request.
type ServerError struct {
	UserID string
	ID     string
	Act    string
	Err    error
}

func (e ServerError) Error() string {
	return fmt.Sprintf("%s ServerError: %s. %s. internal server error is occured. %s", e.UserID, e.Act, e.ID, e.Err)
}

// Unwrap returns a error wrapped by ServerError.
func (e ServerError) Unwrap() error {
	return e.Err
}

// Is checks target is ServerError.
func (e ServerError) Is(target error) bool {
	_, ok := target.(ServerError)
	return ok
}
