package controller

import (
	"errors"

	"github.com/gofiber/fiber/v2"

	"brickbang/internal/model"
	"brickbang/internal/service"
	"brickbang/internal/utility"
)

type IRoleController interface {
	RegisterRoutes(router fiber.Router)
}

type RoleController struct {
	svc service.IRoleService
}

func NewRoleController(svc service.IRoleService) IRoleController {
	return &RoleController{svc: svc}
}

func (rc *RoleController) RegisterRoutes(router fiber.Router) {
	router.Get("/", rc.List)
	router.Get("/:id", rc.Get)
	router.Post("/", rc.Create)
	router.Put("/:id", rc.Update)
	router.Delete("/:id", rc.Delete)
}

// List roles with pagination
func (rc *RoleController) List(ctx *fiber.Ctx) error {
	limit, _ := utility.ParseIntQuery(ctx, "limit", 20)
	offset, _ := utility.ParseIntQuery(ctx, "offset", 0)
	order := ctx.Query("order", "id asc")

	roles, err := rc.svc.List(ctx.Context(), order, int32(limit), int32(offset))
	if err != nil {
		return utility.RespondWithError(ctx, fiber.StatusInternalServerError, err)
	}

	resp := make([]*model.RoleResponse, len(roles))

	for idx, role := range roles {
		resp[idx] = &model.RoleResponse{
			ID:        role.ID,
			Name:      role.Name,
			CreatedAt: utility.FromPGTimestampToString(role.CreatedAt),
			UpdatedAt: utility.FromPGTimestampToString(role.UpdatedAt),
		}
	}

	return ctx.JSON(resp)
}

// Get role by ID
func (rc *RoleController) Get(ctx *fiber.Ctx) error {
	id, err := utility.ParseID(ctx, "id")
	if err != nil {
		return err
	}

	role, err := rc.svc.GetByID(ctx.Context(), id)
	if err != nil {
		return utility.RespondWithError(ctx, fiber.StatusNotFound, errors.New("role not found"))
	}

	return ctx.JSON(&model.RoleResponse{
		ID:        role.ID,
		Name:      role.Name,
		CreatedAt: utility.FromPGTimestampToString(role.CreatedAt),
		UpdatedAt: utility.FromPGTimestampToString(role.UpdatedAt),
	})
}

// Create new role
func (rc *RoleController) Create(ctx *fiber.Ctx) error {
	var req model.RoleCreateRequest

	if err := utility.ValidateBody(ctx, &req); err != nil {
		return err
	}

	role, err := rc.svc.Create(ctx.Context(), req.Name)
	if err != nil {
		return utility.RespondWithError(ctx, fiber.StatusInternalServerError, err)
	}

	return ctx.Status(fiber.StatusCreated).JSON(&model.RoleResponse{
		ID:        role.ID,
		Name:      role.Name,
		CreatedAt: utility.FromPGTimestampToString(role.CreatedAt),
		UpdatedAt: utility.FromPGTimestampToString(role.UpdatedAt),
	})
}

// Update role by ID
func (rc *RoleController) Update(ctx *fiber.Ctx) error {
	id, err := utility.ParseID(ctx, "id")
	if err != nil {
		return err
	}

	var req model.RoleUpdateRequest

	if err := utility.ValidateBody(ctx, &req); err != nil {
		return err
	}

	role, err := rc.svc.UpdateByID(ctx.Context(), id, req.Name)
	if err != nil {
		return utility.RespondWithError(ctx, fiber.StatusInternalServerError, err)
	}

	return ctx.JSON(&model.RoleResponse{
		ID:        role.ID,
		Name:      role.Name,
		CreatedAt: utility.FromPGTimestampToString(role.CreatedAt),
		UpdatedAt: utility.FromPGTimestampToString(role.UpdatedAt),
	})
}

// Delete role by ID
func (rc *RoleController) Delete(ctx *fiber.Ctx) error {
	id, err := utility.ParseID(ctx, "id")
	if err != nil {
		return err
	}

	_, err = rc.svc.DeleteByID(ctx.Context(), id)
	if err != nil {
		return utility.RespondWithError(ctx, fiber.StatusInternalServerError, err)
	}

	return ctx.SendStatus(fiber.StatusNoContent)
}
