package api

import (
	"go-hackathon/src/scoringservice/pkg/expressions/application/adapter"
	"go-hackathon/src/scoringservice/pkg/expressions/application/command"
	adapterImpl "go-hackathon/src/scoringservice/pkg/expressions/infrastructure/adapter"
)

type Api interface {
	Score(url string) int
}

type api struct {
	factory adapter.ExternalServiceAdapterFactory
}

func NewApi() Api {
	return &api{
		adapterImpl.NewExternalServiceAdapterFactory(),
	}
}

func (a *api) Score(url string) int {
	h := command.NewScoreCommandHandler(a.factory)
	return h.Handle(command.ScoreCommand{Endpoint: url})
}
