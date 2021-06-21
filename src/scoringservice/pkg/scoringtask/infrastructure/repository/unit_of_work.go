package repository

import (
	"database/sql"
	"go-hackathon/src/common/application/events"
	eventsImpl "go-hackathon/src/common/infrastructure/events"
	"go-hackathon/src/common/infrastructure/repository"
	"go-hackathon/src/scoringservice/pkg/scoringtask/application/command"
	"go-hackathon/src/scoringservice/pkg/scoringtask/model"
)

type unitOfWork struct {
	db repository.Database
}

func NewUnitOfWork(db *sql.DB) command.UnitOfWork {
	return &unitOfWork{db: repository.Database{DB: db}}
}

func (u *unitOfWork) Execute(job func(rp command.RepositoryProvider) error) error {
	return u.db.Tx(func(tx *sql.Tx) error {
		return job(&repositoryProvider{tx})
	})
}

type repositoryProvider struct {
	tx *sql.Tx
}

func (rp *repositoryProvider) ScoringTaskRepository() model.ScoringTaskRepository {
	return &scoringTaskRepository{rp.tx}
}

func (rp *repositoryProvider) EventStore() events.EventStore {
	return eventsImpl.NewEventStore(rp.tx)
}
