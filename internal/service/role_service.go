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

type roleService struct {
	queries *dbs.Queries
}

func NewRoleService(q *dbs.Queries) IRoleService {
	return &roleService{queries: q}
}

func (rcv *roleService) Create(ctx context.Context, name string) (*dbs.Role, error) {
	role, err := rcv.queries.RoleCreate(ctx, name)
	if err != nil {
		log.Printf("failed to create role: %v", err)
		return nil, err
	}

	return role, nil
}

func (rcv *roleService) GetByID(ctx context.Context, id string) (*dbs.Role, error) {
	role, err := rcv.queries.RoleGetByID(ctx, id)
	if err != nil {
		log.Printf("role not found by id (%s): %v", id, err)
		return nil, errors.New("role not found")
	}

	return role, nil
}

func (rcv *roleService) GetByName(ctx context.Context, name string) (*dbs.Role, error) {
	role, err := rcv.queries.RoleGetByName(ctx, name)
	if err != nil {
		log.Printf("role not found by name (%s): %v", name, err)
		return nil, errors.New("role not found")
	}

	return role, nil
}

func (rcv *roleService) Count(ctx context.Context) (int64, error) {
	count, err := rcv.queries.RolesCount(ctx)
	if err != nil {
		log.Printf("failed to count roles: %v", err)
		return 0, err
	}

	return count, nil
}

func (rcv *roleService) List(ctx context.Context, order string, limit, offset int32) ([]*dbs.Role, error) {
	roles, err := rcv.queries.RolesList(ctx, &dbs.RolesListParams{
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

func (rcv *roleService) UpdateByID(ctx context.Context, id string, name string) (*dbs.Role, error) {
	role, err := rcv.queries.RoleUpdateByID(ctx, &dbs.RoleUpdateByIDParams{
		ID:   id,
		Name: name,
	})
	if err != nil {
		log.Printf("failed to update role: %v", err)
		return nil, err
	}

	return role, nil
}

func (rcv *roleService) DeleteByID(ctx context.Context, id string) (string, error) {
	deletedID, err := rcv.queries.RoleDeleteByID(ctx, id)
	if err != nil {
		log.Printf("failed to delete role: %v", err)
		return "", err
	}

	return deletedID, nil
}
