package middleware

import (
	"fmt"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/rs/xid"

	"brickbang/internal"
)

var (
	ServerVersion = "Brickbang/" + internal.Version + "-" + internal.Staging
)

func GlobalHeaders() fiber.Handler {
	return func(c *fiber.Ctx) error {
		start := time.Now()
		reqID := xid.New().String()

		// Store request ID in context
		c.Locals("request_id", reqID)

		// Continue middleware chain
		err := c.Next()

		// Calc response time
		duration := time.Since(start).Milliseconds()

		// Add response headers globally
		c.Set("Server", ServerVersion)
		c.Set("X-Response-Time", fmt.Sprintf("%dms", duration))
		c.Set("X-Server-Timestamp", time.Now().UTC().Format(time.RFC3339))
		c.Set("Last-Modified", time.Now().UTC().Format(time.RFC1123))
		c.Set("X-Request-ID", reqID)

		return err
	}
}
