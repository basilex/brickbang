package service

import (
	"context"
	"time"

	"brickbang/internal/repository/dbc"
)

type IBlacklistService interface {
	Block(ctx context.Context, jti string, ttl time.Duration) error
	IsBlocked(ctx context.Context, jti string) (bool, error)
	Delete(ctx context.Context, jti string) error
}

type blacklistService struct {
	repo dbc.IBlacklistRepository
}

func NewBlacklistService(repo dbc.IBlacklistRepository) IBlacklistService {
	return &blacklistService{repo: repo}
}

func (s *blacklistService) Block(ctx context.Context, jti string, ttl time.Duration) error {
	return s.repo.Add(ctx, jti, ttl)
}

func (s *blacklistService) IsBlocked(ctx context.Context, jti string) (bool, error) {
	return s.repo.Exists(ctx, jti)
}

func (s *blacklistService) Delete(ctx context.Context, jti string) error {
	return s.repo.Delete(ctx, jti)
}
