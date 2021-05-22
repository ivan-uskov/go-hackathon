package command

import (
	"go-hackaton/src/pkg/sessions/application/errors"
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

func TestAddSessionWithEmptyCode(t *testing.T) {
	h := addSessionCommandHandler{&mockUnitOfWork{}, &mockTasksAdapter{}}
	_, err := h.Handle(AddSessionCommand{
		"",
		mockSession.Name,
		mockValidTaskType,
	})
	if err != errors.InvalidSessionCodeError {
		t.Error("Add session with empty code works")
	}
}

func TestAddSessionWithEmptyName(t *testing.T) {
	h := addSessionCommandHandler{&mockUnitOfWork{}, &mockTasksAdapter{}}
	_, err := h.Handle(AddSessionCommand{
		mockSession.Code,
		"",
		mockValidTaskType,
	})
	if err != errors.InvalidSessionNameError {
		t.Error("Add session with empty name works")
	}
}

func TestAddSessionWithInvalidTaskType(t *testing.T) {
	h := addSessionCommandHandler{&mockUnitOfWork{}, &mockTasksAdapter{}}
	_, err := h.Handle(AddSessionCommand{
		mockSession.Code,
		mockSession.Name,
		mockInvalidTaskType,
	})
	if err != errors.InvalidTaskTypeError {
		t.Error("Add session with invalid task type works")
	}
}

func TestAddDuplicateSession(t *testing.T) {
	h := addSessionCommandHandler{&mockUnitOfWork{}, &mockTasksAdapter{}}
	_, err := h.Handle(AddSessionCommand{
		mockSession.Code,
		mockSession.Name,
		mockValidTaskType,
	})
	if err != nil {
		t.Error("Add session not works")
	}
	_, err = h.Handle(AddSessionCommand{
		mockSession.Code,
		mockSession.Name,
		mockValidTaskType,
	})
	if err != errors.SessionAlreadyExistsError {
		t.Error("Add duplicate session works")
	}
}
