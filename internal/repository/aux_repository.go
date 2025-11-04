package repository

import "time"

type IAuxRepository interface {
	StartTime() time.Time
}

type auxRepository struct {
	startTime time.Time
}

func NewAuxRepository() IAuxRepository {
	return &auxRepository{startTime: time.Now()}
}

func (rcv *auxRepository) StartTime() time.Time {
	return rcv.startTime
}
