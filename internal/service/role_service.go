package service

import (
	"context"
	"errors"

	"brickbang/internal/repository/dbs"
)

type IRoleService interface {
	Create(ctx context.Context, name string) (*dbs.Role, error)
	GetByID(ctx context.Context, id string) (*dbs.Role, error)
	GetByName(ctx context.Context, name string) (*dbs.Role, error)
	Count(ctx context.Context) (int64, error)
	List(ctx context.Context, order string, limit, offset int32) ([]*dbs.Role, error)
	UpdateByID(ctx context.Context, id string, name string) (*dbs.Role, error)
	DeleteByID(ctx context.Context, id string) (string, error)
}

type roleService struct {
	queries *dbs.Queries
}

func NewRoleService(q *dbs.Queries) IRoleService {
	return &roleService{queries: q}
}

// Create a new role
func (s *roleService) Create(ctx context.Context, name string) (*dbs.Role, error) {
	return s.queries.CreateRole(ctx, name)
}

// Get a role by its ID
func (s *roleService) GetByID(ctx context.Context, id string) (*dbs.Role, error) {
	role, err := s.queries.GetRoleByID(ctx, id)
	if err != nil {
		return nil, errors.New("role not found")
	}
	return role, nil
}

// Get a role by its name
func (s *roleService) GetByName(ctx context.Context, name string) (*dbs.Role, error) {
	role, err := s.queries.GetRoleByName(ctx, name)
	if err != nil {
		return nil, errors.New("role not found")
	}
	return role, nil
}

// Count total number of roles
func (s *roleService) Count(ctx context.Context) (int64, error) {
	return s.queries.CountRoles(ctx)
}

// List roles with pagination and order
func (s *roleService) List(ctx context.Context, order string, limit, offset int32) ([]*dbs.Role, error) {
	params := &dbs.ListRolesParams{
		Sqlorder:  order,
		Sqllimit:  limit,
		Sqloffset: offset,
	}
	return s.queries.ListRoles(ctx, params)
}

// Update role name by ID
func (s *roleService) UpdateByID(ctx context.Context, id string, name string) (*dbs.Role, error) {
	params := &dbs.UpdateRoleByIDParams{
		ID:   id,
		Name: name,
	}
	return s.queries.UpdateRoleByID(ctx, params)
}

// Delete role by ID
func (s *roleService) DeleteByID(ctx context.Context, id string) (string, error) {
	role, err := s.queries.DeleteRoleByID(ctx, id)
	if err != nil {
		return "", err
	}
	return role.ID, nil
}
