package cache

import (
    "context"
    "github.com/redis/go-redis/v9"
)

var Ctx = context.Background()

func NewRedisClient(addr, password string, db int) *redis.Client {
    return redis.NewClient(&redis.Options{
	Addr:     addr,
	Password: password,
	DB:       db,
    })
}
