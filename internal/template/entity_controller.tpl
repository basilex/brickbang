package controller

import (
	"github.com/gofiber/fiber/v2"
	"github.com/go-playground/validator/v10"

	"brickbang/internal/mapper"
	"brickbang/internal/service"
)

type I{{.Entity}}Controller interface {
	Register(router fiber.Router)
}

type {{.Entity}}Controller struct {
	validate *validator.Validate
	service  service.I{{.Entity}}Service
}

func New{{.Entity}}Controller(s service.I{{.Entity}}Service, v *validator.Validate) I{{.Entity}}Controller {
	return &{{.Entity}}Controller{
		validate: v, service:  s,
	}
}

func (rcv *{{.Entity}}Controller) Register(router fiber.Router) {
	router.Get("/", rcv.List)
	router.Get("/:id", rcv.Get)
	router.Post("/", rcv.Create)
	router.Put("/:id", rcv.Update)
	router.Delete("/:id", rcv.Delete)
}

func (rcv *{{.Entity}}Controller) Create(ctx *fiber.Ctx) error {
	var req mapper.{{.Entity}}CreateRequest
	if err := ctx.BodyParser(&req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid request body")
	}
	if err := rcv.validate.Struct(&req); err != nil {
		return fiber.NewError(fiber.StatusUnprocessableEntity, err.Error())
	}
	data, err := rcv.service.Create(ctx.Context(), &req)
	if err != nil {
		return err
	}
	return ctx.Status(fiber.StatusCreated).JSON(data)
}

func (rcv *{{.Entity}}Controller) Update(ctx *fiber.Ctx) error {
	var req mapper.{{.Entity}}UpdateRequest
	if err := ctx.BodyParser(&req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid request body")
	}
	if err := rcv.validate.Struct(&req); err != nil {
		return fiber.NewError(fiber.StatusUnprocessableEntity, err.Error())
	}
	req.ID = ctx.Params("id")
	data, err := rcv.service.Update(ctx.Context(), &req)
	if err != nil {
		return err
	}
	return ctx.Status(fiber.StatusOK).JSON(data)
}

func (rcv *{{.Entity}}Controller) Get(ctx *fiber.Ctx) error {
	id := ctx.Params("id")
	data, err := rcv.service.Get(ctx.Context(), id)
	if err != nil {
		return err
	}
	return ctx.Status(fiber.StatusOK).JSON(data)
}

func (rcv *{{.Entity}}Controller) List(ctx *fiber.Ctx) error {
	data, err := rcv.service.List(ctx.Context())
	if err != nil {
		return err
	}
	return ctx.Status(fiber.StatusOK).JSON(data)
}

func (rcv *{{.Entity}}Controller) Delete(ctx *fiber.Ctx) error {
	id := ctx.Params("id")
	if err := rcv.service.Delete(ctx.Context(), id); err != nil {
		return err
	}
	return ctx.SendStatus(fiber.StatusNoContent)
}
