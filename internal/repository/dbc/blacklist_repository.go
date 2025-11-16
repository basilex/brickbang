package dbc

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

const blacklistKeyPrefix = "auth:blacklist:"

type IBlacklistRepository interface {
	Add(ctx context.Context, jti string, ttl time.Duration) error
	Exists(ctx context.Context, jti string) (bool, error)
	Delete(ctx context.Context, jti string) error
}

type blacklistRepository struct {
	rdb *redis.Client
}

func NewBlacklistRepository(rdb *redis.Client) IBlacklistRepository {
	return &blacklistRepository{rdb: rdb}
}

func key(jti string) string {
	return blacklistKeyPrefix + jti
}

func (r *blacklistRepository) Add(ctx context.Context, jti string, ttl time.Duration) error {
	if jti == "" {
		return fmt.Errorf("empty jti")
	}
	return r.rdb.Set(ctx, key(jti), 1, ttl).Err()
}

func (r *blacklistRepository) Exists(ctx context.Context, jti string) (bool, error) {
	if jti == "" {
		return false, fmt.Errorf("empty jti")
	}
	res, err := r.rdb.Exists(ctx, key(jti)).Result()
	return res == 1, err
}

func (r *blacklistRepository) Delete(ctx context.Context, jti string) error {
	if jti == "" {
		return fmt.Errorf("empty jti")
	}
	return r.rdb.Del(ctx, key(jti)).Err()
}
