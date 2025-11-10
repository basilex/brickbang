package service

import (
	"context"
	"errors"
	"log"

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

type RoleService struct {
	q *dbs.Queries
}

func NewRoleService(q *dbs.Queries) IRoleService {
	return &RoleService{q: q}
}

// Create creates a new role
func (s *RoleService) Create(ctx context.Context, name string) (*dbs.Role, error) {
	role, err := s.q.CreateRole(ctx, name)
	if err != nil {
		log.Printf("failed to create role: %v", err)
		return nil, err
	}
	return role, nil
}

// GetByID returns a role by its ID
func (s *RoleService) GetByID(ctx context.Context, id string) (*dbs.Role, error) {
	role, err := s.q.GetRoleByID(ctx, id)
	if err != nil {
		log.Printf("role not found by id (%s): %v", id, err)
		return nil, errors.New("role not found")
	}
	return role, nil
}

// GetByName returns a role by its name
func (s *RoleService) GetByName(ctx context.Context, name string) (*dbs.Role, error) {
	role, err := s.q.GetRoleByName(ctx, name)
	if err != nil {
		log.Printf("role not found by name (%s): %v", name, err)
		return nil, errors.New("role not found")
	}
	return role, nil
}

// Count returns total number of roles
func (s *RoleService) Count(ctx context.Context) (int64, error) {
	count, err := s.q.CountRoles(ctx)
	if err != nil {
		log.Printf("failed to count roles: %v", err)
		return 0, err
	}
	return count, nil
}

// List returns a paginated list of roles
func (s *RoleService) List(ctx context.Context, order string, limit, offset int32) ([]*dbs.Role, error) {
	roles, err := s.q.ListRoles(ctx, &dbs.ListRolesParams{
		SqlOrder:  order,
		SqlLimit:  limit,
		SqlOffset: offset,
	})
	if err != nil {
		log.Printf("failed to list roles: %v", err)
		return nil, err
	}
	return roles, nil
}

// UpdateByID updates a role by ID
func (s *RoleService) UpdateByID(ctx context.Context, id string, name string) (*dbs.Role, error) {
	role, err := s.q.UpdateRoleByID(ctx, &dbs.UpdateRoleByIDParams{
		ID:   id,
		Name: name,
	})
	if err != nil {
		log.Printf("failed to update role: %v", err)
		return nil, err
	}
	return role, nil
}

// DeleteByID deletes a role by ID
func (s *RoleService) DeleteByID(ctx context.Context, id string) (string, error) {
	deletedID, err := s.q.DeleteRoleByID(ctx, id)
	if err != nil {
		log.Printf("failed to delete role: %v", err)
		return "", err
	}
	return deletedID, nil
}
