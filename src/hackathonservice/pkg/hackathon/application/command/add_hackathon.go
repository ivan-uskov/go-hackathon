package command

import (
	"github.com/google/uuid"
	"go-hackathon/src/hackathonservice/pkg/hackathon/application/adapter"
	"go-hackathon/src/hackathonservice/pkg/hackathon/application/errors"
	"go-hackathon/src/hackathonservice/pkg/hackathon/model"
	"time"
)

type AddHackathonCommand struct {
	Code string
	Name string
	Type string
}

type addHackathonCommandHandler struct {
	unitOfWork UnitOfWork
	tasks      adapter.TaskAdapter
}

type AddHackathonCommandHandler interface {
	Handle(command AddHackathonCommand) (*uuid.UUID, error)
}

func NewAddHackathonCommandHandler(unitOfWork UnitOfWork, tasks adapter.TaskAdapter) AddHackathonCommandHandler {
	return &addHackathonCommandHandler{unitOfWork, tasks}
}

func (h *addHackathonCommandHandler) Handle(c AddHackathonCommand) (*uuid.UUID, error) {
	var hackathonId *uuid.UUID
	err := h.unitOfWork.Execute(func(rp RepositoryProvider) error {
		repo := rp.HackathonRepository()

		if c.Code == "" {
			return errors.InvalidHackathonCodeError
		}
		if c.Name == "" {
			return errors.InvalidHackathonNameError
		}

		taskType, valid := h.tasks.TranslateType(c.Type)
		if !valid {
			return errors.InvalidHackathonTypeError
		}

		s, err := repo.GetByCode(c.Code)
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
			Code:      c.Code,
			Name:      c.Name,
			Type:      taskType,
			CreatedAt: time.Now(),
		})
	})

	return hackathonId, err
}
