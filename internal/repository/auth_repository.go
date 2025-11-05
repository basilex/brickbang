package repository

import (
	"context"

	"brickbang/storage/dbs"
)

type IAuthRepository interface {
	FindByUsername(ctx context.Context, username string) (*dbs.AuthSelectUserCredentialsRow, error)
	FindByID(ctx context.Context, id string) (*dbs.AuthSelectUserByIDRow, error)
	CreateUser(ctx context.Context, user *dbs.AuthCreateUserParams) (*dbs.AuthCreateUserRow, error)
	UpdateVisitedAt(ctx context.Context, id string) (*dbs.AuthUpdateVisitedAtRow, error)
}
type AuthRepository struct {
	db *dbs.Queries
}

func NewAuthRepository(db *dbs.Queries) IAuthRepository {
	return &AuthRepository{db: db}
}

func (r *AuthRepository) FindByUsername(ctx context.Context, username string) (*dbs.AuthSelectUserCredentialsRow, error) {
	return r.db.AuthSelectUserCredentials(ctx, username)
}

func (r *AuthRepository) FindByID(ctx context.Context, id string) (*dbs.AuthSelectUserByIDRow, error) {
	return r.db.AuthSelectUserByID(ctx, id)
}

func (r *AuthRepository) CreateUser(ctx context.Context, user *dbs.AuthCreateUserParams) (*dbs.AuthCreateUserRow, error) {
	return r.db.AuthCreateUser(ctx, user)
}

func (r *AuthRepository) UpdateVisitedAt(ctx context.Context, id string) (*dbs.AuthUpdateVisitedAtRow, error) {
	return r.db.AuthUpdateVisitedAt(ctx, id)
}
