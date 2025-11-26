package service

import (
	"context"
	"errors"

	"brickbang/internal/repository/dbs"
)

type IRoleService interface {
	Create(ctx context.Context, name string, description string) (*dbs.Role, error)
	GetByID(ctx context.Context, id string) (*dbs.Role, error)
	GetByName(ctx context.Context, name string) (*dbs.Role, error)
	Count(ctx context.Context) (int64, error)
	List(ctx context.Context, order string, limit, offset int32) ([]*dbs.Role, error)
	UpdateByID(ctx context.Context, id string, name string, description string) (*dbs.Role, error)
	DeleteByID(ctx context.Context, id string) (string, error)
}

type roleService struct {
	queries *dbs.Queries
}

func NewRoleService(q *dbs.Queries) IRoleService {
	return &roleService{queries: q}
}

// Create a new role
func (rcv *roleService) Create(ctx context.Context, name string, description string) (*dbs.Role, error) {
	return rcv.queries.CreateRole(ctx, &dbs.CreateRoleParams{
		Name:        name,
		Description: description,
	})
}

// Get a role by its ID
func (rcv *roleService) GetByID(ctx context.Context, id string) (*dbs.Role, error) {
	role, err := rcv.queries.GetRoleByID(ctx, id)
	if err != nil {
		return nil, errors.New("role not found")
	}
	return role, nil
}

// Get a role by its name
func (rcv *roleService) GetByName(ctx context.Context, name string) (*dbs.Role, error) {
	role, err := rcv.queries.GetRoleByName(ctx, name)
	if err != nil {
		return nil, errors.New("role not found")
	}
	return role, nil
}

// Count total number of roles
func (rcv *roleService) Count(ctx context.Context) (int64, error) {
	return rcv.queries.CountRoles(ctx)
}

// List roles with pagination and order
func (rcv *roleService) List(ctx context.Context, order string, limit, offset int32) ([]*dbs.Role, error) {
	params := &dbs.ListRolesParams{
		SqlOrder:  order,
		SqlLimit:  limit,
		SqlOffset: offset,
	}
	return rcv.queries.ListRoles(ctx, params)
}

// Update role name by ID
func (rcv *roleService) UpdateByID(ctx context.Context, id string, name string, description string) (*dbs.Role, error) {
	params := &dbs.UpdateRoleByIDParams{
		ID:          id,
		Name:        name,
		Description: description,
	}
	return rcv.queries.UpdateRoleByID(ctx, params)
}

// Delete role by ID
func (rcv *roleService) DeleteByID(ctx context.Context, id string) (string, error) {
	role, err := rcv.queries.DeleteRoleByID(ctx, id)
	if err != nil {
		return "", err
	}
	return role.ID, nil
}
