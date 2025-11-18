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
	svc  service.ICurrencyService
	ctrl controller.ICurrencyController
}

func NewCurrencyModule(dbs *dbs.Queries) ICurrencyModule {
	svc := service.NewCurrencyService(dbs)
	ctrl := controller.NewCurrencyController(svc)

	return &currencyModule{svc: svc, ctrl: ctrl}
}

func (rcv *currencyModule) Service() service.ICurrencyService          { return rcv.svc }
func (rcv *currencyModule) Controller() controller.ICurrencyController { return rcv.ctrl }
