package command

import (
	"github.com/google/uuid"
	"go-hackathon/src/hackathonservice/pkg/hackathon/application/adapter"
	"go-hackathon/src/hackathonservice/pkg/hackathon/application/errors"
	"go-hackathon/src/hackathonservice/pkg/hackathon/model"
)

type CloseHackathonCommand struct {
	HackathonID uuid.UUID
}

type closeHackathonCommandHandler struct {
	unitOfWork UnitOfWork
	scoring    adapter.ScoringAdapter
}

type CloseHackathonCommandHandler interface {
	Handle(command CloseHackathonCommand) error
}

func NewCloseHackathonCommandHandler(unitOfWork UnitOfWork, scoring adapter.ScoringAdapter) CloseHackathonCommandHandler {
	return &closeHackathonCommandHandler{unitOfWork, scoring}
}

func (h *closeHackathonCommandHandler) Handle(c CloseHackathonCommand) error {
	var participants []model.Participant
	err := h.unitOfWork.Execute(func(rp RepositoryProvider) error {
		repo := rp.HackathonRepository()
		partRepo := rp.ParticipantRepository()

		hackathon, err := repo.Get(c.HackathonID)
		if err != nil {
			return err
		}
		if hackathon == nil {
			return errors.HackathonNotExistsError
		}

		if hackathon.IsClosed() {
			return errors.HackathonAlreadyClosedError
		}

		participants, err = partRepo.GetByHackathonID(hackathon.ID)
		if err != nil {
			return err
		}

		hackathon.Close()

		return repo.Add(*hackathon)
	})

	if err != nil {
		return err
	}

	var participantIDs []string
	for _, p := range participants {
		participantIDs = append(participantIDs, p.ID.String())
	}

	return h.scoring.RemoveTasks(participantIDs)
}
