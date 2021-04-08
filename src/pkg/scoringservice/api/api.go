package api

import (
	log "github.com/sirupsen/logrus"
	expressions "go-hackaton/src/pkg/expressions/api"
	sessions "go-hackaton/src/pkg/sessions/api"
	"go-hackaton/src/pkg/sessions/api/input"
	"time"
)

const scoreTimeout = time.Minute * 2
const cycleTimeout = time.Second * 10

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
		scored := a.doScore()

		if !scored {
			select {
			case <-time.After(cycleTimeout):
				break
			case <-a.stopScoringChan:
				return
			}
		}

		select {
		case _, ok := <-a.stopScoringChan:
			if ok {
				return
			}
		default:
		}
	}

	a.scoringStoppedChan <- true
}

func (a *api) doScore() bool {
	part, err := a.sessionsApi.GetFirstScoredParticipantBefore(time.Now().Add(-scoreTimeout))
	if err != nil {
		log.Error(err)
		return false
	}

	if part == nil {
		return false
	}

	score := a.expressionsApi.Score(part.Endpoint)
	if score < 0 {
		return false
	}

	log.WithFields(log.Fields{"id": part.ID, "score": score}).Info("score participant")

	err = a.sessionsApi.UpdateSessionParticipantScore(input.UpdateSessionParticipantScoreInput{ID: part.ID, Score: score})
	if err != nil {
		log.Error(err)
		return false
	}

	return true
}

func (a *api) needStopCycle() bool {
	select {
	case _, ok := <-a.stopScoringChan:
		return ok
	default:
		return false
	}
}
