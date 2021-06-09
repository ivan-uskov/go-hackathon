package adapter

import (
	"context"
	scoring "go-hackathon/api/scoringservice"
	"go-hackathon/src/common/infrastructure"
	"go-hackathon/src/hackathonservice/pkg/hackathon/application/adapter"
	"go-hackathon/src/hackathonservice/pkg/hackathon/application/errors"
)

type scoringAdapter struct {
	api scoring.ScoringServiceClient
}

func NewScoringAdapter(api scoring.ScoringServiceClient) adapter.ScoringAdapter {
	return &scoringAdapter{api: api}
}

func (s *scoringAdapter) AddTask(solutionID string, taskType string, endpoint string) error {
	v, ok := scoring.AddTaskRequest_TaskType_value[taskType]
	if !ok {
		return errors.InvalidHackathonTypeError
	}

	taskTypeObj := scoring.AddTaskRequest_TaskType(v)
	_, err := s.api.AddTask(context.Background(), &scoring.AddTaskRequest{SolutionId: solutionID, TaskType: taskTypeObj, Endpoint: endpoint})
	if err != nil {
		return infrastructure.InternalError(err)
	}

	return nil
}

func (s *scoringAdapter) RemoveTasks(solutionIDs []string) error {
	_, err := s.api.RemoveTasks(context.Background(), &scoring.RemoveTasksRequest{SolutionIds: solutionIDs})
	if err != nil {
		return infrastructure.InternalError(err)
	}

	return nil
}

func (s *scoringAdapter) ValidateTaskType(taskType string) bool {
	_, ok := scoring.AddTaskRequest_TaskType_value[taskType]
	return ok
}
