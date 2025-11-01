package exception

import "net/http"

type AppError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func (e *AppError) Error() string {
	return e.Message
}

func NewAppError(code int, msg string) *AppError {
	return &AppError{Code: code, Message: msg}
}

// shortcut helpers
var (
	ErrUnauthorized  = func(msg string) *AppError { return NewAppError(http.StatusUnauthorized, msg) }
	ErrForbidden     = func(msg string) *AppError { return NewAppError(http.StatusForbidden, msg) }
	ErrNotFound      = func(msg string) *AppError { return NewAppError(http.StatusNotFound, msg) }
	ErrBadRequest    = func(msg string) *AppError { return NewAppError(http.StatusBadRequest, msg) }
	ErrUnprocessable = func(msg string) *AppError { return NewAppError(http.StatusUnprocessableEntity, msg) }
	ErrInternal      = func(msg string) *AppError { return NewAppError(http.StatusInternalServerError, msg) }
)
