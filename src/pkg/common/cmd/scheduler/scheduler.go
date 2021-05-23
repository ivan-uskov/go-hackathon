package scheduler

import (
	"time"
)

const cycleTimeout = time.Second * 10

type Job func() bool

type scheduler struct {
	job         Job
	stopChan    chan bool
	stoppedChan chan bool
}

type Api interface {
	Start()
	Stop()
}

func NewScheduler(job Job) Api {
	return &scheduler{
		job:         job,
		stopChan:    make(chan bool, 1),
		stoppedChan: make(chan bool, 1),
	}
}

func (s *scheduler) Start() {
	go s.cycle()
}

func (s *scheduler) Stop() {
	s.stopChan <- true
	<-s.stoppedChan
}

func (s *scheduler) cycle() {
	for !s.needStopCycle() {
		success := s.job()

		if !success {
			select {
			case <-time.After(cycleTimeout):
				break
			case <-s.stopChan:
				return
			}
		}

		select {
		case _, ok := <-s.stoppedChan:
			if ok {
				return
			}
		default:
		}
	}

	s.stoppedChan <- true
}

func (s *scheduler) needStopCycle() bool {
	select {
	case _, ok := <-s.stopChan:
		return ok
	default:
		return false
	}
}
