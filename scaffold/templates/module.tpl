package module

import (
	"{{.ModulePath}}/internal/repository"
	"{{.ModulePath}}/internal/service"
	"{{.ModulePath}}/internal/controller"
)

// I{{.Entity}}Module exposes the interfaces for the entity
type I{{.Entity}}Module interface {
	Repository() repository.I{{.Entity}}Repository
	Service() service.I{{.Entity}}Service
	Controller() controller.I{{.Entity}}Controller
}

// {{.Entity}}ModuleImpl implements I{{.Entity}}Module
type {{.Entity}}ModuleImpl struct {
	repo repository.I{{.Entity}}Repository
	svc  service.I{{.Entity}}Service
	ctrl controller.I{{.Entity}}Controller
}

// New{{.Entity}}Module constructs the module
func New{{.Entity}}Module(container *Container) I{{.Entity}}Module {
	repo := repository.New{{.Entity}}Repository(container.DBPool)
	svc := service.New{{.Entity}}Service(repo)
	ctrl := controller.New{{.Entity}}Controller(svc)

	return &{{.Entity}}ModuleImpl{
		repo:  repo,
		svc:   svc,
		ctrl:  ctrl,
	}
}

func (m *{{.Entity}}ModuleImpl) Repository() repository.I{{.Entity}}Repository { return m.repo }
func (m *{{.Entity}}ModuleImpl) Service() service.I{{.Entity}}Service       { return m.svc }
func (m *{{.Entity}}ModuleImpl) Controller() controller.I{{.Entity}}Controller { return m.ctrl }
