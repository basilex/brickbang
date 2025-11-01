package repository

import "time"

type IAuxRepository interface {
	GetStartTime() time.Time
}

type auxRepository struct {
	startTime time.Time
}

func NewAuxRepository() IAuxRepository {
	return &auxRepository{startTime: time.Now()}
}

func (rcv *auxRepository) GetStartTime() time.Time {
	return rcv.startTime
}
