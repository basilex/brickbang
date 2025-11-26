package service

import (
	"context"
	"errors"

	"brickbang/internal/repository/dbs"
)

type IGrantService interface {
	Create(ctx context.Context, code string, description string) (*dbs.Grant, error)
	GetByID(ctx context.Context, pid string) (*dbs.Grant, error)
	GetByCode(ctx context.Context, code string) (*dbs.Grant, error)
	Count(ctx context.Context) (int64, error)
	List(ctx context.Context, order string, limit, offset int32) ([]*dbs.Grant, error)
	UpdateByID(ctx context.Context, pid string, code string, description string) (*dbs.Grant, error)
	DeleteByID(ctx context.Context, pid string) (string, error)
}

type grantService struct {
	queries *dbs.Queries
}

func NewGrantService(q *dbs.Queries) IGrantService {
	return &grantService{queries: q}
}

// Create a new grant
func (rcv *grantService) Create(ctx context.Context, code string, description string) (*dbs.Grant, error) {
	return rcv.queries.CreateGrant(ctx, &dbs.CreateGrantParams{
		Code:        code,
		Description: description,
	})
}

// Get a grant by its ID
func (rcv *grantService) GetByID(ctx context.Context, id string) (*dbs.Grant, error) {
	grant, err := rcv.queries.GetGrantByID(ctx, id)
	if err != nil {
		return nil, errors.New("grant not found")
	}
	return grant, nil
}

// Get a grant by its code
func (rcv *grantService) GetByCode(ctx context.Context, code string) (*dbs.Grant, error) {
	grant, err := rcv.queries.GetGrantByCode(ctx, code)
	if err != nil {
		return nil, errors.New("grant not found")
	}
	return grant, nil
}

// Count total number of grants
func (rcv *grantService) Count(ctx context.Context) (int64, error) {
	return rcv.queries.CountGrants(ctx)
}

// List grants with pagination and order
func (rcv *grantService) List(ctx context.Context, order string, limit, offset int32) ([]*dbs.Grant, error) {
	params := &dbs.ListGrantsParams{
		SqlOrder:  order,
		SqlLimit:  limit,
		SqlOffset: offset,
	}
	return rcv.queries.ListGrants(ctx, params)
}

// Update grant name and description by ID
func (rcv *grantService) UpdateByID(ctx context.Context, pid string, code string, description string) (*dbs.Grant, error) {
	params := &dbs.UpdateGrantByIDParams{
		Pid:         pid,
		Code:        code,
		Description: description,
	}
	return rcv.queries.UpdateGrantByID(ctx, params)
}

// Delete grant by ID
func (rcv *grantService) DeleteByID(ctx context.Context, id string) (string, error) {
	grant, err := rcv.queries.DeleteGrantByID(ctx, id)
	if err != nil {
		return "", err
	}
	return grant.ID, nil
}
