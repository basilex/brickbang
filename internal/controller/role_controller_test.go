package controller

import (
	"context"
	"encoding/json"
	"net/http/httptest"
	"testing"

	"brickbang/internal/repository/dbs"

	"github.com/gofiber/fiber/v2"
)

type fakeRoleService struct{}

func (f *fakeRoleService) Create(ctx context.Context, name string, description string) (*dbs.Role, error) {
	return &dbs.Role{Pid: "p1", Name: name, Description: description}, nil
}
func (f *fakeRoleService) GetByID(ctx context.Context, pid string) (*dbs.Role, error) {
	return &dbs.Role{Pid: pid, Name: "r"}, nil
}
func (f *fakeRoleService) GetByName(ctx context.Context, name string) (*dbs.Role, error) {
	return &dbs.Role{Pid: "p1", Name: name}, nil
}
func (f *fakeRoleService) Count(ctx context.Context) (int64, error) { return 0, nil }
func (f *fakeRoleService) List(ctx context.Context, order string, limit, offset int32) ([]*dbs.Role, error) {
	return []*dbs.Role{}, nil
}
func (f *fakeRoleService) UpdateByID(ctx context.Context, pid string, name string, description string) (*dbs.Role, error) {
	return &dbs.Role{Pid: pid, Name: name, Description: description}, nil
}
func (f *fakeRoleService) DeleteByID(ctx context.Context, pid string) (string, error) {
	return pid, nil
}

func TestRoleListRoute(t *testing.T) {
	app := fiber.New()
	svc := &fakeRoleService{}
	ctrl := NewRoleController(svc)
	grp := app.Group("/roles")
	ctrl.RegisterRoutes(grp)

	req := httptest.NewRequest("GET", "/roles/", nil)
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
