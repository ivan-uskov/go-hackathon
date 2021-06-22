package repository

import (
	"database/sql"
	"go-hackathon/src/common/infrastructure/repository"
	"go-hackathon/src/hackathonservice/pkg/hackathon/application/command"
	"go-hackathon/src/hackathonservice/pkg/hackathon/model"
)

type unitOfWork struct {
	db repository.Database
}

func NewUnitOfWork(db *sql.DB) command.UnitOfWork {
	return &unitOfWork{db: repository.Database{DB: db}}
}

func (u *unitOfWork) Execute(job command.Job) error {
	return u.db.Tx(func(tx *sql.Tx) error {
		return job(&repositoryProvider{tx})
	})
}

func (u *unitOfWork) WithLock(name string, job command.Job) command.Job {
	return func(rp command.RepositoryProvider) error {
		return u.db.Lock(name, func() error {
			return job(rp)
		})
	}
}

type repositoryProvider struct {
	tx *sql.Tx
}

func (rp *repositoryProvider) HackathonRepository() model.HackathonRepository {
	return &hackathonRepository{rp.tx}
}

func (rp *repositoryProvider) ParticipantRepository() model.ParticipantRepository {
	return &participantRepository{rp.tx}
}
