package controller

import (
	"context"
	"encoding/json"
	"net/http/httptest"
	"testing"

	"brickbang/internal/repository/dbs"

	"github.com/gofiber/fiber/v2"
)

type fakeGrantService struct{}

func (f *fakeGrantService) Create(ctx context.Context, code string, description string) (*dbs.Grant, error) {
	return &dbs.Grant{Pid: "p1", Code: code, Description: description}, nil
}
func (f *fakeGrantService) GetByID(ctx context.Context, pid string) (*dbs.Grant, error) {
	return &dbs.Grant{Pid: pid, Code: "c1"}, nil
}
func (f *fakeGrantService) GetByCode(ctx context.Context, code string) (*dbs.Grant, error) {
	return &dbs.Grant{Pid: "p1", Code: code}, nil
}
func (f *fakeGrantService) Count(ctx context.Context) (int64, error) { return 0, nil }
func (f *fakeGrantService) List(ctx context.Context, order string, limit, offset int32) ([]*dbs.Grant, error) {
	return []*dbs.Grant{}, nil
}
func (f *fakeGrantService) UpdateByID(ctx context.Context, pid string, code string, description string) (*dbs.Grant, error) {
	return &dbs.Grant{Pid: pid, Code: code, Description: description}, nil
}
func (f *fakeGrantService) DeleteByID(ctx context.Context, pid string) (string, error) {
	return pid, nil
}

func TestGrantListRoute(t *testing.T) {
	app := fiber.New()
	svc := &fakeGrantService{}
	ctrl := NewGrantController(svc)
	grp := app.Group("/grants")
	ctrl.RegisterRoutes(grp)

	req := httptest.NewRequest("GET", "/grants/", nil)
	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("app.Test error: %v", err)
	}
	defer resp.Body.Close()

	var body []any
	if err := json.NewDecoder(resp.Body).Decode(&body); err != nil {
		t.Fatalf("decode body: %v", err)
	}
}
