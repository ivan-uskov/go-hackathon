package command

import (
	"github.com/google/uuid"
	"go-hackaton/src/pkg/sessions/application/errors"
)

type CloseSessionCommand struct {
	SessionID uuid.UUID
}

type closeSessionCommandHandler struct {
	unitOfWork UnitOfWork
}

type CloseSessionCommandHandler interface {
	Handle(command CloseSessionCommand) error
}

func NewCloseSessionCommandHandler(unitOfWork UnitOfWork) CloseSessionCommandHandler {
	return &closeSessionCommandHandler{unitOfWork}
}

func (h *closeSessionCommandHandler) Handle(c CloseSessionCommand) error {
	return h.unitOfWork.Execute(func(rp RepositoryProvider) error {
		sessRepo := rp.SessionRepository()

		s, err := sessRepo.Get(c.SessionID)
		if err != nil {
			return err
		}
		if s == nil {
			return errors.SessionNotExistsError
		}

		if s.IsClosed() {
			return errors.SessionAlreadyClosedError
		}

		s.Close()

		return sessRepo.Add(*s)
	})
}
