package supervisor

import (
	"time"
)

type Job func() error

type supervisor struct {
	failTimeout time.Duration
	job         Job
	stopChan    chan struct{}
	stoppedChan chan struct{}
}

type Supervisor interface {
	Stop()
}

func StartSupervisor(job Job, failTimeoutSeconds int) Supervisor {
	s := &supervisor{
		failTimeout: time.Duration(failTimeoutSeconds) * time.Second,
		job:         job,
		stopChan:    make(chan struct{}, 1),
		stoppedChan: make(chan struct{}, 1),
	}

	go s.cycle()

	return s
}

func (s *supervisor) Stop() {
	s.stopChan <- struct{}{}
	<-s.stoppedChan
}

func (s *supervisor) cycle() {
	defer func() {
		s.stoppedChan <- struct{}{}
	}()

	for {
		select {
		case <-s.stopChan:
			return
		default:
		}

		err := s.job()

		if err != nil {
			select {
			case <-time.After(s.failTimeout):
				break
			case <-s.stopChan:
				return
			}
		}
	}
}
