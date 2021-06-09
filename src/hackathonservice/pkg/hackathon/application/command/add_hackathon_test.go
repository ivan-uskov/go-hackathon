package command

import (
	"go-hackathon/src/hackathonservice/pkg/hackathon/application/errors"
	"testing"
)

const mockValidTaskType = "someType"
const mockInvalidTaskType = "invalidTaskType"

type mockScoringAdapter struct{}

func (m mockScoringAdapter) AddTask(solutionID string, taskType string, endpoint string) error {
	return nil
}

func (m mockScoringAdapter) RemoveTasks(solutionIDs []string) error {
	return nil
}

func (m mockScoringAdapter) ValidateTaskType(taskType string) bool {
	return taskType == mockValidTaskType
}

func TestAddHackathonWithEmptyName(t *testing.T) {
	h := addHackathonCommandHandler{&mockUnitOfWork{}, &mockScoringAdapter{}}
	_, err := h.Handle(AddHackathonCommand{
		"",
		mockValidTaskType,
	})
	if err != errors.InvalidHackathonNameError {
		t.Error("Add hackathon with empty name works")
	}
}

func TestAddHackathonWithInvalidTaskType(t *testing.T) {
	h := addHackathonCommandHandler{&mockUnitOfWork{}, &mockScoringAdapter{}}
	_, err := h.Handle(AddHackathonCommand{
		mockHackathon.Name,
		mockInvalidTaskType,
	})
	if err != errors.InvalidHackathonTypeError {
		t.Error("Add hackathon with invalid task type works")
	}
}

func TestAddDuplicateHackathon(t *testing.T) {
	h := addHackathonCommandHandler{&mockUnitOfWork{}, &mockScoringAdapter{}}
	_, err := h.Handle(AddHackathonCommand{
		mockHackathon.Name,
		mockValidTaskType,
	})
	if err != nil {
		t.Error("Add hackathon not works")
	}
	_, err = h.Handle(AddHackathonCommand{
		mockHackathon.Name,
		mockValidTaskType,
	})
	if err != errors.HackathonAlreadyExistsError {
		t.Error("Add duplicate hackathon works")
	}
}
