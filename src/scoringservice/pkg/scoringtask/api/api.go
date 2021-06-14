package api

import (
	"database/sql"
	"go-hackathon/src/scoringservice/pkg/scoringtask/api/errors"
	"go-hackathon/src/scoringservice/pkg/scoringtask/api/input"
	"go-hackathon/src/scoringservice/pkg/scoringtask/application/command"
	"go-hackathon/src/scoringservice/pkg/scoringtask/infrastructure/repository"
)

type Api interface {
	AddTask(input input.AddScoringTaskInput) error
	RemoveTasks(input input.RemoveScoringTasksInput) error

	ScoreOnce() error
}

type api struct {
	unitOfWork command.UnitOfWork
}

func NewApi(db *sql.DB) Api {
	return &api{repository.NewUnitOfWork(db)}
}

func (a *api) AddTask(in input.AddScoringTaskInput) error {
	c, err := in.Command()
	if err != nil {
		return err
	}

	h := command.NewAddTaskCommandHandler(a.unitOfWork)
	return errors.WrapError(h.Handle(*c))
}

func (a *api) RemoveTasks(in input.RemoveScoringTasksInput) error {
	c, err := in.Command()
	if err != nil {
		return err
	}

	h := command.NewRemoveTasksCommandHandler(a.unitOfWork)
	return errors.WrapError(h.Handle(*c))
}

func (a *api) ScoreOnce() error {
	return nil
}
