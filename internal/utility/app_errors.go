package utility

import (
	"net/http"

	"github.com/gofiber/fiber/v2"
)

type AppError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func (rcv *AppError) Error() string {
	return rcv.Message
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
	ErrConflict      = func(msg string) *AppError { return NewAppError(http.StatusConflict, msg) }
	ErrUnprocessable = func(msg string) *AppError { return NewAppError(http.StatusUnprocessableEntity, msg) }
	ErrInternal      = func(msg string) *AppError { return NewAppError(http.StatusInternalServerError, msg) }
)

func RespondWithError(ctx *fiber.Ctx, status int, err error) error {
	return ctx.Status(status).JSON(fiber.Map{"error": err.Error()})
}
