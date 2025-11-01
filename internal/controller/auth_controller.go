package controller

import (
	"brickbang/internal/exception"
	"brickbang/internal/service"

	"github.com/gofiber/fiber/v2"
)

type IAuthController interface {
	Register(router fiber.Router)
}

type AuthController struct {
	service service.IAuthService
}

func NewAuthController(s service.IAuthService) IAuthController {
	return &AuthController{service: s}
}

func (rcv *AuthController) Register(router fiber.Router) {
	auth := router.Group("/auth")

	auth.Post("/login", rcv.Login)
	auth.Post("/refresh", rcv.Refresh)
	auth.Get("/me", rcv.Me)

	auth.Post("/apikey", rcv.CreateAPIKey)
	auth.Delete("/apikey/:id", rcv.DeleteAPIKey)
	auth.Get("/apikey", rcv.ListAPIKeys)
}

func handleError(ctx *fiber.Ctx, err error) error {
	status := fiber.StatusInternalServerError
	msg := err.Error()

	if ae, ok := err.(*exception.AppError); ok {
		status = ae.Code
		msg = ae.Message
	}

	return ctx.Status(status).JSON(map[string]any{"error": msg})
}

func (rcv *AuthController) Login(ctx *fiber.Ctx) error {
	data, err := rcv.service.Login(ctx)
	if err != nil {
		return handleError(ctx, err)
	}
	return ctx.Status(fiber.StatusOK).JSON(data)
}

func (rcv *AuthController) Refresh(ctx *fiber.Ctx) error {
	data, err := rcv.service.Refresh(ctx)
	if err != nil {
		return handleError(ctx, err)
	}
	return ctx.Status(fiber.StatusOK).JSON(data)
}

func (rcv *AuthController) Me(ctx *fiber.Ctx) error {
	data, err := rcv.service.Me(ctx)
	if err != nil {
		return handleError(ctx, err)
	}
	return ctx.Status(fiber.StatusOK).JSON(data)
}

func (rcv *AuthController) CreateAPIKey(ctx *fiber.Ctx) error {
	data, err := rcv.service.CreateAPIKey(ctx)
	if err != nil {
		return handleError(ctx, err)
	}
	return ctx.Status(fiber.StatusCreated).JSON(data)
}

func (rcv *AuthController) DeleteAPIKey(ctx *fiber.Ctx) error {
	data, err := rcv.service.DeleteAPIKey(ctx)
	if err != nil {
		return handleError(ctx, err)
	}
	return ctx.Status(fiber.StatusOK).JSON(data)
}

func (rcv *AuthController) ListAPIKeys(ctx *fiber.Ctx) error {
	data, err := rcv.service.ListAPIKeys(ctx)
	if err != nil {
		return handleError(ctx, err)
	}
	return ctx.Status(fiber.StatusOK).JSON(data)
}
