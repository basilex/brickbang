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

type fakeAuthServicePrivate struct{}

func (f *fakeAuthServicePrivate) Register(ctx context.Context, req *transfer.AuthRegisterRequest) (*transfer.AuthUserResponse, error) {
	return &transfer.AuthUserResponse{ID: "u1", Username: req.Username}, nil
}
func (f *fakeAuthServicePrivate) Login(ctx context.Context, req *transfer.AuthLoginRequest) (*transfer.AuthLoginResponse, error) {
	return &transfer.AuthLoginResponse{User: &transfer.AuthUserResponse{ID: "u1", Username: req.Username}}, nil
}
func (f *fakeAuthServicePrivate) Logout(ctx context.Context, sessionID string) error { return nil }
func (f *fakeAuthServicePrivate) Refresh(ctx context.Context, userID, refreshToken string) (*transfer.AuthLoginResponse, error) {
	return &transfer.AuthLoginResponse{User: &transfer.AuthUserResponse{ID: userID}}, nil
}
func (f *fakeAuthServicePrivate) Me(ctx context.Context, userID string) (*transfer.AuthMeResponse, error) {
	return &transfer.AuthMeResponse{User: &transfer.AuthUserResponse{ID: userID}}, nil
}
func (f *fakeAuthServicePrivate) Block(ctx context.Context, userID string, blocked bool) (*transfer.AuthUserResponse, error) {
	return &transfer.AuthUserResponse{ID: userID}, nil
}

func TestAuthControllerPrivateRoutes(t *testing.T) {
	app := fiber.New()
	svc := &fakeAuthServicePrivate{}
	ctrl := NewAuthController(svc)
	grp := app.Group("/auth")

	// add middleware to inject user_id (register before routes so it applies)
	app.Use(func(c *fiber.Ctx) error {
		c.Locals("user_id", "u1")
		return c.Next()
	})

	ctrl.RegisterPrivateRoutes(grp)

	// Test /auth/me
	req := httptest.NewRequest("GET", "/auth/me", nil)
	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("app.Test error: %v", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		t.Fatalf("unexpected status for me: %v", resp.StatusCode)
	}

	// Test /auth/logout (POST with body)
	logoutBody := transfer.AuthLogoutRequest{SessionID: "s1"}
	b, _ := json.Marshal(logoutBody)
	req = httptest.NewRequest("POST", "/auth/logout", bytes.NewReader(b))
	req.Header.Set("Content-Type", "application/json")
	resp, err = app.Test(req)
	if err != nil {
		t.Fatalf("app.Test error: %v", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != 204 && resp.StatusCode != 200 {
		t.Fatalf("unexpected status for logout: %v", resp.StatusCode)
	}

	// Test /auth/block/:id
	blockBody := map[string]bool{"blocked": true}
	bb, _ := json.Marshal(blockBody)
	req = httptest.NewRequest("POST", "/auth/block/u1", bytes.NewReader(bb))
	req.Header.Set("Content-Type", "application/json")
	resp, err = app.Test(req)
	if err != nil {
		t.Fatalf("app.Test error: %v", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		t.Fatalf("unexpected status for block: %v", resp.StatusCode)
	}
}

func TestAuthControllerNegativeCases(t *testing.T) {
	app := fiber.New()
	svc := &fakeAuthServicePrivate{}
	ctrl := NewAuthController(svc)
	grp := app.Group("/auth")
	ctrl.RegisterPublicRoutes(grp)
	ctrl.RegisterPrivateRoutes(grp)

	// 1) Invalid JSON for register
	req := httptest.NewRequest("POST", "/auth/register", bytes.NewReader([]byte("{badjson")))
	req.Header.Set("Content-Type", "application/json")
	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("app.Test error: %v", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != 400 {
		t.Fatalf("expected 400 for bad json, got %d", resp.StatusCode)
	}

	// 2) Validation error: empty login body
	req = httptest.NewRequest("POST", "/auth/login", bytes.NewReader([]byte(`{"username":"","password":""}`)))
	req.Header.Set("Content-Type", "application/json")
	resp, err = app.Test(req)
	if err != nil {
		t.Fatalf("app.Test error: %v", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != 400 {
		t.Fatalf("expected 400 for validation error, got %d", resp.StatusCode)
	}

	// 3) Me without user_id should be unauthorized
	// create a fresh app without user_id injection
	app2 := fiber.New()
	ctrl2 := NewAuthController(svc)
	grp2 := app2.Group("/auth")
	ctrl2.RegisterPrivateRoutes(grp2)

	req = httptest.NewRequest("GET", "/auth/me", nil)
	resp, err = app2.Test(req)
	if err != nil {
		t.Fatalf("app.Test error: %v", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != 401 {
		t.Fatalf("expected 401 for unauthenticated me, got %d", resp.StatusCode)
	}
}
