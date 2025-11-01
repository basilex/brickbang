package controller

import (
	"github.com/gofiber/fiber/v2"
)

func RegisterUserRoutes(router fiber.Router) {
	r := router.Group("/users")

	r.Get("/", listUsers)
	r.Get("/:id", getUser)
}

func listUsers(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{"users": []string{}})
}

func getUser(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{"id": c.Params("id")})
}
