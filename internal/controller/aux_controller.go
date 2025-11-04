package controller

import (
	"brickbang/internal/service"

	"github.com/gofiber/fiber/v2"
)

type IAuxController interface {
	Register(router fiber.Router)
}
type AuxController struct {
	service service.IAuxService
}

func NewAuxController(service service.IAuxService) IAuxController {
	return &AuxController{service: service}
}

func (rcv *AuxController) Register(router fiber.Router) {
	router.Get("/health", rcv.health)
	router.Get("/uptime", rcv.uptime)
	router.Get("/metadata", rcv.metadata)
}

func (rcv *AuxController) health(ctx *fiber.Ctx) error {
	return ctx.JSON(rcv.service.Health())
}

func (rcv *AuxController) uptime(ctx *fiber.Ctx) error {
	return ctx.JSON(rcv.service.Uptime())
}

func (rcv *AuxController) metadata(ctx *fiber.Ctx) error {
	return ctx.JSON(rcv.service.Metadata())
}
