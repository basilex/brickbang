package module

import (
	"context"

	"brickbang/internal/controller"
	"brickbang/internal/service"
)

// IAuxModule exposes the interfaces for the entity
type IAuxModule interface {
	Service() service.IAuxService
	Controller() controller.IAuxController
}

type AuxModule struct {
	svc  service.IAuxService
	ctrl controller.IAuxController
}

func NewAuxModule(ctx context.Context, meta map[string]string) IAuxModule {
	svc := service.NewAuxService(meta)
	ctrl := controller.NewAuxController(svc)

	return &AuxModule{
		svc: svc, ctrl: ctrl,
	}
}

func (m *AuxModule) Service() service.IAuxService          { return m.svc }
func (m *AuxModule) Controller() controller.IAuxController { return m.ctrl }
