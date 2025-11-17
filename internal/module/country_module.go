package module

import (
	"brickbang/internal/controller"
	"brickbang/internal/repository/dbs"
	"brickbang/internal/service"
)

type ICountryModule interface {
	Service() service.ICountryService
	Controller() controller.ICountryController
}
type countryModule struct {
	svc  service.ICountryService
	ctrl controller.ICountryController
}

func NewCountryModule(dbs *dbs.Queries) ICountryModule {
	svc := service.NewCountryService(dbs)
	ctrl := controller.NewCountryController(svc)

	return &countryModule{
		svc: svc, ctrl: ctrl,
	}
}

func (rcv *countryModule) Service() service.ICountryService          { return rcv.svc }
func (rcv *countryModule) Controller() controller.ICountryController { return rcv.ctrl }
