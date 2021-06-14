package command

import (
	"go-hackathon/src/scoringservice/pkg/expressions/application/adapter"
	"go-hackathon/src/scoringservice/pkg/expressions/model"
	"math"
)

const epsilon = 0.01

type ScoreCommand struct {
	Endpoint string
}

type ScoreCommandHandler interface {
	Handle(command ScoreCommand) int
}

func NewScoreCommandHandler(factory adapter.ExternalServiceAdapterFactory) ScoreCommandHandler {
	return &scoreCommandHandler{factory: factory}
}

type scoreCommandHandler struct {
	factory adapter.ExternalServiceAdapterFactory
}

func (s *scoreCommandHandler) Handle(command ScoreCommand) int {
	service := s.factory.Create(command.Endpoint)

	if !service.HealthCheck() {
		return 0
	}

	score := model.HealthCheckScore
	for _, expr := range model.GetExpressions() {
		res, err := service.Calculate(expr.Expr)
		if err == nil && s.floatsEquals(res, expr.Result) {
			score += expr.Score
		}
	}

	return score
}

func (s *scoreCommandHandler) floatsEquals(left, right float64) bool {
	return math.Abs(left-right) <= epsilon
}
