package command

import (
	"github.com/google/uuid"
	"go-hackaton/src/pkg/sessions/application/adapter"
	"go-hackaton/src/pkg/sessions/application/errors"
	"go-hackaton/src/pkg/sessions/model"
	"time"
)

type AddSessionCommand struct {
	Code string
	Name string
	Type string
}

type addSessionCommandHandler struct {
	unitOfWork UnitOfWork
	tasks      adapter.TaskAdapter
}

type AddSessionCommandHandler interface {
	Handle(command AddSessionCommand) (*uuid.UUID, error)
}

func NewAddSessionCommandHandler(unitOfWork UnitOfWork, tasks adapter.TaskAdapter) AddSessionCommandHandler {
	return &addSessionCommandHandler{unitOfWork, tasks}
}

func (h *addSessionCommandHandler) Handle(c AddSessionCommand) (*uuid.UUID, error) {
	var sessionId *uuid.UUID
	err := h.unitOfWork.Execute(func(rp RepositoryProvider) error {
		sessRepo := rp.SessionRepository()

		if c.Code == "" {
			return errors.InvalidSessionCodeError
		}
		if c.Name == "" {
			return errors.InvalidSessionNameError
		}

		taskType, valid := h.tasks.TranslateType(c.Type)
		if !valid {
			return errors.InvalidTaskTypeError
		}

		s, err := sessRepo.GetBySessionCode(c.Code)
		if err != nil {
			return err
		}
		if s != nil {
			return errors.SessionAlreadyExistsError
		}

		id := uuid.New()
		sessionId = &id
		return sessRepo.Add(model.Session{
			ID:        id,
			Code:      c.Code,
			Name:      c.Name,
			Type:      taskType,
			CreatedAt: time.Now(),
		})
	})

	return sessionId, err
}
