package client

import (
	"context"
	"fmt"

	"github.com/redis/go-redis/v9"
)

// NewRedisClient creates and validates a new Redis client connection.
// Returns an error if the connection cannot be established.
func NewRedisClient(ctx context.Context, addr, password string, db int) (*redis.Client, error) {
	if addr == "" {
		return nil, fmt.Errorf("redis address cannot be empty")
	}

	if db < 0 {
		return nil, fmt.Errorf("redis db must be non-negative")
	}

	client := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
		DB:       db,
	})

	if err := client.Ping(ctx).Err(); err != nil {
		return nil, fmt.Errorf("failed to connect to redis: %w", err)
	}

	return client, nil

}
