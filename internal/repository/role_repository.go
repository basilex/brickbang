package repository

import (
    "context"

    "brickbang/storage/dbs"
)

type IRoleRepository interface {
    Create(ctx context.Context, params *dbs.) (*dbs., error)
    List(ctx context.Context, params *dbs.ListRolesParams) ([]*dbs., error)
    GetByID(ctx context.Context, id string) (*dbs., error)
    UpdateByID(ctx context.Context, params *dbs.UpdateRoleByIDParams) (*dbs., error)
    DeleteByID(ctx context.Context, id string) error
    Count(ctx context.Context) (int64, error)
}

type RoleRepository struct {
    q *dbs.Queries
}

func NewRoleRepository(q *dbs.Queries) IRoleRepository {
    return &RoleRepository{q: q}
}

func (r *RoleRepository) Create(ctx context.Context, params *dbs.) (*dbs., error) {
    return r.q.CreateRole(ctx, *params)
}

func (r *RoleRepository) List(ctx context.Context, params *dbs.ListRolesParams) ([]*dbs., error) {
    return r.q.ListRoles(ctx, *params)
}

func (r *RoleRepository) GetByID(ctx context.Context, id string) (*dbs., error) {
    return r.q.GetRoleByID(ctx, id)
}

func (r *RoleRepository) UpdateByID(ctx context.Context, params *dbs.UpdateRoleByIDParams) (*dbs., error) {
    return r.q.UpdateRoleByID(ctx, *params)
}

func (r *RoleRepository) DeleteByID(ctx context.Context, id string) error {
    return r.q.DeleteRoleByID(ctx, id)
}

func (r *RoleRepository) Count(ctx context.Context) (int64, error) {
    return r.q.CountRoles(ctx)
}
