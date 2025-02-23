package user_errors

import (
	"fmt"
	"net/http"
)

// Define custom error codes
const (
	ErrCodeUserNotFound   = 1001
	ErrCodeUsernameExists = 1002
	ErrCodeInternal       = 1003
)

// UserError structure with code, message, and optional context (cause)
type UserError struct {
	Code    int    // Error code
	Message string // Human-readable error message
	Err     error  // Optional wrapped error
}

// Implement the error interface for CustomError
func (e *UserError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("error %d: %s - %v", e.Code, e.Message, e.Err)
	}
	return fmt.Sprintf("error %d: %s", e.Code, e.Message)
}

// New creates a new error with the given code and message
func New(code int, message string) *UserError {
	return &UserError{
		Code:    code,
		Message: message,
	}
}

// Wrap allows you to add context to an existing error
func Wrap(err error, code int, message string) *UserError {
	return &UserError{
		Code:    code,
		Message: message,
		Err:     err,
	}
}

// Predefined errors using codes and messages
var (
	ErrUserNotFound   = New(ErrCodeUserNotFound, "user not found")
	ErrUsernameExists = New(ErrCodeUsernameExists, "username already exists")
	ErrInternal       = New(ErrCodeInternal, "internal server error")
)

func MapErrorCodeToHTTPStatus(code int) int {
	switch code {
	case ErrCodeUserNotFound:
		return http.StatusNotFound
	case ErrCodeUsernameExists:
		return http.StatusConflict
	case ErrCodeInternal:
		return http.StatusInternalServerError
	default:
		return http.StatusInternalServerError
	}
}
