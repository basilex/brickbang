package module

import (
	"brickbang/internal/controller"
	"brickbang/internal/repository/dbs"
	"brickbang/internal/service"

	"github.com/go-playground/validator/v10"
)

type IAuthModule interface {
	Service() service.IAuthService
	Controller() controller.IAuthController
}
type authModule struct {
	svc  service.IAuthService
	ctrl controller.IAuthController
}

func NewAuthModule(queries *dbs.Queries, validator *validator.Validate) IAuthModule {
	svc := service.NewAuthService(queries)
	ctrl := controller.NewAuthController(svc, validator)

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
