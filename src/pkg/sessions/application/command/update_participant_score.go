package command

import (
	"github.com/google/uuid"
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

		part, err := repo.Get(command.ID)
		if err != nil {
			return err
		}

		part.UpdateScore(command.Score)

		return repo.Add(*part)
	})
}
