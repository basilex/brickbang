package utility

import "github.com/gofiber/fiber/v2"

func RespondWithError(ctx *fiber.Ctx, status int, err error) error {
	return ctx.Status(status).JSON(fiber.Map{"error": err.Error()})
}
