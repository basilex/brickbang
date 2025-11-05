package controller

import (
	"strconv"

	"github.com/gofiber/fiber/v2"

	"brickbang/internal/service"
	"brickbang/storage/dbs"
)

type ICountryController interface {
	Register(router fiber.Router)
}

type CountryController struct {
	svc service.ICountryService
}

func NewCountryController(svc service.ICountryService) ICountryController {
	return &CountryController{svc: svc}
}

func (rcv *CountryController) Register(router fiber.Router) {
	router.Get("/", rcv.GetAll)
	router.Get("/with-currencies", rcv.GetAllWithCurrencies)
	router.Get("/:id", rcv.GetByID)
	router.Post("/", rcv.Create)
	router.Put("/:id", rcv.Update)
	router.Delete("/:id", rcv.Delete)
	router.Get("/count", rcv.Count)
}

func (rcv *CountryController) GetAll(ctx *fiber.Ctx) error {
	order := ctx.Query("order", "name")
	offset, _ := strconv.Atoi(ctx.Query("offset", "0"))
	limit, _ := strconv.Atoi(ctx.Query("limit", "20"))

	data, err := rcv.svc.GetAll(ctx.Context(), order, int32(offset), int32(limit))
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}
	return ctx.JSON(data)
}

func (rcv *CountryController) GetAllWithCurrencies(ctx *fiber.Ctx) error {
	order := ctx.Query("order", "name")
	offset, _ := strconv.Atoi(ctx.Query("offset", "0"))
	limit, _ := strconv.Atoi(ctx.Query("limit", "20"))

	data, err := rcv.svc.GetAllWithCurrencies(ctx.Context(), order, int32(offset), int32(limit))
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}
	return ctx.JSON(data)
}

func (rcv *CountryController) GetByID(ctx *fiber.Ctx) error {
	id := ctx.Params("id")
	data, err := rcv.svc.GetByID(ctx.Context(), id)
	if err != nil {
		return fiber.NewError(fiber.StatusNotFound, err.Error())
	}
	return ctx.JSON(data)
}

func (rcv *CountryController) Create(ctx *fiber.Ctx) error {
	var dto dbs.CountryNewParams
	if err := ctx.BodyParser(&dto); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}
	data, err := rcv.svc.Create(ctx.Context(), &dto)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}
	return ctx.Status(fiber.StatusCreated).JSON(data)
}

func (rcv *CountryController) Update(ctx *fiber.Ctx) error {
	id := ctx.Params("id")
	var dto dbs.CountryUpdateByIDParams
	if err := ctx.BodyParser(&dto); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}
	dto.ID = id

	data, err := rcv.svc.Update(ctx.Context(), &dto)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}
	return ctx.JSON(data)
}

func (rcv *CountryController) Delete(ctx *fiber.Ctx) error {
	id := ctx.Params("id")
	deletedID, err := rcv.svc.Delete(ctx.Context(), id)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}
	return ctx.JSON(fiber.Map{"deleted_id": deletedID})
}

func (rcv *CountryController) Count(ctx *fiber.Ctx) error {
	count, err := rcv.svc.Count(ctx.Context())
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}
	return ctx.JSON(fiber.Map{"count": count})
}
