package repository

import (
	"database/sql"
	"go-hackaton/src/pkg/common/infrastructure/repository"
	"go-hackaton/src/pkg/sessions/application/command"
	"go-hackaton/src/pkg/sessions/model"
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

func (rp *repositoryProvider) SessionRepository() model.SessionRepository {
	return &sessionRepository{rp.tx}
}

func (rp *repositoryProvider) ParticipantRepository() model.ParticipantRepository {
	return &participantRepository{rp.tx}
}
