package controller

import (
	"errors"

	"github.com/gofiber/fiber/v2"

	"brickbang/internal/service"
	"brickbang/internal/transfer"
	"brickbang/internal/utility"
)

type IRoleController interface {
	RegisterRoutes(router fiber.Router)
}

type roleController struct {
	svc service.IRoleService
}

func NewRoleController(svc service.IRoleService) IRoleController {
	return &roleController{svc: svc}
}

func (rcv *roleController) RegisterRoutes(router fiber.Router) {
	router.Get("/", rcv.List)
	router.Get("/:id", rcv.Get)
	router.Post("/", rcv.Create)
	router.Put("/:id", rcv.Update)
	router.Delete("/:id", rcv.Delete)
}

// List roles with pagination
func (rcv *roleController) List(ctx *fiber.Ctx) error {
	limit, _ := utility.ParseIntQuery(ctx, "limit", 20)
	offset, _ := utility.ParseIntQuery(ctx, "offset", 0)
	order := ctx.Query("order", "id asc")

	roles, err := rcv.svc.List(ctx.Context(), order, int32(limit), int32(offset))
	if err != nil {
		return utility.RespondWithError(ctx, fiber.StatusInternalServerError, err)
	}

	resp := make([]*transfer.RoleResponse, len(roles))

	for idx, role := range roles {
		resp[idx] = &transfer.RoleResponse{
			Pid:         role.Pid,
			Name:        role.Name,
			Description: role.Description,
			CreatedAt:   utility.FromPGTimestampToString(role.CreatedAt),
			UpdatedAt:   utility.FromPGTimestampToString(role.UpdatedAt),
		}
	}

	return ctx.JSON(resp)
}

// Get role by ID
func (rcv *roleController) Get(ctx *fiber.Ctx) error {
	id, err := utility.ParseID(ctx, "id")
	if err != nil {
		return err
	}

	role, err := rcv.svc.GetByID(ctx.Context(), id)
	if err != nil {
		return utility.RespondWithError(ctx, fiber.StatusNotFound, errors.New("role not found"))
	}

	return ctx.JSON(&transfer.RoleResponse{
		Pid:         role.Pid,
		Name:        role.Name,
		Description: role.Description,
		CreatedAt:   utility.FromPGTimestampToString(role.CreatedAt),
		UpdatedAt:   utility.FromPGTimestampToString(role.UpdatedAt),
	})
}

// Create new role
func (rcv *roleController) Create(ctx *fiber.Ctx) error {
	var req transfer.RoleCreateRequest

	if err := utility.ValidateBody(ctx, &req); err != nil {
		return err
	}

	role, err := rcv.svc.Create(ctx.Context(), req.Name, req.Description)
	if err != nil {
		return utility.RespondWithError(ctx, fiber.StatusInternalServerError, err)
	}

	return ctx.Status(fiber.StatusCreated).JSON(&transfer.RoleResponse{
		Pid:         role.Pid,
		Name:        role.Name,
		Description: role.Description,
		CreatedAt:   utility.FromPGTimestampToString(role.CreatedAt),
		UpdatedAt:   utility.FromPGTimestampToString(role.UpdatedAt),
	})
}

// Update role by ID
func (rcv *roleController) Update(ctx *fiber.Ctx) error {
	id, err := utility.ParseID(ctx, "id")
	if err != nil {
		return err
	}

	var req transfer.RoleUpdateRequest

	if err := utility.ValidateBody(ctx, &req); err != nil {
		return err
	}

	role, err := rcv.svc.UpdateByID(ctx.Context(), id, req.Name, req.Description)
	if err != nil {
		return utility.RespondWithError(ctx, fiber.StatusInternalServerError, err)
	}

	return ctx.JSON(&transfer.RoleResponse{
		Pid:         role.Pid,
		Name:        role.Name,
		Description: role.Description,
		CreatedAt:   utility.FromPGTimestampToString(role.CreatedAt),
		UpdatedAt:   utility.FromPGTimestampToString(role.UpdatedAt),
	})
}

// Delete role by ID
func (rcv *roleController) Delete(ctx *fiber.Ctx) error {
	id, err := utility.ParseID(ctx, "id")
	if err != nil {
		return err
	}

	_, err = rcv.svc.DeleteByID(ctx.Context(), id)
	if err != nil {
		return utility.RespondWithError(ctx, fiber.StatusInternalServerError, err)
	}

	return ctx.SendStatus(fiber.StatusNoContent)
}
