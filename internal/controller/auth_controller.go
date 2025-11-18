// internal/controller/auth_controller.go
package controller

import (
	"fmt"

	"github.com/gofiber/fiber/v2"

	"brickbang/internal/service"
	"brickbang/internal/transfer"
	"brickbang/internal/utility"
)

type IAuthController interface {
	RegisterPublicRoutes(router fiber.Router)
	RegisterPrivateRoutes(router fiber.Router)
}
type authController struct {
	svc service.IAuthService
}

func NewAuthController(svc service.IAuthService) IAuthController {
	return &authController{svc: svc}
}

func (c *authController) RegisterPublicRoutes(router fiber.Router) {
	router.Post("/register", c.RegisterUser)
	router.Post("/login", c.Login)
	router.Post("/refresh", c.Refresh)
}

func (c *authController) RegisterPrivateRoutes(router fiber.Router) {
	router.Get("/me", c.Me)
	router.Post("/logout", c.Logout)
	router.Post("/block/:id", c.Block)
}

func (rcv *authController) RegisterUser(ctx *fiber.Ctx) error {
	var req transfer.AuthRegisterRequest

	if err := utility.ValidateBody(ctx, &req); err != nil {
		return utility.RespondWithError(ctx, fiber.StatusBadRequest, err)
	}

	resp, err := rcv.svc.Register(ctx.Context(), &req)
	if err != nil {
		return utility.RespondWithError(ctx, fiber.StatusInternalServerError, err)
	}

	return ctx.Status(fiber.StatusCreated).JSON(resp)
}

func (rcv *authController) Login(ctx *fiber.Ctx) error {
	var req transfer.AuthLoginRequest

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

func (rcv *authController) Refresh(ctx *fiber.Ctx) error {
	var req transfer.AuthRefreshRequest

	if err := utility.ValidateBody(ctx, &req); err != nil {
		return utility.RespondWithError(ctx, fiber.StatusBadRequest, err)
	}

	resp, err := rcv.svc.Refresh(ctx.Context(), req.UserID, req.RefreshToken)
	if err != nil {
		return utility.RespondWithError(ctx, fiber.StatusUnauthorized, err)
	}

	return ctx.JSON(resp)
}

func (rcv *authController) Me(ctx *fiber.Ctx) error {
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

func (rcv *authController) Logout(ctx *fiber.Ctx) error {
	var req transfer.AuthLogoutRequest

	if err := utility.ValidateBody(ctx, &req); err != nil {
		return utility.RespondWithError(ctx, fiber.StatusBadRequest, err)
	}

	if err := rcv.svc.Logout(ctx.Context(), req.SessionID); err != nil {
		return utility.RespondWithError(ctx, fiber.StatusInternalServerError, err)
	}
	return ctx.SendStatus(fiber.StatusNoContent)
}

func (rcv *authController) Block(ctx *fiber.Ctx) error {
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
