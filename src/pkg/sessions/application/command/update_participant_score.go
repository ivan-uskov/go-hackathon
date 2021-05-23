package command

import (
	"github.com/google/uuid"
	"go-hackaton/src/pkg/sessions/application/errors"
)

type UpdateParticipantScoreCommand struct {
	ID    uuid.UUID
	Score int
}

type UpdateParticipantScoreCommandHandler interface {
	Handle(command UpdateParticipantScoreCommand) error
}

type updateParticipantScoreCommandHandler struct {
	unitOfWork UnitOfWork
}

func NewUpdateParticipantScoreCommandHandler(unitOfWork UnitOfWork) UpdateParticipantScoreCommandHandler {
	return &updateParticipantScoreCommandHandler{unitOfWork}
}

func (h *updateParticipantScoreCommandHandler) Handle(command UpdateParticipantScoreCommand) error {
	return h.unitOfWork.Execute(func(rp RepositoryProvider) error {
		repo := rp.ParticipantRepository()

		if command.Score < 0 {
			return errors.InvalidParticipantScoreError
		}

		part, err := repo.Get(command.ID)
		if err != nil {
			return err
		}

		part.UpdateScore(command.Score)

		return repo.Add(*part)
	})
}
