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
	service    service.ICountryService
	controller controller.ICountryController
}

func NewCountryModule(dbs *dbs.Queries) ICountryModule {
	service := service.NewCountryService(dbs)
	controller := controller.NewCountryController(service)

	return &countryModule{service: service, controller: controller}
}

func (rcv *countryModule) Service() service.ICountryService          { return rcv.service }
func (rcv *countryModule) Controller() controller.ICountryController { return rcv.controller }
