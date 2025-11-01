package middleware

import "github.com/gofiber/fiber/v2"

func RBACMiddleware(c *fiber.Ctx) error {
    role := c.Get("X-Role")
    if role == "" {
        role = "guest"
    }
    c.Locals("role", role)
    return c.Next()
}
