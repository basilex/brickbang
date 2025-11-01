package controller

import (
	"github.com/gofiber/fiber/v2"

	"brickbang/internal/service"
)

// IAuthController defines interface for AuthController
type IAuthController interface {
	Register(router fiber.Router)
}

// AuthController handles authentication endpoints (login, refresh, me, apikeys)
type AuthController struct {
	service service.IAuthService
}

// NewAuthController constructor
func NewAuthController(s service.IAuthService) IAuthController {
	return &AuthController{service: s}
}

// Register attaches routes to the given router
func (rcv *AuthController) Register(router fiber.Router) {
	auth := router.Group("/auth")

	// JWT login / refresh
	auth.Post("/login", rcv.Login)
	auth.Post("/refresh", rcv.Refresh)
	auth.Get("/me", rcv.Me)

	// API Key management (admin only)
	auth.Post("/apikey", rcv.CreateAPIKey)
	auth.Delete("/apikey/:id", rcv.DeleteAPIKey)
	auth.Get("/apikey", rcv.ListAPIKeys)
}

// ===== Handlers =====

func (rcv *AuthController) Login(ctx *fiber.Ctx) error {
	data, err := rcv.service.Login(ctx)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	return ctx.JSON(data)
}

func (rcv *AuthController) Refresh(ctx *fiber.Ctx) error {
	data, err := rcv.service.Refresh(ctx)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	return ctx.JSON(data)
}

func (rcv *AuthController) Me(ctx *fiber.Ctx) error {
	data, err := rcv.service.Me(ctx)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	return ctx.JSON(data)
}

func (rcv *AuthController) CreateAPIKey(ctx *fiber.Ctx) error {
	data, err := rcv.service.CreateAPIKey(ctx)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	return ctx.JSON(data)
}

func (rcv *AuthController) DeleteAPIKey(ctx *fiber.Ctx) error {
	data, err := rcv.service.DeleteAPIKey(ctx)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	return ctx.JSON(data)
}

func (rcv *AuthController) ListAPIKeys(ctx *fiber.Ctx) error {
	data, err := rcv.service.ListAPIKeys(ctx)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	return ctx.JSON(data)
}
