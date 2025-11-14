package controller

import (
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"

	"brickbang/internal/model"
	"brickbang/internal/service"
)

type IAuthController interface {
	RegisterPublicRoutes(router fiber.Router)
	RegisterPrivateRoutes(router fiber.Router)
}

type AuthController struct {
	svc       service.IAuthService
	validator *validator.Validate
}

func NewAuthController(svc service.IAuthService) IAuthController {
	return &AuthController{
		svc:       svc,
		validator: validator.New(),
	}
}

func (c *AuthController) RegisterPublicRoutes(router fiber.Router) {
	router.Post("/register", c.RegisterUser)
	router.Post("/login", c.Login)
	router.Post("/refresh", c.Refresh)
}

func (c *AuthController) RegisterPrivateRoutes(router fiber.Router) {
	router.Get("/me", c.Me)
	router.Post("/logout", c.Logout)
	router.Post("/block/:id", c.Block)
}

func (c *AuthController) RegisterUser(ctx *fiber.Ctx) error {
	var req model.AuthRegisterRequest
	if err := ctx.BodyParser(&req); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	if err := c.validator.Struct(req); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	resp, err := c.svc.Register(ctx.Context(), &req)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return ctx.Status(fiber.StatusCreated).JSON(resp)
}

func (c *AuthController) Login(ctx *fiber.Ctx) error {
	var req model.AuthLoginRequest

	if err := ctx.BodyParser(&req); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	req.IpAddress = ctx.IP()
	req.UserAgent = ctx.Get("User-Agent")

	resp, err := c.svc.Login(ctx.Context(), &req)
	if err != nil {
		return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": err.Error()})
	}

	return ctx.JSON(resp)
}

func (c *AuthController) Logout(ctx *fiber.Ctx) error {
	var req struct {
		SessionID string `json:"session_id" validate:"required"`
	}

	if err := ctx.BodyParser(&req); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	if err := c.validator.Struct(req); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	if err := c.svc.Logout(ctx.Context(), req.SessionID); err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return ctx.SendStatus(fiber.StatusNoContent)
}

func (c *AuthController) Refresh(ctx *fiber.Ctx) error {
	var req struct {
		UserID       string `json:"user_id" validate:"required"`
		RefreshToken string `json:"refresh_token" validate:"required"`
	}

	if err := ctx.BodyParser(&req); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	if err := c.validator.Struct(req); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	resp, err := c.svc.Refresh(ctx.Context(), req.UserID, req.RefreshToken)
	if err != nil {
		return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": err.Error()})
	}

	return ctx.JSON(resp)
}

func (c *AuthController) Me(ctx *fiber.Ctx) error {
	userID := ctx.Locals("user_id")
	if userID == nil {
		return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "unauthorized"})
	}

	resp, err := c.svc.Me(ctx.Context(), userID.(string))
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return ctx.JSON(resp)
}

func (c *AuthController) Block(ctx *fiber.Ctx) error {
	id := ctx.Params("id")

	var req struct {
		Blocked bool `json:"blocked" validate:"required"`
	}

	if err := ctx.BodyParser(&req); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	if err := c.validator.Struct(req); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	resp, err := c.svc.Block(ctx.Context(), id, req.Blocked)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return ctx.JSON(resp)
}
