package api

import (
	log "github.com/sirupsen/logrus"
	expressions "go-hackaton/src/pkg/expressions/api"
	sessions "go-hackaton/src/pkg/sessions/api"
	"go-hackaton/src/pkg/sessions/api/input"
	"time"
)

const scoreTimeout = time.Minute * 3

type Api interface {
	StartScoring()
	StopScoring()
}

type api struct {
	sessionsApi        sessions.Api
	expressionsApi     expressions.Api
	stopScoringChan    chan bool
	scoringStoppedChan chan bool
}

func NewApi(sessionsApi sessions.Api, expressionsApi expressions.Api) Api {
	return &api{
		sessionsApi:        sessionsApi,
		expressionsApi:     expressionsApi,
		stopScoringChan:    make(chan bool, 1),
		scoringStoppedChan: make(chan bool, 1),
	}
}

func (a *api) StartScoring() {
	go a.scoreCycle()
}

func (a *api) StopScoring() {
	a.stopScoringChan <- true
	<-a.scoringStoppedChan
}

func (a *api) scoreCycle() {
	for !a.needStopCycle() {
		a.doScore()
	}

	a.scoringStoppedChan <- true
}

func (a *api) doScore() {
	part, err := a.sessionsApi.GetFirstScoredParticipantBefore(time.Now().Add(-scoreTimeout))
	if err != nil {
		log.Error(err)
		return
	}

	if part == nil {
		return
	}

	score := a.expressionsApi.Score(part.Endpoint)
	if score > 0 {
		err = a.sessionsApi.UpdateSessionParticipantScore(input.UpdateSessionParticipantScoreInput{ID: part.ID, Score: score})
		if err != nil {
			log.Error(err)
		}
	}
}

func (a *api) needStopCycle() bool {
	select {
	case _, ok := <-a.stopScoringChan:
		return ok
	default:
		return false
	}
}
