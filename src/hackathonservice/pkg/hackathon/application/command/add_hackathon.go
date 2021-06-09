package command

import (
	"github.com/google/uuid"
	"go-hackathon/src/hackathonservice/pkg/hackathon/application/adapter"
	"go-hackathon/src/hackathonservice/pkg/hackathon/application/errors"
	"go-hackathon/src/hackathonservice/pkg/hackathon/model"
	"time"
)

type AddHackathonCommand struct {
	Name string
	Type string
}

type addHackathonCommandHandler struct {
	unitOfWork UnitOfWork
	scoring    adapter.ScoringAdapter
}

type AddHackathonCommandHandler interface {
	Handle(command AddHackathonCommand) (*uuid.UUID, error)
}

func NewAddHackathonCommandHandler(unitOfWork UnitOfWork, scoring adapter.ScoringAdapter) AddHackathonCommandHandler {
	return &addHackathonCommandHandler{unitOfWork, scoring}
}

func (h *addHackathonCommandHandler) Handle(c AddHackathonCommand) (*uuid.UUID, error) {
	var hackathonId *uuid.UUID
	err := h.unitOfWork.Execute(func(rp RepositoryProvider) error {
		repo := rp.HackathonRepository()

		if c.Name == "" {
			return errors.InvalidHackathonNameError
		}

		if !h.scoring.ValidateTaskType(c.Type) {
			return errors.InvalidHackathonTypeError
		}

		s, err := repo.GetByName(c.Name)
		if err != nil {
			return err
		}
		if s != nil {
			return errors.HackathonAlreadyExistsError
		}

		id := uuid.New()
		hackathonId = &id
		return repo.Add(model.Hackathon{
			ID:        id,
			Name:      c.Name,
			Type:      c.Type,
			CreatedAt: time.Now(),
		})
	})

	return hackathonId, err
}
