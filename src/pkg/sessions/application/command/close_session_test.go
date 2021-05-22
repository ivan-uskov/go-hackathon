package command

import (
	"go-hackaton/src/pkg/sessions/application/errors"
	"testing"
)

func TestCloseSessionWithInvalidId(t *testing.T) {
	h := closeSessionCommandHandler{&mockUnitOfWork{}}
	err := h.Handle(CloseSessionCommand{mockSession.ID})
	if err != errors.SessionNotExistsError {
		t.Error("Close session with invalid id works")
	}
}

func TestCloseSession(t *testing.T) {
	uow := &mockUnitOfWork{}
	_ = uow.mockSessionRepository.Add(mockSession)
	h := closeSessionCommandHandler{uow}
	err := h.Handle(CloseSessionCommand{mockSession.ID})
	if err != nil {
		t.Error("Close session not works")
	}

	s, _ := uow.mockSessionRepository.Get(mockSession.ID)
	if s == nil {
		t.Error("Session not exists after close")
	} else if !s.IsClosed() {
		t.Error("Session is not closed after close")
	}
}

func TestCloseAlreadyClosedSession(t *testing.T) {
	uow := &mockUnitOfWork{}
	session := mockSession
	session.Close()
	_ = uow.mockSessionRepository.Add(session)

	h := closeSessionCommandHandler{uow}
	err := h.Handle(CloseSessionCommand{session.ID})
	if err != errors.SessionAlreadyClosedError {
		t.Error("Close already closed session works")
	}
}
