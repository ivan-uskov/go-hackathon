package command

import (
	"github.com/google/uuid"
	"go-hackaton/src/pkg/sessions/model"
	"time"
)

type UpdateParticipantScoreCommand struct {
	ID    uuid.UUID
	Score int
}

type UpdateParticipantScoreCommandHandler interface {
	Handle(command UpdateParticipantScoreCommand) error
}

type updateParticipantScoreCommandHandler struct {
	repo model.ParticipantRepository
}

func NewUpdateParticipantScoreCommandHandler(repo model.ParticipantRepository) UpdateParticipantScoreCommandHandler {
	return &updateParticipantScoreCommandHandler{repo}
}

func (h updateParticipantScoreCommandHandler) Handle(command UpdateParticipantScoreCommand) error {
	part, err := h.repo.Get(command.ID)
	if err != nil {
		return err
	}

	now := time.Now()
	part.Score = command.Score
	part.ScoredAt = &now
	return h.repo.Add(*part)
}
