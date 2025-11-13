package dbc

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

type IBlacklistRepository interface {
	Add(jti string, ttl time.Duration) error
	Exists(jti string) (bool, error)
}

type blacklistRepository struct {
	ctx    context.Context
	client *redis.Client
}

func NewBlacklistRepository(ctx context.Context, client *redis.Client) IBlacklistRepository {
	return &blacklistRepository{
		ctx:    ctx,
		client: client,
	}
}

func (rcv *blacklistRepository) Add(jti string, ttl time.Duration) error {
	return rcv.client.Set(rcv.ctx, "blacklist:"+jti, true, ttl).Err()
}

func (rcv *blacklistRepository) Exists(jti string) (bool, error) {
	res, err := rcv.client.Exists(rcv.ctx, "blacklist:"+jti).Result()
	return res > 0, err
}
