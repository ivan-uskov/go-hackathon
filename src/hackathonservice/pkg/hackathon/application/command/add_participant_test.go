package command

import (
	"github.com/google/uuid"
	"go-hackathon/src/hackathonservice/pkg/hackathon/application/errors"
	"testing"
)

func TestAddParticipantToNotExistentHackathon(t *testing.T) {
	uow := &mockUnitOfWork{}
	_ = uow.mockHackathonRepository.Add(mockHackathon)

	h := addParticipantCommandHandler{uow}
	err := h.Handle(AddParticipantCommand{
		uuid.New(),
		"Test hackathon",
		mockParticipant.Name,
		mockParticipant.Endpoint,
	})
	if err != errors.HackathonNotExistsError {
		t.Error("Add participant to not exists hackathon works")
	}
}

func TestAddParticipantWithInvalidHackathonCode(t *testing.T) {
	uow := &mockUnitOfWork{}
	_ = uow.mockHackathonRepository.Add(mockHackathon)

	h := addParticipantCommandHandler{uow}
	err := h.Handle(AddParticipantCommand{
		mockHackathon.ID,
		"Test hackathon",
		mockParticipant.Name,
		mockParticipant.Endpoint,
	})
	if err != errors.InvalidHackathonCodeError {
		t.Error("Add participant to hackathon with invalid code works")
	}
}

func TestAddParticipantToClosedHackathon(t *testing.T) {
	uow := &mockUnitOfWork{}
	hackathon := mockHackathon
	hackathon.Close()

	_ = uow.mockHackathonRepository.Add(hackathon)
	h := addParticipantCommandHandler{uow}
	err := h.Handle(AddParticipantCommand{
		hackathon.ID,
		hackathon.Code,
		mockParticipant.Name,
		mockParticipant.Endpoint,
	})
	if err != errors.HackathonClosedError {
		t.Error("Add participant to closed hackathon works")
	}
}

func TestAddParticipantWithEmptyName(t *testing.T) {
	uow := &mockUnitOfWork{}
	_ = uow.mockHackathonRepository.Add(mockHackathon)

	h := addParticipantCommandHandler{uow}
	err := h.Handle(AddParticipantCommand{
		mockHackathon.ID,
		mockHackathon.Code,
		"",
		mockParticipant.Endpoint,
	})
	if err != errors.ParticipantNameIsEmptyError {
		t.Error("Add participant with empty name works")
	}
}

func TestAddParticipantWithEmptyEndpoint(t *testing.T) {
	uow := &mockUnitOfWork{}
	_ = uow.mockHackathonRepository.Add(mockHackathon)

	h := addParticipantCommandHandler{uow}
	err := h.Handle(AddParticipantCommand{
		mockHackathon.ID,
		mockHackathon.Code,
		mockParticipant.Name,
		"",
	})
	if err != errors.ParticipantEndpointIsEmptyError {
		t.Error("Add participant with empty endpoint works")
	}
}

func TestAddDuplicateParticipant(t *testing.T) {
	uow := &mockUnitOfWork{}
	_ = uow.mockHackathonRepository.Add(mockHackathon)

	h := addParticipantCommandHandler{uow}
	err := h.Handle(AddParticipantCommand{
		mockHackathon.ID,
		mockHackathon.Code,
		mockParticipant.Name,
		mockParticipant.Endpoint,
	})
	if err != nil {
		t.Error("Add participant not works")
	}
	err = h.Handle(AddParticipantCommand{
		mockHackathon.ID,
		mockHackathon.Code,
		mockParticipant.Name,
		mockParticipant.Endpoint,
	})
	if err != errors.ParticipantAlreadyExistsError {
		t.Error("Add duplicate participant works")
	}
}
