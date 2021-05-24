package input

import (
	"github.com/google/uuid"
	"go-hackathon/src/scoringservice/pkg/scoringtask/api/errors"
	"go-hackathon/src/scoringservice/pkg/scoringtask/application/command"
)

type AddScoringTaskInput struct {
	SolutionID string
	TaskType   int
	Endpoint   string
}

func (i *AddScoringTaskInput) Command() (*command.AddTaskCommand, error) {
	id, err := uuid.Parse(i.SolutionID)
	if err != nil {
		return nil, errors.InvalidSolutionIdError
	}

	return &command.AddTaskCommand{
		SolutionID: id,
		TaskType:   i.TaskType,
		Endpoint:   i.Endpoint,
	}, nil
}

type RemoveScoringTasksInput struct {
	SolutionIDs []string
}

func (i *RemoveScoringTasksInput) Command() (*command.RemoveTasksCommand, error) {
	ids := make([]uuid.UUID, len(i.SolutionIDs))
	for i, stringID := range i.SolutionIDs {
		id, err := uuid.Parse(stringID)
		if err != nil {
			return nil, errors.InvalidSolutionIdError
		}

		ids[i] = id
	}

	return &command.RemoveTasksCommand{SolutionIDs: ids}, nil
}
