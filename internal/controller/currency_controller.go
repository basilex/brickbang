package controller

import (
	"errors"

	"github.com/gofiber/fiber/v2"

	"brickbang/internal/service"
	"brickbang/internal/transfer"
	"brickbang/internal/utility"
)

type ICurrencyController interface {
	RegisterRoutes(router fiber.Router)
}

type currencyController struct {
	svc service.ICurrencyService
}

func NewCurrencyController(svc service.ICurrencyService) ICurrencyController {
	return &currencyController{svc: svc}
}

func (rcv *currencyController) RegisterRoutes(router fiber.Router) {
	router.Get("/", rcv.List)
	router.Get("/:id", rcv.Get)
	router.Post("/", rcv.Create)
	router.Put("/:id", rcv.Update)
	router.Delete("/:id", rcv.Delete)
}

// List currencies with pagination
func (rcv *currencyController) List(ctx *fiber.Ctx) error {
	limit, _ := utility.ParseIntQuery(ctx, "limit", 20)
	offset, _ := utility.ParseIntQuery(ctx, "offset", 0)
	order := ctx.Query("order", "id asc")

	currencies, err := rcv.svc.List(ctx.Context(), order, int32(limit), int32(offset))
	if err != nil {
		return utility.RespondWithError(ctx, fiber.StatusInternalServerError, err)
	}

	resp := make([]*transfer.CurrencyResponse, len(currencies))

	for idx, currency := range currencies {
		resp[idx] = &transfer.CurrencyResponse{
			ID:        currency.ID,
			Name:      currency.Name,
			Code:      currency.Code,
			NumCode:   currency.NumCode,
			Symbol:    currency.Symbol,
			CreatedAt: utility.FromPGTimestampToString(currency.CreatedAt),
			UpdatedAt: utility.FromPGTimestampToString(currency.UpdatedAt),
		}
	}

	return ctx.JSON(resp)
}

// Get currency by ID
func (rcv *currencyController) Get(ctx *fiber.Ctx) error {
	id, err := utility.ParseID(ctx, "id")
	if err != nil {
		return err
	}

	currency, err := rcv.svc.GetByID(ctx.Context(), id)
	if err != nil {
		return utility.RespondWithError(ctx, fiber.StatusNotFound, errors.New("currency not found"))
	}

	return ctx.JSON(&transfer.CurrencyResponse{
		ID:        currency.ID,
		Name:      currency.Name,
		CreatedAt: utility.FromPGTimestampToString(currency.CreatedAt),
		UpdatedAt: utility.FromPGTimestampToString(currency.UpdatedAt),
	})
}

// Create new currency
func (rcv *currencyController) Create(ctx *fiber.Ctx) error {
	var req transfer.CurrencyCreateRequest

	if err := utility.ValidateBody(ctx, &req); err != nil {
		return err
	}

	currency, err := rcv.svc.Create(ctx.Context(), &req)
	if err != nil {
		return utility.RespondWithError(ctx, fiber.StatusInternalServerError, err)
	}

	return ctx.Status(fiber.StatusCreated).JSON(&transfer.CurrencyResponse{
		ID:        currency.ID,
		Name:      currency.Name,
		Code:      currency.Code,
		NumCode:   currency.NumCode,
		Symbol:    currency.Symbol,
		CreatedAt: utility.FromPGTimestampToString(currency.CreatedAt),
		UpdatedAt: utility.FromPGTimestampToString(currency.UpdatedAt),
	})
}

// Update currency by ID
func (rcv *currencyController) Update(ctx *fiber.Ctx) error {
	id, err := utility.ParseID(ctx, "id")
	if err != nil {
		return err
	}

	var req transfer.CurrencyUpdateRequest

	if err := utility.ValidateBody(ctx, &req); err != nil {
		return err
	}

	currency, err := rcv.svc.UpdateByID(ctx.Context(), id, &req)
	if err != nil {
		return utility.RespondWithError(ctx, fiber.StatusInternalServerError, err)
	}

	return ctx.JSON(&transfer.CurrencyResponse{
		ID:        currency.ID,
		Name:      currency.Name,
		Code:      currency.Code,
		NumCode:   currency.NumCode,
		Symbol:    currency.Symbol,
		CreatedAt: utility.FromPGTimestampToString(currency.CreatedAt),
		UpdatedAt: utility.FromPGTimestampToString(currency.UpdatedAt),
	})
}

// Delete currency by ID
func (rcv *currencyController) Delete(ctx *fiber.Ctx) error {
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
