package internal

import "errors"

// App Errors
var (
	ErrInvalidID          = errors.New("invalid ID")
	ErrMethodNotAllowed   = errors.New("method not allowed")
	ErrInvalidRequestBody = errors.New("invalid request body")
)

// Book Errors
var (
	ErrBookNotFound         = errors.New("book not found")
	ErrBookNotUpdated       = errors.New("book can not be updated")
	ErrBookStatusNotUpdated = errors.New("book status can not be updated")
	ErrBookCannotBeDeleted  = errors.New("book can not be deleted")
)
