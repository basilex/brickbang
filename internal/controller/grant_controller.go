package controller

import (
	"errors"

	"github.com/gofiber/fiber/v2"

	"brickbang/internal/service"
	"brickbang/internal/transfer"
	"brickbang/internal/utility"
)

type IGrantController interface {
	RegisterRoutes(router fiber.Router)
}

type grantController struct {
	svc service.IGrantService
}

func NewGrantController(svc service.IGrantService) IGrantController {
	return &grantController{svc: svc}
}

func (rcv *grantController) RegisterRoutes(router fiber.Router) {
	router.Get("/", rcv.List)
	router.Get("/:id", rcv.Get)
	router.Post("/", rcv.Create)
	router.Put("/:id", rcv.Update)
	router.Delete("/:id", rcv.Delete)
}

// List grants with pagination
func (rcv *grantController) List(ctx *fiber.Ctx) error {
	limit, _ := utility.ParseIntQuery(ctx, "limit", 20)
	offset, _ := utility.ParseIntQuery(ctx, "offset", 0)
	order := ctx.Query("order", "id asc")

	grants, err := rcv.svc.List(ctx.Context(), order, int32(limit), int32(offset))
	if err != nil {
		return utility.RespondWithError(ctx, fiber.StatusInternalServerError, err)
	}

	resp := make([]*transfer.GrantResponse, len(grants))

	for idx, grant := range grants {
		resp[idx] = &transfer.GrantResponse{
			ID:          grant.ID,
			Code:        grant.Code,
			Description: grant.Description,
			CreatedAt:   utility.FromPGTimestampToString(grant.CreatedAt),
			UpdatedAt:   utility.FromPGTimestampToString(grant.UpdatedAt),
		}
	}

	return ctx.JSON(resp)
}

// Get grant by ID
func (rcv *grantController) Get(ctx *fiber.Ctx) error {
	id, err := utility.ParseID(ctx, "id")
	if err != nil {
		return err
	}

	grant, err := rcv.svc.GetByID(ctx.Context(), id)
	if err != nil {
		return utility.RespondWithError(ctx, fiber.StatusNotFound, errors.New("grant not found"))
	}

	return ctx.JSON(&transfer.GrantResponse{
		ID:          grant.ID,
		Code:        grant.Code,
		Description: grant.Description,
		CreatedAt:   utility.FromPGTimestampToString(grant.CreatedAt),
		UpdatedAt:   utility.FromPGTimestampToString(grant.UpdatedAt),
	})
}

// Create new grant
func (rcv *grantController) Create(ctx *fiber.Ctx) error {
	var req transfer.GrantCreateRequest

	if err := utility.ValidateBody(ctx, &req); err != nil {
		return err
	}

	grant, err := rcv.svc.Create(ctx.Context(), req.Code, req.Description)
	if err != nil {
		return utility.RespondWithError(ctx, fiber.StatusInternalServerError, err)
	}

	return ctx.Status(fiber.StatusCreated).JSON(&transfer.RoleResponse{
		ID:          grant.ID,
		Name:        grant.Code,
		Description: grant.Description,
		CreatedAt:   utility.FromPGTimestampToString(grant.CreatedAt),
		UpdatedAt:   utility.FromPGTimestampToString(grant.UpdatedAt),
	})
}

// Update grant by ID
func (rcv *grantController) Update(ctx *fiber.Ctx) error {
	id, err := utility.ParseID(ctx, "id")
	if err != nil {
		return err
	}

	var req transfer.GrantUpdateRequest

	if err := utility.ValidateBody(ctx, &req); err != nil {
		return err
	}

	grant, err := rcv.svc.UpdateByID(ctx.Context(), id, req.Code, req.Description)
	if err != nil {
		return utility.RespondWithError(ctx, fiber.StatusInternalServerError, err)
	}

	return ctx.JSON(&transfer.GrantResponse{
		ID:          grant.ID,
		Code:        grant.Code,
		Description: grant.Description,
		CreatedAt:   utility.FromPGTimestampToString(grant.CreatedAt),
		UpdatedAt:   utility.FromPGTimestampToString(grant.UpdatedAt),
	})
}

// Delete grant by ID
func (rcv *grantController) Delete(ctx *fiber.Ctx) error {
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
