package controller

import (
	"strconv"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"

	"brickbang/internal/service"
)

type IRoleController interface {
	Register(router fiber.Router)
}

type RoleController struct {
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

// GET /roles?limit=20&offset=0&order=name
func (c *RoleController) List(ctx *fiber.Ctx) error {
	limit, _ := strconv.Atoi(ctx.Query("limit", "20"))
	offset, _ := strconv.Atoi(ctx.Query("offset", "0"))
	order := ctx.Query("order", "id asc")

	roles, err := c.svc.List(ctx.Context(), order, int32(limit), int32(offset))
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	count, err := c.svc.Count(ctx.Context())
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return ctx.JSON(fiber.Map{
		"data":  roles,
		"count": count,
	})
}

// GET /roles/:id
func (c *RoleController) Get(ctx *fiber.Ctx) error {
	id := ctx.Params("id")
	role, err := c.svc.GetByID(ctx.Context(), id)
	if err != nil {
		return ctx.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": err.Error()})
	}
	return ctx.JSON(role)
}

// POST /roles
func (c *RoleController) Create(ctx *fiber.Ctx) error {
	var req struct {
		Name string `json:"name" validate:"required,min=2"`
	}
	if err := ctx.BodyParser(&req); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}
	if err := c.validator.Struct(req); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	role, err := c.svc.Create(ctx.Context(), req.Name)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return ctx.Status(fiber.StatusCreated).JSON(role)
}

// PUT /roles/:id
func (c *RoleController) Update(ctx *fiber.Ctx) error {
	id := ctx.Params("id")
	var req struct {
		Name string `json:"name" validate:"required,min=2"`
	}
	if err := ctx.BodyParser(&req); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}
	if err := c.validator.Struct(req); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	role, err := c.svc.UpdateByID(ctx.Context(), id, req.Name)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return ctx.JSON(role)
}

// DELETE /roles/:id
func (c *RoleController) Delete(ctx *fiber.Ctx) error {
	id := ctx.Params("id")
	_, err := c.svc.DeleteByID(ctx.Context(), id)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return ctx.SendStatus(fiber.StatusNoContent)
}
