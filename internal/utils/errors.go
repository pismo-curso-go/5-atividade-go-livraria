package utils

import "fmt"

type APIError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func (e APIError) Error() string {
	return fmt.Sprintf("API Error %d: %s", e.Code, e.Message)
}

func NewAPIError(code int, message string) *APIError {
	return &APIError{
		Code:    code,
		Message: message,
	}
}

var (
	ErrBookNotFound     = NewAPIError(404, "Book not found")
	ErrInvalidBookData  = NewAPIError(400, "Invalid book data")
	ErrInvalidJSON      = NewAPIError(400, "Invalid JSON")
	ErrMethodNotAllowed = NewAPIError(405, "Method not allowed")
	ErrInvalidID        = NewAPIError(400, "Invalid ID")
)
