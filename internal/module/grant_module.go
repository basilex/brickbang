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
	svc  service.IGrantService
	ctrl controller.IGrantController
}

func NewGrantModule(dbs *dbs.Queries) IGrantModule {
	svc := service.NewGrantService(dbs)
	ctrl := controller.NewGrantController(svc)

	return &grantModule{svc: svc, ctrl: ctrl}
}

func (rcv *grantModule) Service() service.IGrantService          { return rcv.svc }
func (rcv *grantModule) Controller() controller.IGrantController { return rcv.ctrl }
