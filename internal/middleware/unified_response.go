package middleware

import (
	"encoding/json"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

// UnifiedResponse provides a global middleware that wraps all responses (both success and error)
// into a unified JSON structure. It ensures a consistent format across all endpoints.
func UnifiedResponse() fiber.Handler {
	return func(c *fiber.Ctx) error {
		start := time.Now()

		// Execute the main request handler
		err := c.Next()

		// Skip wrapping if the "raw" flag is set in context
		if raw, ok := c.Locals("raw").(bool); ok && raw {
			return err
		}

		var (
			status  = c.Response().StatusCode()
			content any
		)

		// Handle errors and wrap them into a structured JSON response
		if err != nil {
			// Fiber has a built-in *fiber.Error type for HTTP-related errors
			if fe, ok := err.(*fiber.Error); ok {
				status = fe.Code
				content = fiber.Map{
					"error": fiber.Map{
						"message": fe.Message,
						"code":    fe.Code,
					},
				}
			} else {
				// Handle all other types of errors
				status = fiber.StatusInternalServerError
				content = fiber.Map{
					"error": fiber.Map{
						"message": err.Error(),
						"code":    status,
					},
				}
			}
		} else {
			// Parse the original response body (if present)
			body := c.Response().Body()
			if len(body) > 0 {
				_ = json.Unmarshal(body, &content)
			}
		}

		// Build the unified response structure
		response := fiber.Map{
			"content": content,
			"metadata": fiber.Map{
				"status":    status,
				"origin":    c.OriginalURL(),
				"request":   uuid.New().String(),
				"timeout":   time.Since(start).Milliseconds(),
				"timestamp": time.Now().UTC().Format(time.RFC3339),
			},
		}

		// Set response headers and return JSON
		c.Set(fiber.HeaderContentType, fiber.MIMEApplicationJSONCharsetUTF8)
		return c.Status(status).JSON(response)
	}
}
