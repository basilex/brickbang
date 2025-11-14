// internal/controller/auth_controller.go
package controller

import (
	"fmt"

	"github.com/gofiber/fiber/v2"

	"brickbang/internal/model"
	"brickbang/internal/service"
	"brickbang/internal/utility"
)

type IAuthController interface {
	RegisterPublicRoutes(router fiber.Router)
	RegisterPrivateRoutes(router fiber.Router)
}
type AuthController struct {
	svc service.IAuthService
}

func NewAuthController(svc service.IAuthService) IAuthController {
	return &AuthController{svc: svc}
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

func (rcv *AuthController) RegisterUser(ctx *fiber.Ctx) error {
	var req model.AuthRegisterRequest

	if err := utility.ValidateBody(ctx, &req); err != nil {
		return utility.RespondWithError(ctx, fiber.StatusBadRequest, err)
	}

	resp, err := rcv.svc.Register(ctx.Context(), &req)
	if err != nil {
		return utility.RespondWithError(ctx, fiber.StatusInternalServerError, err)
	}

	return ctx.Status(fiber.StatusCreated).JSON(resp)
}

func (rcv *AuthController) Login(ctx *fiber.Ctx) error {
	var req model.AuthLoginRequest

	if err := utility.ValidateBody(ctx, &req); err != nil {
		return utility.RespondWithError(ctx, fiber.StatusBadRequest, err)
	}

	req.IpAddress = ctx.IP()
	req.UserAgent = ctx.Get("User-Agent")

	resp, err := rcv.svc.Login(ctx.Context(), &req)
	if err != nil {
		return utility.RespondWithError(ctx, fiber.StatusUnauthorized, err)
	}

	return ctx.JSON(resp)
}

func (rcv *AuthController) Refresh(ctx *fiber.Ctx) error {
	var req model.AuthRefreshRequest

	if err := utility.ValidateBody(ctx, &req); err != nil {
		return utility.RespondWithError(ctx, fiber.StatusBadRequest, err)
	}

	resp, err := rcv.svc.Refresh(ctx.Context(), req.UserID, req.RefreshToken)
	if err != nil {
		return utility.RespondWithError(ctx, fiber.StatusUnauthorized, err)
	}

	return ctx.JSON(resp)
}

func (rcv *AuthController) Me(ctx *fiber.Ctx) error {
	userID := ctx.Locals("user_id")
	if userID == nil {
		return utility.RespondWithError(ctx, fiber.StatusUnauthorized, fmt.Errorf("unauthorized"))
	}

	resp, err := rcv.svc.Me(ctx.Context(), userID.(string))
	if err != nil {
		return utility.RespondWithError(ctx, fiber.StatusInternalServerError, err)
	}

	return ctx.JSON(resp)
}

func (rcv *AuthController) Logout(ctx *fiber.Ctx) error {
	var req model.AuthLogoutRequest

	if err := utility.ValidateBody(ctx, &req); err != nil {
		return utility.RespondWithError(ctx, fiber.StatusBadRequest, err)
	}

	if err := rcv.svc.Logout(ctx.Context(), req.SessionID); err != nil {
		return utility.RespondWithError(ctx, fiber.StatusInternalServerError, err)
	}
	return ctx.SendStatus(fiber.StatusNoContent)
}

func (rcv *AuthController) Block(ctx *fiber.Ctx) error {
	id := ctx.Params("id")

	var req struct {
		Blocked bool `json:"blocked" validate:"required"`
	}
	if err := utility.ValidateBody(ctx, &req); err != nil {
		return utility.RespondWithError(ctx, fiber.StatusBadRequest, err)
	}

	resp, err := rcv.svc.Block(ctx.Context(), id, req.Blocked)
	if err != nil {
		return utility.RespondWithError(ctx, fiber.StatusInternalServerError, err)
	}
	return ctx.JSON(resp)
}
