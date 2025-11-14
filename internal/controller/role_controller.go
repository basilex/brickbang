// controller/role_controller.go
package controller

import (
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"

	"brickbang/internal/mapper"
	"brickbang/internal/service"
	"brickbang/internal/utility"
)

type IRoleController interface {
	RegisterRoutes(router fiber.Router)
}
type RoleController struct {
	svc       service.IRoleService
	validator *validator.Validate
}

func NewRoleController(svc service.IRoleService, validator *validator.Validate) IRoleController {
	return &RoleController{
		svc:       svc,
		validator: validator,
	}
}

func (rc *RoleController) RegisterRoutes(router fiber.Router) {
	router.Get("/", rc.List)
	router.Get("/:id", rc.Get)
	router.Post("/", rc.Create)
	router.Put("/:id", rc.Update)
	router.Delete("/:id", rc.Delete)
}

func (rc *RoleController) List(ctx *fiber.Ctx) error {
	limit, _ := utility.ParseIntQuery(ctx, "limit", 20)
	offset, _ := utility.ParseIntQuery(ctx, "offset", 0)
	order := ctx.Query("order", "id asc")

	list, err := rc.svc.List(ctx.Context(), order, int32(limit), int32(offset))
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	count, err := rc.svc.Count(ctx.Context())
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return ctx.JSON(fiber.Map{
		"data":  mapper.RoleToResponseList(list),
		"count": count,
	})
}

func (rc *RoleController) Get(ctx *fiber.Ctx) error {
	id, err := utility.ParseID(ctx, "id")
	if err != nil {
		return err // already a *fiber.Error
	}

	role, err := rc.svc.GetByID(ctx.Context(), id)
	if err != nil {
		return ctx.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": err.Error()})
	}

	return ctx.JSON(mapper.RoleToResponse(role))
}

func (rc *RoleController) Create(ctx *fiber.Ctx) error {
	var req mapper.RoleCreateRequest

	if err := utility.ValidateBody(ctx, &req, rc.validator); err != nil {
		return err // already a *fiber.Error
	}

	role, err := rc.svc.Create(ctx.Context(), req.Name)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return ctx.Status(fiber.StatusCreated).JSON(mapper.RoleToResponse(role))
}

func (rc *RoleController) Update(ctx *fiber.Ctx) error {
	id, err := utility.ParseID(ctx, "id")
	if err != nil {
		return err
	}

	var req mapper.RoleUpdateRequest

	if err := utility.ValidateBody(ctx, &req, rc.validator); err != nil {
		return err
	}

	role, err := rc.svc.UpdateByID(ctx.Context(), id, req.Name)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return ctx.JSON(mapper.RoleToResponse(role))
}

func (rc *RoleController) Delete(ctx *fiber.Ctx) error {
	id, err := utility.ParseID(ctx, "id")
	if err != nil {
		return err
	}

	_, err = rc.svc.DeleteByID(ctx.Context(), id)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return ctx.SendStatus(fiber.StatusNoContent)
}
