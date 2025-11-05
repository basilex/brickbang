package controller

import (
	"{{.ModulePath}}/internal/service"

	"github.com/gofiber/fiber/v2"
)

type I{{.Entity}}Controller interface {
	Register(router fiber.Router)
}

type {{.Entity}}Controller struct {
	service service.I{{.Entity}}Service
}

func New{{.Entity}}Controller(s service.I{{.Entity}}Service) I{{.Entity}}Controller {
	return &{{.Entity}}Controller{service: s}
}

func (rcv *{{.Entity}}Controller) Register(router fiber.Router) {
	router.Get("/", rcv.GetAll)
	router.Get("/:id", rcv.GetByID)
	router.Post("/", rcv.Create)
	router.Put("/:id", rcv.Update)
	router.Delete("/:id", rcv.Delete)
}

// GET /api/v1/{{.EntityLower}}
func (rcv *{{.Entity}}Controller) GetAll(ctx *fiber.Ctx) error {
	data, err := rcv.service.GetAll(ctx)
	if err != nil {
		return err
	}
	return ctx.Status(fiber.StatusOK).JSON(data)
}

// GET /api/v1/{{.EntityLower}}/:id
func (rcv *{{.Entity}}Controller) GetByID(ctx *fiber.Ctx) error {
	id := ctx.Params("id")
	data, err := rcv.service.GetByID(ctx, id)
	if err != nil {
		return err
	}
	return ctx.Status(fiber.StatusOK).JSON(data)
}

// POST /api/v1/{{.EntityLower}}
func (rcv *{{.Entity}}Controller) Create(ctx *fiber.Ctx) error {
	data, err := rcv.service.Create(ctx)
	if err != nil {
		return err
	}
	return ctx.Status(fiber.StatusCreated).JSON(data)
}

// PUT /api/v1/{{.EntityLower}}/:id
func (rcv *{{.Entity}}Controller) Update(ctx *fiber.Ctx) error {
	id := ctx.Params("id")
	data, err := rcv.service.Update(ctx, id)
	if err != nil {
		return err
	}
	return ctx.Status(fiber.StatusOK).JSON(data)
}

// DELETE /api/v1/{{.EntityLower}}/:id
func (rcv *{{.Entity}}Controller) Delete(ctx *fiber.Ctx) error {
	id := ctx.Params("id")
	if err := rcv.service.Delete(ctx, id); err != nil {
		return err
	}
	return ctx.SendStatus(fiber.StatusNoContent)
}
