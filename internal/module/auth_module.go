package module

import (
	"brickbang/internal/controller"
	"brickbang/internal/repository/dbs"
	"brickbang/internal/service"
)

type IAuthModule interface {
	AuthService() service.IAuthService
	BlacklistService() service.IBlacklistService
	Controller() controller.IAuthController
}

type authModule struct {
	authSvc      service.IAuthService
	blacklistSvc service.IBlacklistService
	ctrl         controller.IAuthController
}

func NewAuthModule(queries *dbs.Queries, blacklistSvc service.IBlacklistService) IAuthModule {
	authSvc := service.NewAuthService(queries, blacklistSvc)
	ctrl := controller.NewAuthController(authSvc)

	return &authModule{
		authSvc:      authSvc,
		blacklistSvc: blacklistSvc,
		ctrl:         ctrl,
	}
}

func (m *authModule) AuthService() service.IAuthService           { return m.authSvc }
func (m *authModule) BlacklistService() service.IBlacklistService { return m.blacklistSvc }
func (m *authModule) Controller() controller.IAuthController      { return m.ctrl }
