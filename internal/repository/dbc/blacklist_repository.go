package dbc

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

const redisBlacklistKey = "blacklist:"

type IBlacklistRepository interface {
	Add(ctx context.Context, jti string, ttl time.Duration) error
	Exists(ctx context.Context, jti string) (bool, error)
}

type blacklistRepository struct {
	client *redis.Client
}

func NewBlacklistRepository(client *redis.Client) IBlacklistRepository {
	return &blacklistRepository{
		client: client,
	}
}

func (rcv *blacklistRepository) Add(ctx context.Context, jti string, ttl time.Duration) error {
	if jti == "" {
		return fmt.Errorf("jti cannot be empty")
	}

	return rcv.client.Set(ctx, "blacklist:"+jti, true, ttl).Err()
}

func (rcv *blacklistRepository) Exists(ctx context.Context, jti string) (bool, error) {
	if jti == "" {
		return false, fmt.Errorf("jti cannot be empty")
	}

	res, err := rcv.client.Exists(ctx, redisBlacklistKey+jti).Result()

	return res > 0, err
}
