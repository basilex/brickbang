package service

import (
	"time"

	"brickbang/internal/repository"
)

type IAuxService interface {
	GetHealth() map[string]string
	GetVersion() map[string]string
	GetUptime() map[string]string
}

type auxService struct {
	repo repository.IAuxRepository
}

func NewAuxService(repo repository.IAuxRepository) IAuxService {
	return &auxService{repo: repo}
}

func (rcv *auxService) GetHealth() map[string]string {
	return map[string]string{"status": "ok"}
}

func (rcv *auxService) GetVersion() map[string]string {
	return map[string]string{"version": "1.0.0"}
}

func (rcv *auxService) GetUptime() map[string]string {
	uptime := time.Since(rcv.repo.GetStartTime()).String()
	return map[string]string{"uptime": uptime}
}
