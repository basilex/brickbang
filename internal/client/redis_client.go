package client

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"brickbang/internal/config"

	"github.com/redis/go-redis/v9"
)

type LoggingHook struct{}

func (h LoggingHook) ProcessHook(next redis.ProcessHook) redis.ProcessHook {
	return func(ctx context.Context, cmd redis.Cmder) error {
		slog.Debug("Redis start", "cmd", cmd.String())
		err := next(ctx, cmd)
		slog.Debug("Redis done", "cmd", cmd.String(), "err:", err)
		return err
	}
}

func (h LoggingHook) ProcessPipelineHook(next redis.ProcessPipelineHook) redis.ProcessPipelineHook {
	return func(ctx context.Context, cmds []redis.Cmder) error {
		slog.Debug("Redis pipeline start", "cmds", cmds)
		err := next(ctx, cmds)
		slog.Debug("Redis pipeline done", "cmds", cmds, "err:", err)
		return err
	}
}

func (h LoggingHook) DialHook(next redis.DialHook) redis.DialHook {
	return next
}

type RedisClient struct {
	Client *redis.Client
}

func NewRedisClient(ctx context.Context, cfg *config.Config) (*RedisClient, error) {
	var lastErr error

	if cfg.RedisAddr == "" {
		return nil, fmt.Errorf("redis address cannot be empty")
	}

	client := redis.NewClient(&redis.Options{
		Addr:            cfg.RedisAddr,
		Password:        cfg.RedisPassword,
		DB:              cfg.RedisDatabase,
		DialTimeout:     cfg.RedisDialTimeout,
		ReadTimeout:     cfg.RedisReadTimeout,
		WriteTimeout:    cfg.RedisWriteTimeout,
		MinRetryBackoff: cfg.RedisMinRetryBackoff,
		MaxRetryBackoff: cfg.RedisMaxRetryBackoff,
	})

	client.AddHook(LoggingHook{})

	for i := 0; i < cfg.RedisReconnectAttempts; i++ {
		if err := client.Ping(ctx).Err(); err != nil {
			lastErr = err
			slog.Warn("Redis ping failed, retrying...", "attempts", i+1, "/", cfg.RedisReconnectAttempts, "err", err)
			time.Sleep(500 * time.Millisecond)
			continue
		}

		lastErr = nil
		break
	}

	if lastErr != nil {
		return nil, fmt.Errorf("failed to connect to redis after %d attempts: %w", cfg.RedisReconnectAttempts, lastErr)
	}

	// slog.Info("Redis connected:", "address", cfg.RedisAddr, "db", cfg.RedisDatabase)
	return &RedisClient{Client: client}, nil
}

func (r *RedisClient) Close() error {
	return r.Client.Close()
}
