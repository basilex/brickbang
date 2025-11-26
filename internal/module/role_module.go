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
	service    service.IRoleService
	controller controller.IRoleController
}

func NewRoleModule(dbs *dbs.Queries) IRoleModule {
	service := service.NewRoleService(dbs)
	controller := controller.NewRoleController(service)

	return &roleModule{service: service, controller: controller}
}

func (rcv *roleModule) Service() service.IRoleService          { return rcv.service }
func (rcv *roleModule) Controller() controller.IRoleController { return rcv.controller }
