package controller

import (
	"errors"

	"github.com/gofiber/fiber/v2"

	"brickbang/internal/service"
	"brickbang/internal/transfer"
	"brickbang/internal/utility"
)

type ICountryController interface {
	RegisterRoutes(router fiber.Router)
}

type CountryController struct {
	svc service.ICountryService
}

func NewCountryController(svc service.ICountryService) ICountryController {
	return &CountryController{svc: svc}
}

func (rcv *CountryController) RegisterRoutes(router fiber.Router) {
	router.Get("/", rcv.List)
	router.Get("/:id", rcv.Get)
	router.Post("/", rcv.Create)
	router.Put("/:id", rcv.Update)
	router.Delete("/:id", rcv.Delete)
}

// List countries with pagination
func (rcv *CountryController) List(ctx *fiber.Ctx) error {
	limit, _ := utility.ParseIntQuery(ctx, "limit", 20)
	offset, _ := utility.ParseIntQuery(ctx, "offset", 0)
	order := ctx.Query("order", "id asc")

	countries, err := rcv.svc.List(ctx.Context(), order, int32(limit), int32(offset))
	if err != nil {
		return utility.RespondWithError(ctx, fiber.StatusInternalServerError, err)
	}

	resp := make([]*transfer.CountryResponse, len(countries))

	for idx, country := range countries {
		resp[idx] = &transfer.CountryResponse{
			ID:        country.ID,
			Name:      country.Name,
			Iso2:      country.Iso2,
			Iso3:      country.Iso3,
			NumCode:   country.NumCode,
			CreatedAt: utility.FromPGTimestampToString(country.CreatedAt),
			UpdatedAt: utility.FromPGTimestampToString(country.UpdatedAt),
		}
	}

	return ctx.JSON(resp)
}

// Get country by ID
func (rcv *CountryController) Get(ctx *fiber.Ctx) error {
	id, err := utility.ParseID(ctx, "id")
	if err != nil {
		return err
	}

	country, err := rcv.svc.GetByID(ctx.Context(), id)
	if err != nil {
		return utility.RespondWithError(ctx, fiber.StatusNotFound, errors.New("country not found"))
	}

	return ctx.JSON(&transfer.RoleResponse{
		ID:        country.ID,
		Name:      country.Name,
		CreatedAt: utility.FromPGTimestampToString(country.CreatedAt),
		UpdatedAt: utility.FromPGTimestampToString(country.UpdatedAt),
	})
}

// Create new country
func (rcv *CountryController) Create(ctx *fiber.Ctx) error {
	var req transfer.CountryCreateRequest

	if err := utility.ValidateBody(ctx, &req); err != nil {
		return err
	}

	country, err := rcv.svc.Create(ctx.Context(), &req)
	if err != nil {
		return utility.RespondWithError(ctx, fiber.StatusInternalServerError, err)
	}

	return ctx.Status(fiber.StatusCreated).JSON(&transfer.CountryResponse{
		ID:        country.ID,
		Name:      country.Name,
		CreatedAt: utility.FromPGTimestampToString(country.CreatedAt),
		UpdatedAt: utility.FromPGTimestampToString(country.UpdatedAt),
	})
}

// Update country by ID
func (rcv *CountryController) Update(ctx *fiber.Ctx) error {
	id, err := utility.ParseID(ctx, "id")
	if err != nil {
		return err
	}

	var req transfer.CountryUpdateRequest

	if err := utility.ValidateBody(ctx, &req); err != nil {
		return err
	}

	country, err := rcv.svc.UpdateByID(ctx.Context(), id, &req)
	if err != nil {
		return utility.RespondWithError(ctx, fiber.StatusInternalServerError, err)
	}

	return ctx.JSON(&transfer.CountryResponse{
		ID:        country.ID,
		Name:      country.Name,
		Iso2:      country.Iso2,
		Iso3:      country.Iso3,
		NumCode:   country.NumCode,
		CreatedAt: utility.FromPGTimestampToString(country.CreatedAt),
		UpdatedAt: utility.FromPGTimestampToString(country.UpdatedAt),
	})
}

// Delete country by ID
func (rcv *CountryController) Delete(ctx *fiber.Ctx) error {
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
