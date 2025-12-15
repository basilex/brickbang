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
	router.Get("/:pid", rcv.Get)
	router.Post("/", rcv.Create)
	router.Put("/:pid", rcv.Update)
	router.Delete("/:pid", rcv.Delete)
}

// List grants with pagination
func (rcv *grantController) List(ctx *fiber.Ctx) error {
	limit, _ := utility.ParseIntQuery(ctx, "limit", 20)
	offset, _ := utility.ParseIntQuery(ctx, "offset", 0)
	order := ctx.Query("order", "pid asc")

	grants, err := rcv.svc.List(ctx.Context(), order, int32(limit), int32(offset))
	if err != nil {
		return utility.RespondWithError(ctx, fiber.StatusInternalServerError, err)
	}

	resp := make([]*transfer.GrantResponse, len(grants))

	for idx, grant := range grants {
		resp[idx] = &transfer.GrantResponse{
			Pid:         grant.Pid,
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
	pid, err := utility.ParseID(ctx, "pid")
	if err != nil {
		return err
	}

	grant, err := rcv.svc.GetByID(ctx.Context(), pid)
	if err != nil {
		return utility.RespondWithError(ctx, fiber.StatusNotFound, errors.New("grant not found"))
	}

	return ctx.JSON(&transfer.GrantResponse{
		Pid:         grant.Pid,
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

	return ctx.Status(fiber.StatusCreated).JSON(
		&transfer.GrantResponse{
			Pid:         grant.Pid,
			Code:        grant.Code,
			Description: grant.Description,
			CreatedAt:   utility.FromPGTimestampToString(grant.CreatedAt),
			UpdatedAt:   utility.FromPGTimestampToString(grant.UpdatedAt),
		})
}

// Update grant by ID
func (rcv *grantController) Update(ctx *fiber.Ctx) error {
	pid, err := utility.ParseID(ctx, "pid")
	if err != nil {
		return err
	}

	var req transfer.GrantUpdateRequest

	if err := utility.ValidateBody(ctx, &req); err != nil {
		return err
	}

	grant, err := rcv.svc.UpdateByID(ctx.Context(), pid, req.Code, req.Description)
	if err != nil {
		return utility.RespondWithError(ctx, fiber.StatusInternalServerError, err)
	}

	return ctx.JSON(&transfer.GrantResponse{
		Pid:         grant.Pid,
		Code:        grant.Code,
		Description: grant.Description,
		CreatedAt:   utility.FromPGTimestampToString(grant.CreatedAt),
		UpdatedAt:   utility.FromPGTimestampToString(grant.UpdatedAt),
	})
}

// Delete grant by Pid
func (rcv *grantController) Delete(ctx *fiber.Ctx) error {
	pid, err := utility.ParseID(ctx, "pid")
	if err != nil {
		return err
	}

	_, err = rcv.svc.DeleteByID(ctx.Context(), pid)
	if err != nil {
		return utility.RespondWithError(ctx, fiber.StatusInternalServerError, err)
	}

	return ctx.SendStatus(fiber.StatusNoContent)
}
