package service

import (
	"context"
	"time"

	"brickbang/internal/repository/dbc"
)

type IBlacklistService interface {
	BlockToken(ctx context.Context, jti string, ttl time.Duration) error
	IsBlocked(ctx context.Context, jti string) (bool, error)
}

type blacklistService struct {
	repo dbc.IBlacklistRepository
}

func NewBlacklistService(repo dbc.IBlacklistRepository) IBlacklistService {
	return &blacklistService{repo: repo}
}

func (rcv *blacklistService) BlockToken(ctx context.Context, jti string, ttl time.Duration) error {
	return rcv.repo.Add(ctx, jti, ttl)
}

func (rcv *blacklistService) IsBlocked(ctx context.Context, jti string) (bool, error) {
	return rcv.repo.Exists(ctx, jti)
}
