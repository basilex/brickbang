package controller

import (
	"encoding/json"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"
)

type fakeAuxService struct{}

func (f *fakeAuxService) Health() map[string]string {
	return map[string]string{"status": "ok"}
}
func (f *fakeAuxService) Uptime() map[string]string {
	return map[string]string{"uptime": "42s"}
}
func (f *fakeAuxService) Metadata() map[string]string {
	return map[string]string{"version": "1.2.3"}
}

func TestAuxControllerRoutes(t *testing.T) {
	app := fiber.New()
	svc := &fakeAuxService{}
	ctrl := NewAuxController(svc)
	group := app.Group("/api")
	ctrl.RegisterRoutes(group)

	// test /api/health
	req := httptest.NewRequest("GET", "/api/health", nil)
	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("app.Test error: %v", err)
	}
	defer resp.Body.Close()

	var body map[string]string
	if err := json.NewDecoder(resp.Body).Decode(&body); err != nil {
		t.Fatalf("decode body: %v", err)
	}
	if body["status"] != "ok" {
		t.Fatalf("expected status ok, got %v", body)
	}
}
