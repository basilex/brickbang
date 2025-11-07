package service

import (
    "context"

    "brickbang/storage/dbs"
)

// Интерфейс по конвенции I...
type IRoleService interface {
    Create(ctx context.Context, params dbs.) (*dbs.Role, error)
    List(ctx context.Context, params dbs.ListRolesParams) ([]*dbs.Role, error)
    GetByID(ctx context.Context, id string) (*dbs.Role, error)
    GetByName(ctx context.Context, name string) (*dbs.Role, error)
    UpdateByID(ctx context.Context, params dbs.UpdateRoleByIDParams) (*dbs.Role, error)
    DeleteByID(ctx context.Context, id string) error
    Count(ctx context.Context) (int64, error)
}

// Реализация сервиса
type RoleService struct {
    q *dbs.Queries
}

func NewRoleService(q *dbs.Queries) IRoleService {
    return &RoleService{q: q}
}


func (s *RoleService) Create(ctx context.Context, params dbs.) (*dbs.Role, error) {
    return s.q.CreateRole(ctx, params)
}



func (s *RoleService) List(ctx context.Context, params dbs.ListRolesParams) ([]*dbs.Role, error) {
    return s.q.ListRoles(ctx, params)
}



func (s *RoleService) GetByID(ctx context.Context, id string) (*dbs.Role, error) {
    return s.q.GetRoleByID(ctx, id)
}



func (s *RoleService) GetByName(ctx context.Context, name string) (*dbs.Role, error) {
    return s.q.GetRoleByName(ctx, name)
}



func (s *RoleService) UpdateByID(ctx context.Context, params dbs.UpdateRoleByIDParams) (*dbs.Role, error) {
    return s.q.UpdateRoleByID(ctx, params)
}



func (s *RoleService) DeleteByID(ctx context.Context, id string) error {
    return s.q.DeleteRoleByID(ctx, id)
}



func (s *RoleService) Count(ctx context.Context) (int64, error) {
    return s.q.CountRole(ctx)
}

