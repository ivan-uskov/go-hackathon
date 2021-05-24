package command

import (
	"go-hackathon/src/hackathonservice/pkg/hackathon/application/errors"
	"testing"
)

const mockValidTaskType = "someType"
const mockInvalidTaskType = "invalidTaskType"

type mockTasksAdapter struct{}

func (m mockTasksAdapter) TranslateType(t string) (int, bool) {
	if t == mockInvalidTaskType {
		return 0, false
	}
	return 1, true
}

func TestAddHackathonWithEmptyCode(t *testing.T) {
	h := addHackathonCommandHandler{&mockUnitOfWork{}, &mockTasksAdapter{}}
	_, err := h.Handle(AddHackathonCommand{
		"",
		mockHackathon.Name,
		mockValidTaskType,
	})
	if err != errors.InvalidHackathonCodeError {
		t.Error("Add hackathon with empty code works")
	}
}

func TestAddHackathonWithEmptyName(t *testing.T) {
	h := addHackathonCommandHandler{&mockUnitOfWork{}, &mockTasksAdapter{}}
	_, err := h.Handle(AddHackathonCommand{
		mockHackathon.Code,
		"",
		mockValidTaskType,
	})
	if err != errors.InvalidHackathonNameError {
		t.Error("Add hackathon with empty name works")
	}
}

func TestAddHackathonWithInvalidTaskType(t *testing.T) {
	h := addHackathonCommandHandler{&mockUnitOfWork{}, &mockTasksAdapter{}}
	_, err := h.Handle(AddHackathonCommand{
		mockHackathon.Code,
		mockHackathon.Name,
		mockInvalidTaskType,
	})
	if err != errors.InvalidHackathonTypeError {
		t.Error("Add hackathon with invalid task type works")
	}
}

func TestAddDuplicateHackathon(t *testing.T) {
	h := addHackathonCommandHandler{&mockUnitOfWork{}, &mockTasksAdapter{}}
	_, err := h.Handle(AddHackathonCommand{
		mockHackathon.Code,
		mockHackathon.Name,
		mockValidTaskType,
	})
	if err != nil {
		t.Error("Add hackathon not works")
	}
	_, err = h.Handle(AddHackathonCommand{
		mockHackathon.Code,
		mockHackathon.Name,
		mockValidTaskType,
	})
	if err != errors.HackathonAlreadyExistsError {
		t.Error("Add duplicate hackathon works")
	}
}
