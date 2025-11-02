package controller

import (
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
	// Registration route
	router.Post("/register", rcv.Registering)

	// User authentication routes
	router.Post("/login", rcv.Login)
	router.Post("/refresh", rcv.Refresh)
	router.Get("/me", rcv.Me)

	// API key management routes
	router.Post("/apikey", rcv.CreateAPIKey)
	router.Delete("/apikey/:id", rcv.DeleteAPIKey)
	router.Get("/apikey", rcv.ListAPIKeys)
}

// ============ Handlers ============

func (rcv *AuthController) Registering(ctx *fiber.Ctx) error {
	data, err := rcv.service.Register(ctx)
	if err != nil {
		return err // глобальный ErrorHandler поймает
	}
	return ctx.Status(fiber.StatusCreated).JSON(data)
}

func (rcv *AuthController) Login(ctx *fiber.Ctx) error {
	data, err := rcv.service.Login(ctx)
	if err != nil {
		return err
	}
	return ctx.Status(fiber.StatusOK).JSON(data)
}

func (rcv *AuthController) Refresh(ctx *fiber.Ctx) error {
	data, err := rcv.service.Refresh(ctx)
	if err != nil {
		return err
	}
	return ctx.Status(fiber.StatusOK).JSON(data)
}

func (rcv *AuthController) Me(ctx *fiber.Ctx) error {
	data, err := rcv.service.Me(ctx)
	if err != nil {
		return err
	}
	return ctx.Status(fiber.StatusOK).JSON(data)
}

func (rcv *AuthController) CreateAPIKey(ctx *fiber.Ctx) error {
	data, err := rcv.service.CreateAPIKey(ctx)
	if err != nil {
		return err
	}
	return ctx.Status(fiber.StatusCreated).JSON(data)
}

func (rcv *AuthController) DeleteAPIKey(ctx *fiber.Ctx) error {
	data, err := rcv.service.DeleteAPIKey(ctx)
	if err != nil {
		return err
	}
	return ctx.Status(fiber.StatusOK).JSON(data)
}

func (rcv *AuthController) ListAPIKeys(ctx *fiber.Ctx) error {
	data, err := rcv.service.ListAPIKeys(ctx)
	if err != nil {
		return err
	}
	return ctx.Status(fiber.StatusOK).JSON(data)
}
