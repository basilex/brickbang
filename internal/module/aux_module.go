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
	svc  service.IAuxService
	ctrl controller.IAuxController
}

func NewAuxModule(meta map[string]string) IAuxModule {
	svc := service.NewAuxService(meta)
	ctrl := controller.NewAuxController(svc)

	return &auxModule{svc: svc, ctrl: ctrl}
}

func (rcv *auxModule) Service() service.IAuxService          { return rcv.svc }
func (rcv *auxModule) Controller() controller.IAuxController { return rcv.ctrl }
