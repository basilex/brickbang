package controller

import (
    "strconv"

    "github.com/gofiber/fiber/v2"
    "github.com/go-playground/validator/v10"

    "brickbang/internal/service"
    "brickbang/internal/mapper"
)

type IRoleController struct {
    svc       service.IRoleService
    validator *validator.Validate
}

func NewRoleController(svc service.IRoleService) IRoleController {
    return &RoleController{
	svc:       svc,
	validator: validator.New(),
    }
}

func (c *RoleController) Register(router fiber.Router) {
    r := router.Group("/roles")
    r.Get("/", c.List)
    r.Get("/:id", c.Get)
    r.Post("/", c.Create)
    r.Put("/:id", c.Update)
    r.Delete("/:id", c.Delete)
}

func (c *RoleController) List(ctx *fiber.Ctx) error {
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

func (c *RoleController) Get(ctx *fiber.Ctx) error {
    id := ctx.Params("id")
    res, err := c.svc.GetByID(ctx.Context(), id)
    if err != nil {
	return ctx.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": err.Error()})
    }
    return ctx.JSON(res)
}

func (c *RoleController) Create(ctx *fiber.Ctx) error {
    var req mapper.RoleCreateRequest
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

func (c *RoleController) Update(ctx *fiber.Ctx) error {
    id := ctx.Params("id")
    var req mapper.RoleUpdateRequest
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

func (c *RoleController) Delete(ctx *fiber.Ctx) error {
    id := ctx.Params("id")
    if err := c.svc.Delete(ctx.Context(), id); err != nil {
	return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
    }
    return ctx.SendStatus(fiber.StatusNoContent)
}
