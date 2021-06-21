package command

import (
	"go-hackathon/src/common/application/events"
	"go-hackathon/src/scoringservice/pkg/scoringtask/model"
)

type UnitOfWork interface {
	Execute(func(rp RepositoryProvider) error) error
}

type RepositoryProvider interface {
	ScoringTaskRepository() model.ScoringTaskRepository
	EventStore() events.EventStore
}
