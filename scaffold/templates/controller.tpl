package controller

import (
    "strconv"

    "github.com/gofiber/fiber/v2"
    "github.com/go-playground/validator/v10"

    "{{.ModulePath}}/internal/service"
    "{{.ModulePath}}/internal/mapper"
)

type I{{.Entity}}Controller struct {
    svc       service.I{{.Entity}}Service
    validator *validator.Validate
}

func New{{.Entity}}Controller(svc service.I{{.Entity}}Service) I{{.Entity}}Controller {
    return &{{.Entity}}Controller{
	svc:       svc,
	validator: validator.New(),
    }
}

func (c *{{.Entity}}Controller) Register(router fiber.Router) {
    r := router.Group("/{{.EntityPluralLo}}")
    r.Get("/", c.List)
    r.Get("/:id", c.Get)
    r.Post("/", c.Create)
    r.Put("/:id", c.Update)
    r.Delete("/:id", c.Delete)
}

func (c *{{.Entity}}Controller) List(ctx *fiber.Ctx) error {
    limit, _ := strconv.Atoi(ctx.Query("limit", "20"))
    offset, _ := strconv.Atoi(ctx.Query("offset", "0"))
    order := ctx.Query("order", "asc")

    list, count, err := c.svc.List(ctx.Context(), int32(limit), int32(offset), order)
    if err != nil {
	return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
    }

    return ctx.JSON(fiber.Map{
	"data":  list,
	"count": count,
    })
}

func (c *{{.Entity}}Controller) Get(ctx *fiber.Ctx) error {
    id := ctx.Params("id")
    res, err := c.svc.GetByID(ctx.Context(), id)
    if err != nil {
	return ctx.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": err.Error()})
    }
    return ctx.JSON(res)
}

func (c *{{.Entity}}Controller) Create(ctx *fiber.Ctx) error {
    var req mapper.{{.Entity}}CreateRequest
    if err := ctx.BodyParser(&req); err != nil {
	return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
    }

    if err := c.validator.Struct(req); err != nil {
	return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
    }

    res, err := c.svc.Create(ctx.Context(), &req)
    if err != nil {
	return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
    }

    return ctx.Status(fiber.StatusCreated).JSON(res)
}

func (c *{{.Entity}}Controller) Update(ctx *fiber.Ctx) error {
    id := ctx.Params("id")
    var req mapper.{{.Entity}}UpdateRequest
    if err := ctx.BodyParser(&req); err != nil {
	return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
    }

    if err := c.validator.Struct(req); err != nil {
	return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
    }

    res, err := c.svc.Update(ctx.Context(), id, &req)
    if err != nil {
	return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
    }

    return ctx.JSON(res)
}

func (c *{{.Entity}}Controller) Delete(ctx *fiber.Ctx) error {
    id := ctx.Params("id")
    if err := c.svc.Delete(ctx.Context(), id); err != nil {
	return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
    }
    return ctx.SendStatus(fiber.StatusNoContent)
}
