package controller

import (
	"strconv"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"

	"brickbang/internal/mapper"
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

	list, err := c.svc.List(ctx.Context(), order, int32(limit), int32(offset))
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	resp := mapper.RoleToResponseList(list)
	count, _ := c.svc.Count(ctx.Context())

	return ctx.JSON(fiber.Map{
		"data":  resp,
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

	return ctx.JSON(mapper.RoleToResponse(role))
}

// POST /roles
func (c *RoleController) Create(ctx *fiber.Ctx) error {
	var req mapper.RoleCreateRequest

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

	return ctx.Status(fiber.StatusCreated).JSON(mapper.RoleToResponse(role))
}

// PUT /roles/:id
func (c *RoleController) Update(ctx *fiber.Ctx) error {
	id := ctx.Params("id")
	var req mapper.RoleUpdateRequest

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

	return ctx.JSON(mapper.RoleToResponse(role))
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
