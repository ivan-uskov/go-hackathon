package command

import (
	"github.com/google/uuid"
	"go-hackathon/src/hackathonservice/pkg/hackathon/application/errors"
)

type CloseHackathonCommand struct {
	HackathonID uuid.UUID
}

type closeHackathonCommandHandler struct {
	unitOfWork UnitOfWork
}

type CloseHackathonCommandHandler interface {
	Handle(command CloseHackathonCommand) error
}

func NewCloseHackathonCommandHandler(unitOfWork UnitOfWork) CloseHackathonCommandHandler {
	return &closeHackathonCommandHandler{unitOfWork}
}

func (h *closeHackathonCommandHandler) Handle(c CloseHackathonCommand) error {
	return h.unitOfWork.Execute(func(rp RepositoryProvider) error {
		repo := rp.HackathonRepository()

		s, err := repo.Get(c.HackathonID)
		if err != nil {
			return err
		}
		if s == nil {
			return errors.HackathonNotExistsError
		}

		if s.IsClosed() {
			return errors.HackathonAlreadyClosedError
		}

		s.Close()

		return repo.Add(*s)
	})
}
