package module

import (
	"brickbang/internal/controller"
	"brickbang/internal/repository/dbs"
	"brickbang/internal/service"
)

type IAuthModule interface {
	Controller() controller.IAuthController
	AuthService() service.IAuthService
	BlacklistService() service.IBlacklistService
}

type authModule struct {
	controller       controller.IAuthController
	authService      service.IAuthService
	blacklistService service.IBlacklistService
}

func NewAuthModule(queries *dbs.Queries, blacklistService service.IBlacklistService) IAuthModule {
	authService := service.NewAuthService(queries, blacklistService)
	controller := controller.NewAuthController(authService)

	return &authModule{
		authService: authService, blacklistService: blacklistService, controller: controller,
	}
}

func (rcv *authModule) Controller() controller.IAuthController      { return rcv.controller }
func (rcv *authModule) AuthService() service.IAuthService           { return rcv.authService }
func (rcv *authModule) BlacklistService() service.IBlacklistService { return rcv.blacklistService }
