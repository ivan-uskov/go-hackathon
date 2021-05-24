package command

import (
	"github.com/google/uuid"
	"go-hackathon/src/scoringservice/pkg/scoringtask/application/errors"
	"testing"
)

func TestRemoveNotExistentScoringTasks(t *testing.T) {
	h := NewRemoveTasksCommandHandler(&mockUnitOfWork{})
	err := h.Handle(RemoveTasksCommand{[]uuid.UUID{mockScoringTask.SolutionID}})
	if err != errors.TaskNotExistError {
		t.Error("Remove not existent scoring task works")
	}
}

func TestRemoveScoringTasks(t *testing.T) {
	uow := &mockUnitOfWork{}
	_ = uow.Add(mockScoringTask)
	h := NewRemoveTasksCommandHandler(uow)

	err := h.Handle(RemoveTasksCommand{[]uuid.UUID{mockScoringTask.SolutionID}})
	if err != nil {
		t.Error("Remove scoring tasks not works")
	}

	task, _ := uow.Get(mockScoringTask.ID)
	if !task.IsDeleted() {
		t.Error("Remove scoring tasks not works")
	}
}
