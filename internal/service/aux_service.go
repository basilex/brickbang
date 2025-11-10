package service

import (
	"time"
)

type IAuxService interface {
	Health() map[string]string
	Uptime() map[string]string
	Metadata() map[string]string
}
type AuxService struct {
	metadata  map[string]string
	startTime time.Time
}

func NewAuxService(meta map[string]string) IAuxService {
	return &AuxService{
		metadata:  meta,
		startTime: time.Now(),
	}
}

func (rcv *AuxService) Health() map[string]string {
	return map[string]string{"status": "ok"}
}

func (rcv *AuxService) Uptime() map[string]string {
	return map[string]string{
		"uptime": time.Since(rcv.startTime).String(),
	}
}

func (rcv *AuxService) Metadata() map[string]string {
	return rcv.metadata
}
