package repository

import (
	"context"

	"brickbang/storage/dbs"
)

type IKeyRepository interface {
	APIKeysByUser(ctx context.Context, userID string) ([]*dbs.Apikey, error)
	CreateAPIKey(ctx context.Context, key *dbs.AuthCreateAPIKeyParams) (*dbs.Apikey, error)
	DeleteAPIKey(ctx context.Context, keyID string) error
}

type KeyRepository struct {
	db *dbs.Queries
}

func NewKeyRepository(db *dbs.Queries) IKeyRepository {
	return &KeyRepository{db: db}
}

func (rcv *KeyRepository) APIKeysByUser(ctx context.Context, userID string) ([]*dbs.Apikey, error) {
	return rcv.db.AuthSelectAPIKeysByUser(ctx, userID)
}

func (rcv *KeyRepository) CreateAPIKey(ctx context.Context, key *dbs.AuthCreateAPIKeyParams) (*dbs.Apikey, error) {
	return rcv.db.AuthCreateAPIKey(ctx, key)
}

func (rcv *KeyRepository) DeleteAPIKey(ctx context.Context, keyID string) error {
	return rcv.db.AuthDeleteAPIKey(ctx, keyID)
}
