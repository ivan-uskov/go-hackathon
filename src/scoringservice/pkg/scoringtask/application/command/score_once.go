package command

import (
	"go-hackathon/src/scoringservice/pkg/scoringtask/application/adapter"
	"go-hackathon/src/scoringservice/pkg/scoringtask/application/errors"
	"go-hackathon/src/scoringservice/pkg/scoringtask/model"
	"time"
)

const scoreTimeout = 5 * time.Second

type ScoreOnceCommandHandler interface {
	Handle() error
}

type scoreOnceCommandHandler struct {
	uow     UnitOfWork
	factory adapter.ScorerFactory
}

func NewScoreOnceCommandHandler(uow UnitOfWork, factory adapter.ScorerFactory) ScoreOnceCommandHandler {
	return &scoreOnceCommandHandler{uow, factory}
}

func (s *scoreOnceCommandHandler) Handle() error {
	var task *model.ScoringTask
	err := s.uow.Execute(func(rp RepositoryProvider) (err error) {
		task, err = rp.ScoringTaskRepository().GetFirstScoringTaskBefore(time.Now().Add(-scoreTimeout))
		return err
	})
	if err != nil {
		return err
	}
	if task == nil {
		return errors.TaskNotExistError
	}

	scorer, err := s.factory.GetScorer(task.Type)
	if err != nil {
		return err
	}

	task.UpdateScore(scorer.Score(task.Endpoint))

	return s.uow.Execute(func(rp RepositoryProvider) error {
		err := rp.ScoringTaskRepository().Add(*task)
		if err != nil {
			return err
		}

		return rp.EventStore().Add(&model.ScoringCompleteEvent{
			SolutionID: task.SolutionID,
			Score:      task.Score,
			ScoredAt:   *task.ScoredAt,
		})
	})
}
