package command

import (
	"go-hackathon/src/common/application/events"
	"go-hackathon/src/scoringservice/pkg/scoringtask/model"
)

type Job func(rp RepositoryProvider) error

type UnitOfWork interface {
	Execute(Job) error
	WithLock(name string, job Job) Job
}

type RepositoryProvider interface {
	ScoringTaskRepository() model.ScoringTaskRepository
	EventStore() events.EventStore
}
