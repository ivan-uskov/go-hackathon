package command

import (
	"github.com/google/uuid"
	"go-hackaton/src/pkg/sessions/application/errors"
	"go-hackaton/src/pkg/sessions/model"
	"time"
)

type AddSessionCommand struct {
	Code string
	Name string
	Type int
}

type addSessionCommandHandler struct {
	unitOfWork UnitOfWork
}

type AddSessionCommandHandler interface {
	Handle(command AddSessionCommand) (*uuid.UUID, error)
}

func NewAddSessionCommandHandler(unitOfWork UnitOfWork) AddSessionCommandHandler {
	return &addSessionCommandHandler{unitOfWork}
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
			Type:      c.Type,
			CreatedAt: time.Now(),
		})
	})

	return sessionId, err
}
