package command

import (
	"go-hackathon/src/scoringservice/pkg/scoringtask/application/errors"
	"testing"
)

func TestAddDuplicateScoringTask(t *testing.T) {
	h := NewAddTaskCommandHandler(&mockUnitOfWork{})
	err := h.Handle(AddTaskCommand{
		mockScoringTask.SolutionID,
		mockScoringTask.Type,
		mockScoringTask.Endpoint,
	})
	if err != nil {
		t.Error("Add scoring task not works")
	}

	err = h.Handle(AddTaskCommand{
		mockScoringTask.SolutionID,
		mockScoringTask.Type,
		mockScoringTask.Endpoint,
	})
	if err != errors.TaskAlreadyExistError {
		t.Error("Add duplicate scoring task works")
	}
}

func TestAddScoringTaskWithInvalidType(t *testing.T) {
	h := NewAddTaskCommandHandler(&mockUnitOfWork{})
	err := h.Handle(AddTaskCommand{
		mockScoringTask.SolutionID,
		"invalid type",
		mockScoringTask.Endpoint,
	})
	if err != errors.InvalidTaskTypeError {
		t.Error("Add scoring task with invalid type works")
	}
}
