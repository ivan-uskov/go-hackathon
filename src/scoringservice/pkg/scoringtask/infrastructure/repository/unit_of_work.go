package repository

import (
	"database/sql"
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

func (u *unitOfWork) Execute(job func(rp model.ScoringTaskRepository) error) error {
	return u.db.Tx(func(tx *sql.Tx) error {
		return job(&scoringTaskRepository{tx})
	})
}
