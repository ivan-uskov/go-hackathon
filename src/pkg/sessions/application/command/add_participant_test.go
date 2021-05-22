package command

import (
	"github.com/google/uuid"
	"go-hackaton/src/pkg/sessions/application/errors"
	"testing"
)

func TestAddParticipantToNotExistentSession(t *testing.T) {
	uow := &mockUnitOfWork{}
	_ = uow.mockSessionRepository.Add(mockSession)

	h := addParticipantCommandHandler{uow}
	err := h.Handle(AddParticipantCommand{
		uuid.New(),
		"Test session",
		mockParticipant.Name,
		mockParticipant.Endpoint,
	})
	if err != errors.SessionNotExistsError {
		t.Error("Add participant to not exists session works")
	}
}

func TestAddParticipantWithInvalidSessionCode(t *testing.T) {
	uow := &mockUnitOfWork{}
	_ = uow.mockSessionRepository.Add(mockSession)

	h := addParticipantCommandHandler{uow}
	err := h.Handle(AddParticipantCommand{
		mockSession.ID,
		"Test session",
		mockParticipant.Name,
		mockParticipant.Endpoint,
	})
	if err != errors.InvalidSessionCodeError {
		t.Error("Add participant to session with invalid code works")
	}
}

func TestAddParticipantToClosedSession(t *testing.T) {
	uow := &mockUnitOfWork{}
	session := mockSession
	session.Close()

	_ = uow.mockSessionRepository.Add(session)
	h := addParticipantCommandHandler{uow}
	err := h.Handle(AddParticipantCommand{
		session.ID,
		session.Code,
		mockParticipant.Name,
		mockParticipant.Endpoint,
	})
	if err != errors.SessionClosedError {
		t.Error("Add participant to closed session works")
	}
}

func TestAddParticipantWithEmptyName(t *testing.T) {
	uow := &mockUnitOfWork{}
	_ = uow.mockSessionRepository.Add(mockSession)

	h := addParticipantCommandHandler{uow}
	err := h.Handle(AddParticipantCommand{
		mockSession.ID,
		mockSession.Code,
		"",
		mockParticipant.Endpoint,
	})
	if err != errors.ParticipantNameIsEmptyError {
		t.Error("Add participant with empty name works")
	}
}

func TestAddParticipantWithEmptyEndpoint(t *testing.T) {
	uow := &mockUnitOfWork{}
	_ = uow.mockSessionRepository.Add(mockSession)

	h := addParticipantCommandHandler{uow}
	err := h.Handle(AddParticipantCommand{
		mockSession.ID,
		mockSession.Code,
		mockParticipant.Name,
		"",
	})
	if err != errors.ParticipantEndpointIsEmptyError {
		t.Error("Add participant with empty endpoint works")
	}
}

func TestAddDuplicateParticipant(t *testing.T) {
	uow := &mockUnitOfWork{}
	_ = uow.mockSessionRepository.Add(mockSession)

	h := addParticipantCommandHandler{uow}
	err := h.Handle(AddParticipantCommand{
		mockSession.ID,
		mockSession.Code,
		mockParticipant.Name,
		mockParticipant.Endpoint,
	})
	if err != nil {
		t.Error("Add participant not works")
	}
	err = h.Handle(AddParticipantCommand{
		mockSession.ID,
		mockSession.Code,
		mockParticipant.Name,
		mockParticipant.Endpoint,
	})
	if err != errors.ParticipantAlreadyExistsError {
		t.Error("Add duplicate participant works")
	}
}
