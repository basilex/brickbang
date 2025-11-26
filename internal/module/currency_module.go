package module

import (
	"brickbang/internal/controller"
	"brickbang/internal/repository/dbs"
	"brickbang/internal/service"
)

type ICurrencyModule interface {
	Service() service.ICurrencyService
	Controller() controller.ICurrencyController
}
type currencyModule struct {
	service    service.ICurrencyService
	controller controller.ICurrencyController
}

func NewCurrencyModule(dbs *dbs.Queries) ICurrencyModule {
	service := service.NewCurrencyService(dbs)
	controller := controller.NewCurrencyController(service)

	return &currencyModule{service: service, controller: controller}
}

func (rcv *currencyModule) Service() service.ICurrencyService          { return rcv.service }
func (rcv *currencyModule) Controller() controller.ICurrencyController { return rcv.controller }
