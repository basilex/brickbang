package service

import (
	"time"
)

type IAuxService interface {
	Health() map[string]string
	Uptime() map[string]string
	Metadata() map[string]string
}
type auxService struct {
	metadata  map[string]string
	startTime time.Time
}

func NewAuxService(meta map[string]string) IAuxService {
	return &auxService{
		metadata:  meta,
		startTime: time.Now(),
	}
}

func (rcv *auxService) Health() map[string]string {
	return map[string]string{"status": "ok"}
}

func (rcv *auxService) Uptime() map[string]string {
	return map[string]string{
		"uptime": time.Since(rcv.startTime).String(),
	}
}

func (rcv *auxService) Metadata() map[string]string {
	return rcv.metadata
}
