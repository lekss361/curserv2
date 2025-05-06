package errors

import "net/http"

type NotFoundError struct {
	Message string
}

func (e *NotFoundError) Error() string {
	return e.Message
}

func (e *NotFoundError) StatusCode() int {
	return http.StatusNotFound
}

func NewNotFoundError(msg string) *NotFoundError {
	return &NotFoundError{Message: msg}
}
