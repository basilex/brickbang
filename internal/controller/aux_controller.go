package controller

import (
	"time"

	"github.com/gofiber/fiber/v2"
)

var startTime = time.Now()

func RegisterAuxRoutes(router fiber.Router) {
	r := router.Group("/aux")

	r.Get("/health", health)
	r.Get("/version", version)
	r.Get("/uptime", uptime)
}

func health(rcv *fiber.Ctx) error {
	return rcv.JSON(fiber.Map{"status": "ok"})
}

func version(rcv *fiber.Ctx) error {
	return rcv.JSON(fiber.Map{"version": "1.0.0"})
}

func uptime(rcv *fiber.Ctx) error {
	uptime := time.Since(startTime).String()
	return rcv.JSON(fiber.Map{"uptime": uptime})
}
