package command

import "go-hackaton/src/pkg/scoringtask/model"

type UnitOfWork interface {
	Execute(func(r model.ScoringTaskRepository) error) error
}
