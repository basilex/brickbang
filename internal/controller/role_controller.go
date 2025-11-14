package controller

import (
	"strconv"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"

	"brickbang/internal/mapper"
	"brickbang/internal/service"
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

func (rcv *RoleController) RegisterRoutes(router fiber.Router) {
	router.Get("/", rcv.List)
	router.Get("/:id", rcv.Get)
	router.Post("/", rcv.Create)
	router.Put("/:id", rcv.Update)
	router.Delete("/:id", rcv.Delete)
}

func (rcv *RoleController) List(ctx *fiber.Ctx) error {
	limit, _ := strconv.Atoi(ctx.Query("limit", "20"))
	offset, _ := strconv.Atoi(ctx.Query("offset", "0"))
	order := ctx.Query("order", "id asc")

	list, err := rcv.svc.List(ctx.Context(), order, int32(limit), int32(offset))
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	count, err := rcv.svc.Count(ctx.Context())
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return ctx.JSON(fiber.Map{
		"data":  mapper.RoleToResponseList(list),
		"count": count,
	})
}

func (rcv *RoleController) Get(ctx *fiber.Ctx) error {
	id := ctx.Params("id")

	role, err := rcv.svc.GetByID(ctx.Context(), id)
	if err != nil {
		return ctx.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": err.Error()})
	}

	return ctx.JSON(mapper.RoleToResponse(role))
}

func (rcv *RoleController) Create(ctx *fiber.Ctx) error {
	var req mapper.RoleCreateRequest

	if err := ctx.BodyParser(&req); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	if err := rcv.validator.Struct(req); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	role, err := rcv.svc.Create(ctx.Context(), req.Name)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return ctx.Status(fiber.StatusCreated).JSON(mapper.RoleToResponse(role))
}

func (rcv *RoleController) Update(ctx *fiber.Ctx) error {
	id := ctx.Params("id")
	var req mapper.RoleUpdateRequest

	if err := ctx.BodyParser(&req); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	if err := rcv.validator.Struct(req); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	role, err := rcv.svc.UpdateByID(ctx.Context(), id, req.Name)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return ctx.JSON(mapper.RoleToResponse(role))
}

func (rcv *RoleController) Delete(ctx *fiber.Ctx) error {
	id := ctx.Params("id")

	_, err := rcv.svc.DeleteByID(ctx.Context(), id)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return ctx.SendStatus(fiber.StatusNoContent)
}
