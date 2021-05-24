package command

import (
	"go-hackathon/src/scoringservice/pkg/scoringtask/model"
)

type UnitOfWork interface {
	Execute(func(r model.ScoringTaskRepository) error) error
}
