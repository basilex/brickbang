package middleware

import (
	"slices"

	"github.com/gofiber/fiber/v2"
)

var publicRoutes = map[string][]string{
	"/api/v1/aux":           {fiber.MethodGet},
	"/api/v1/auth/login":    {fiber.MethodPost},
	"/api/v1/auth/register": {fiber.MethodPost},
}

func AuthMiddleware(c *fiber.Ctx) error {
	path := c.Path()
	method := c.Method()

	if methods, exists := publicRoutes[path]; exists && slices.Contains(methods, method) {
		return c.Next()
	}

	return c.Next()
}
