package middleware

import (
	"log/slog"
	"time"

	"github.com/gofiber/fiber/v2"
)

func RequestLogger(logger *slog.Logger) fiber.Handler {
	return func(c *fiber.Ctx) error {
		start := time.Now()

		err := c.Next()

		latency := time.Since(start)
		status := c.Response().StatusCode()
		if status == 0 {
			status = 200 // Default to 200 if not set
		}

		log := logger.With(
			"method", c.Method(),
			"path", c.Path(),
			"status", status,
			"latency", latency.String(),
			"ip", c.IP(),
			"user_agent", string(c.Context().UserAgent()),
		)

		if err != nil {
			log.Error("error", "error", err)
			return err
		}

		log.Info("done")
		return nil
	}
}
