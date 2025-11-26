package module

import (
	"brickbang/internal/controller"
	"brickbang/internal/service"
)

type IAuxModule interface {
	Service() service.IAuxService
	Controller() controller.IAuxController
}

type auxModule struct {
	service    service.IAuxService
	controller controller.IAuxController
}

func NewAuxModule(meta map[string]string) IAuxModule {
	service := service.NewAuxService(meta)
	controller := controller.NewAuxController(service)

	return &auxModule{service: service, controller: controller}
}

func (rcv *auxModule) Service() service.IAuxService          { return rcv.service }
func (rcv *auxModule) Controller() controller.IAuxController { return rcv.controller }
