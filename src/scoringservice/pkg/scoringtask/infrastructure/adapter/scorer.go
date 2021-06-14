package adapter

import (
	"go-hackathon/src/scoringservice/pkg/expressions/api"
	"go-hackathon/src/scoringservice/pkg/scoringtask/application/adapter"
	"go-hackathon/src/scoringservice/pkg/scoringtask/application/errors"
	"go-hackathon/src/scoringservice/pkg/scoringtask/model"
)

type scorerFactory struct {
	expressionsScorer adapter.Scorer
}

func NewScorerFactory(api api.Api) adapter.ScorerFactory {
	return &scorerFactory{api}
}

func (s *scorerFactory) GetScorer(taskType string) (adapter.Scorer, error) {
	switch taskType {
	case model.TaskTypeArithmeticExpression:
		return s.expressionsScorer, nil
	default:
		return nil, errors.ScorerNotExistError
	}
}
