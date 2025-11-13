package module

import (
	"brickbang/internal/controller"
	"brickbang/internal/repository/dbs"
	"brickbang/internal/service"
)

type IRoleModule interface {
	Service() service.IRoleService
	Controller() controller.IRoleController
}
type roleModule struct {
	svc  service.IRoleService
	ctrl controller.IRoleController
}

func NewRoleModule(dbs *dbs.Queries) IRoleModule {
	svc := service.NewRoleService(dbs)
	ctrl := controller.NewRoleController(svc)

	return &roleModule{
		svc: svc, ctrl: ctrl,
	}
}

func (rcv *roleModule) Service() service.IRoleService          { return rcv.svc }
func (rcv *roleModule) Controller() controller.IRoleController { return rcv.ctrl }
