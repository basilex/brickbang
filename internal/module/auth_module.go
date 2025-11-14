package module

import (
	"brickbang/internal/controller"
	"brickbang/internal/repository/dbs"
	"brickbang/internal/service"
)

type IAuthModule interface {
	Service() service.IAuthService
	Controller() controller.IAuthController
}
type authModule struct {
	svc  service.IAuthService
	ctrl controller.IAuthController
}

func NewAuthModule(queries *dbs.Queries) IAuthModule {
	svc := service.NewAuthService(queries)
	ctrl := controller.NewAuthController(svc)

	return &authModule{
		svc: svc, ctrl: ctrl,
	}
}

func (rcv *authModule) Service() service.IAuthService {
	return rcv.svc
}

func (rcv *authModule) Controller() controller.IAuthController {
	return rcv.ctrl
}
