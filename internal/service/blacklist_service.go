package service

import (
	"time"

	"brickbang/internal/repository/dbc"
)

type IBlacklistService interface {
	BlockToken(jti string, ttl time.Duration) error
	IsBlocked(jti string) (bool, error)
}

type blacklistService struct {
	repo dbc.IBlacklistRepository
}

func NewBlacklistService(repo dbc.IBlacklistRepository) IBlacklistService {
	return &blacklistService{repo: repo}
}

func (rcv *blacklistService) BlockToken(jti string, ttl time.Duration) error {
	return rcv.repo.Add(jti, ttl)
}

func (rcv *blacklistService) IsBlocked(jti string) (bool, error) {
	return rcv.repo.Exists(jti)
}
