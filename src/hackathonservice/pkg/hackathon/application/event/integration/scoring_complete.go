package integration

import (
	"github.com/google/uuid"
	"go-hackathon/src/hackathonservice/pkg/hackathon/application/command"
)

type ScoringCompleteEventHandler interface {
	Handle(SolutionID uuid.UUID, Score int) error
}

type scoringCompleteEventHandler struct {
	uow command.UnitOfWork
}

func NewScoringCompleteEventHandler(uow command.UnitOfWork) ScoringCompleteEventHandler {
	return &scoringCompleteEventHandler{uow: uow}
}

func (s *scoringCompleteEventHandler) Handle(SolutionID uuid.UUID, Score int) error {
	h := command.NewUpdateParticipantScoreCommandHandler(s.uow)
	return h.Handle(command.UpdateParticipantScoreCommand{ID: SolutionID, Score: Score})
}
