package module

import (
	"brickbang/internal/controller"
	"brickbang/internal/repository/dbs"
	"brickbang/internal/service"
)

type IGrantModule interface {
	Service() service.IGrantService
	Controller() controller.IGrantController
}

type grantModule struct {
	service    service.IGrantService
	controller controller.IGrantController
}

func NewGrantModule(dbs *dbs.Queries) IGrantModule {
	service := service.NewGrantService(dbs)
	controller := controller.NewGrantController(service)

	return &grantModule{service: service, controller: controller}
}

func (rcv *grantModule) Service() service.IGrantService          { return rcv.service }
func (rcv *grantModule) Controller() controller.IGrantController { return rcv.controller }
