package command

import (
	"github.com/google/uuid"
	"go-hackaton/src/pkg/sessions/application/errors"
	"go-hackaton/src/pkg/sessions/model"
)

type CloseSessionCommand struct {
	SessionID uuid.UUID
}

type closeSessionCommandHandler struct {
	sessRepo model.SessionRepository
}

type CloseSessionCommandHandler interface {
	Handle(command CloseSessionCommand) error
}

func NewCloseSessionCommandHandler(sessRepo model.SessionRepository) CloseSessionCommandHandler {
	return &closeSessionCommandHandler{sessRepo}
}

func (h *closeSessionCommandHandler) Handle(c CloseSessionCommand) error {
	s, err := h.sessRepo.Get(c.SessionID)
	if err != nil {
		return err
	}
	if s == nil {
		return errors.SessionNotExistsError
	}

	s.Close()

	return h.sessRepo.Add(*s)
}
