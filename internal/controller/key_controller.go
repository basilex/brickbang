package controller

import (
	"github.com/gofiber/fiber/v2"

	"brickbang/internal/service"
)

type IKeyController interface {
	Register(router fiber.Router)
}

type KeyController struct {
	service service.IKeyService
}

func NewKeyController(s service.IKeyService) IKeyController {
	return &KeyController{service: s}
}

func (rcv *KeyController) Register(router fiber.Router) {
	router.Get("/key", rcv.ListAPIKeys)
	router.Post("/key", rcv.CreateAPIKey)
	router.Delete("/key/:id", rcv.DeleteAPIKey)
}

func (rcv *KeyController) ListAPIKeys(ctx *fiber.Ctx) error {
	data, err := rcv.service.ListAPIKeys(ctx)
	if err != nil {
		return err
	}
	return ctx.Status(fiber.StatusOK).JSON(data)
}

func (rcv *KeyController) CreateAPIKey(ctx *fiber.Ctx) error {
	data, err := rcv.service.CreateAPIKey(ctx)
	if err != nil {
		return err
	}
	return ctx.Status(fiber.StatusCreated).JSON(data)
}

func (rcv *KeyController) DeleteAPIKey(ctx *fiber.Ctx) error {
	data, err := rcv.service.DeleteAPIKey(ctx)
	if err != nil {
		return err
	}
	return ctx.Status(fiber.StatusOK).JSON(data)
}
