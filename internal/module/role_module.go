package module

import (
	"context"

	"brickbang/internal/controller"
	"brickbang/internal/repository/dbs"
	"brickbang/internal/service"
)

// IRoleModule exposes the interfaces for the entity
type IRoleModule interface {
	Service() service.IRoleService
	Controller() controller.IRoleController
}

type RoleModule struct {
	svc  service.IRoleService
	ctrl controller.IRoleController
}

func NewRoleModule(ctx context.Context, dbs *dbs.Queries) IRoleModule {
	svc := service.NewRoleService(dbs)
	ctrl := controller.NewRoleController(svc)

	return &RoleModule{
		svc: svc, ctrl: ctrl,
	}
}

func (m *RoleModule) Service() service.IRoleService          { return m.svc }
func (m *RoleModule) Controller() controller.IRoleController { return m.ctrl }
