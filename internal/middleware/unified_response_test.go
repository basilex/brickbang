package middleware

import (
	"encoding/json"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"
)

func TestUnifiedResponseWrapsJSON(t *testing.T) {
	app := fiber.New()
	app.Use(UnifiedResponse())
	app.Get("/ok", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{"hello": "world"})
	})

	req := httptest.NewRequest("GET", "/ok", nil)
	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("app.Test failed: %v", err)
	}
	defer resp.Body.Close()

	var body map[string]any
	if err := json.NewDecoder(resp.Body).Decode(&body); err != nil {
		t.Fatalf("decode: %v", err)
	}

	if _, ok := body["content"]; !ok {
		t.Fatalf("expected content key in wrapper, got: %v", body)
	}
	if _, ok := body["metadata"]; !ok {
		t.Fatalf("expected metadata key in wrapper, got: %v", body)
	}
}
