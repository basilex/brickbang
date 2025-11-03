package middleware

import (
	"encoding/json"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

func UnifiedResponse() fiber.Handler {
	return func(c *fiber.Ctx) error {
		start := time.Now()

		// Выполняем основной обработчик
		err := c.Next()

		// Если был установлен "raw" — не оборачиваем
		if raw, ok := c.Locals("raw").(bool); ok && raw {
			return err
		}

		var (
			status  = c.Response().StatusCode()
			content any
		)

		// Если произошла ошибка — готовим ответ в виде content.error
		if err != nil {
			// Fiber уже содержит тип *fiber.Error для HTTP ошибок
			if fe, ok := err.(*fiber.Error); ok {
				status = fe.Code
				content = fiber.Map{
					"error": fiber.Map{
						"message": fe.Message,
						"code":    fe.Code,
					},
				}
			} else {
				// Прочие ошибки
				status = fiber.StatusInternalServerError
				content = fiber.Map{
					"error": fiber.Map{
						"message": err.Error(),
						"code":    status,
					},
				}
			}
		} else {
			// Парсим ответ тела (если он есть)
			body := c.Response().Body()
			if len(body) > 0 {
				_ = json.Unmarshal(body, &content)
			}
		}

		// Формируем унифицированный ответ
		response := fiber.Map{
			"content": content,
			"metadata": fiber.Map{
				"timestamp":          time.Now().UTC().Format(time.RFC3339),
				"request_id":         uuid.New().String(),
				"path":               c.OriginalURL(),
				"status":             status,
				"processing_time_ms": time.Since(start).Milliseconds(),
			},
		}

		c.Set(fiber.HeaderContentType, fiber.MIMEApplicationJSONCharsetUTF8)
		return c.Status(status).JSON(response)
	}
}
