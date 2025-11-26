package middleware

import (
	"brickbang/internal/utility"
	"log"
	"net/http"

	"github.com/gofiber/fiber/v2"
)

func ErrorHandler(c *fiber.Ctx, err error) error {
	if appErr, ok := err.(*utility.AppError); ok {
		return c.Status(appErr.Code).JSON(fiber.Map{
			"error":   true,
			"message": appErr.Message,
		})
	}
	if fiberErr, ok := err.(*fiber.Error); ok {
		return c.Status(fiberErr.Code).JSON(fiber.Map{
			"error":   true,
			"message": fiberErr.Message,
		})
	}
	log.Printf("[ERROR] %v\n", err)
	return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
		"error":   true,
		"message": "Internal Server Error",
	})
}
