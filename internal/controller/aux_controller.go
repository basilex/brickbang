package controller

import (
	"brickbang/internal/service"

	"github.com/gofiber/fiber/v2"
)

type IAuxController interface {
	RegisterRoutes(router fiber.Router)
}
type auxController struct {
	service service.IAuxService
}

func NewAuxController(service service.IAuxService) IAuxController {
	return &auxController{service: service}
}

func (rcv *auxController) RegisterRoutes(router fiber.Router) {
	router.Get("/health", rcv.health)
	router.Get("/uptime", rcv.uptime)
	router.Get("/metadata", rcv.metadata)
}

func (rcv *auxController) health(ctx *fiber.Ctx) error {
	return ctx.JSON(rcv.service.Health())
}

func (rcv *auxController) uptime(ctx *fiber.Ctx) error {
	return ctx.JSON(rcv.service.Uptime())
}

func (rcv *auxController) metadata(ctx *fiber.Ctx) error {
	return ctx.JSON(rcv.service.Metadata())
}
