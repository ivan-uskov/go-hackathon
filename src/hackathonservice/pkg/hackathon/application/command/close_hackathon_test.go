package command

import (
	"go-hackathon/src/hackathonservice/pkg/hackathon/application/errors"
	"testing"
)

func TestCloseHackathonWithInvalidID(t *testing.T) {
	h := closeHackathonCommandHandler{&mockUnitOfWork{}, mockScoringAdapter{}}
	err := h.Handle(CloseHackathonCommand{mockHackathon.ID})
	if err != errors.HackathonNotExistsError {
		t.Error("Close hackathon with invalid id works")
	}
}

func TestCloseHackathon(t *testing.T) {
	uow := &mockUnitOfWork{}
	_ = uow.mockHackathonRepository.Add(mockHackathon)
	hh := closeHackathonCommandHandler{uow, mockScoringAdapter{}}
	err := hh.Handle(CloseHackathonCommand{mockHackathon.ID})
	if err != nil {
		t.Error("Close hackathon not works")
	}

	h, _ := uow.mockHackathonRepository.Get(mockHackathon.ID)
	if h == nil {
		t.Error("Hackathon not exists after close")
	} else if !h.IsClosed() {
		t.Error("Hackathon is not closed after close")
	}
}

func TestCloseAlreadyClosedHackathon(t *testing.T) {
	uow := &mockUnitOfWork{}
	hackathon := mockHackathon
	hackathon.Close()
	_ = uow.mockHackathonRepository.Add(hackathon)

	h := closeHackathonCommandHandler{uow, mockScoringAdapter{}}
	err := h.Handle(CloseHackathonCommand{hackathon.ID})
	if err != errors.HackathonAlreadyClosedError {
		t.Error("Close already closed hackathon works")
	}
}
