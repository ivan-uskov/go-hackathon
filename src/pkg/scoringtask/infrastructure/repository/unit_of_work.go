package repository

import (
	"database/sql"
	"go-hackaton/src/pkg/common/infrastructure/repository"
	"go-hackaton/src/pkg/scoringtask/application/command"
	"go-hackaton/src/pkg/scoringtask/model"
)

type unitOfWork struct {
	db repository.Database
}

func NewUnitOfWork(db *sql.DB) command.UnitOfWork {
	return &unitOfWork{db: repository.Database{DB: db}}
}

func (u *unitOfWork) Execute(job func(rp model.ScoringTaskRepository) error) error {
	return u.db.Tx(func(tx *sql.Tx) error {
		return job(&scoringTaskRepository{tx})
	})
}
