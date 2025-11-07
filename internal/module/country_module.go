package module

import (
	"brickbang/internal/repository"
	"brickbang/internal/service"
	"brickbang/internal/controller"
)

// ICountryModule exposes the interfaces for the entity
type ICountryModule interface {
	Repository() repository.ICountryRepository
	Service() service.ICountryService
	Controller() controller.ICountryController
}

// CountryModuleImpl implements ICountryModule
type CountryModuleImpl struct {
	repo repository.ICountryRepository
	svc  service.ICountryService
	ctrl controller.ICountryController
}

// NewCountryModule constructs the module
func NewCountryModule(container *Container) ICountryModule {
	repo := repository.NewCountryRepository(container.DBPool)
	svc := service.NewCountryService(repo)
	ctrl := controller.NewCountryController(svc)

	return &CountryModuleImpl{
		repo:  repo,
		svc:   svc,
		ctrl:  ctrl,
	}
}

func (m *CountryModuleImpl) Repository() repository.ICountryRepository { return m.repo }
func (m *CountryModuleImpl) Service() service.ICountryService       { return m.svc }
func (m *CountryModuleImpl) Controller() controller.ICountryController { return m.ctrl }
