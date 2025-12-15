package controller

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http/httptest"
	"testing"

	"brickbang/internal/transfer"

	"github.com/gofiber/fiber/v2"
)

type fakeAuthService struct{}

func (f *fakeAuthService) Register(ctx context.Context, req *transfer.AuthRegisterRequest) (*transfer.AuthUserResponse, error) {
	return &transfer.AuthUserResponse{ID: "u1", Username: req.Username}, nil
}
func (f *fakeAuthService) Login(ctx context.Context, req *transfer.AuthLoginRequest) (*transfer.AuthLoginResponse, error) {
	return &transfer.AuthLoginResponse{User: &transfer.AuthUserResponse{ID: "u1", Username: req.Username}}, nil
}
func (f *fakeAuthService) Logout(ctx context.Context, sessionID string) error { return nil }
func (f *fakeAuthService) Refresh(ctx context.Context, userID, refreshToken string) (*transfer.AuthLoginResponse, error) {
	return &transfer.AuthLoginResponse{User: &transfer.AuthUserResponse{ID: userID}}, nil
}
func (f *fakeAuthService) Me(ctx context.Context, userID string) (*transfer.AuthMeResponse, error) {
	return &transfer.AuthMeResponse{User: &transfer.AuthUserResponse{ID: userID}}, nil
}
func (f *fakeAuthService) Block(ctx context.Context, userID string, blocked bool) (*transfer.AuthUserResponse, error) {
	return &transfer.AuthUserResponse{ID: userID}, nil
}

// adapt signatures to match IAuthService in tests by using type assertions where needed

func TestAuthControllerPublicRoutes(t *testing.T) {
	app := fiber.New()
	svc := &fakeAuthService{}
	ctrl := NewAuthController(svc)
	grp := app.Group("/auth")
	ctrl.RegisterPublicRoutes(grp)

	// Register
	reg := transfer.AuthRegisterRequest{Username: "bob", Password: "secret"}
	body, _ := json.Marshal(reg)
	req := httptest.NewRequest("POST", "/auth/register", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("app.Test error: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 201 && resp.StatusCode != 200 {
		t.Fatalf("unexpected status: %v", resp.StatusCode)
	}
}
