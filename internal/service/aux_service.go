package service

import (
	"time"

	"brickbang/internal/repository"
)

type IAuxService interface {
	Health() map[string]string
	Uptime() map[string]string
	Metadata() map[string]string
}

type auxService struct {
	meta map[string]string
	repo repository.IAuxRepository
}

func NewAuxService(meta map[string]string, repo repository.IAuxRepository) IAuxService {
	return &auxService{meta: meta, repo: repo}
}

func (rcv *auxService) Health() map[string]string {
	return map[string]string{"status": "ok"}
}

func (rcv *auxService) Uptime() map[string]string {
	uptime := time.Since(rcv.repo.StartTime()).String()
	return map[string]string{"uptime": uptime}
}

func (rcv *auxService) Metadata() map[string]string {
	return rcv.meta
}
