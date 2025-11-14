package module

import (
	"brickbang/internal/controller"
	"brickbang/internal/repository/dbs"
	"brickbang/internal/service"

	"github.com/go-playground/validator/v10"
)

type IRoleModule interface {
	Service() service.IRoleService
	Controller() controller.IRoleController
}
type roleModule struct {
	svc  service.IRoleService
	ctrl controller.IRoleController
}

func NewRoleModule(dbs *dbs.Queries, validator *validator.Validate) IRoleModule {
	svc := service.NewRoleService(dbs)
	ctrl := controller.NewRoleController(svc, validator)

	return &roleModule{
		svc: svc, ctrl: ctrl,
	}
}

func (rcv *roleModule) Service() service.IRoleService          { return rcv.svc }
func (rcv *roleModule) Controller() controller.IRoleController { return rcv.ctrl }
