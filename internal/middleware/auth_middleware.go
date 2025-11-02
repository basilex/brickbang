package middleware

import (
	"strings"

	"github.com/gofiber/fiber/v2"
)

func AuthMiddleware(c *fiber.Ctx) error {
	path := c.Path()
	method := c.Method()

	// public routes
	if strings.HasPrefix(path, "/api/v1/aux") && method == fiber.MethodGet {
		return c.Next()
	}
	if strings.HasPrefix(path, "/api/v1/auth") && method == fiber.MethodGet {
		return c.Next()
	}

	// API key authentication
	apiKey := c.Get("X-API-Key")
	if apiKey == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "missing API key",
		})
	}
	// TODO: validate apiKey here
	return c.Next()
}
