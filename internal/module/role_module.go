package module

import (
	"brickbang/internal/repository"
	"brickbang/internal/service"
	"brickbang/internal/controller"
)

// IRoleModule exposes the interfaces for the entity
type IRoleModule interface {
	Repository() repository.IRoleRepository
	Service() service.IRoleService
	Controller() controller.IRoleController
}

// RoleModuleImpl implements IRoleModule
type RoleModuleImpl struct {
	repo repository.IRoleRepository
	svc  service.IRoleService
	ctrl controller.IRoleController
}

// NewRoleModule constructs the module
func NewRoleModule(container *Container) IRoleModule {
	repo := repository.NewRoleRepository(container.DBPool)
	svc := service.NewRoleService(repo)
	ctrl := controller.NewRoleController(svc)

	return &RoleModuleImpl{
		repo:  repo,
		svc:   svc,
		ctrl:  ctrl,
	}
}

func (m *RoleModuleImpl) Repository() repository.IRoleRepository { return m.repo }
func (m *RoleModuleImpl) Service() service.IRoleService       { return m.svc }
func (m *RoleModuleImpl) Controller() controller.IRoleController { return m.ctrl }
