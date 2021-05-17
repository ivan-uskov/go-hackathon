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
	sessRepo model.SessionRepository
}

type AddSessionCommandHandler interface {
	Handle(command AddSessionCommand) (*uuid.UUID, error)
}

func NewAddSessionCommandHandler(sessRepo model.SessionRepository) AddSessionCommandHandler {
	return &addSessionCommandHandler{sessRepo}
}

func (apc *addSessionCommandHandler) Handle(c AddSessionCommand) (*uuid.UUID, error) {
	if c.Code == "" {
		return nil, errors.InvalidSessionCodeError
	}
	if c.Name == "" {
		return nil, errors.InvalidSessionNameError
	}

	s, err := apc.sessRepo.GetBySessionCode(c.Code)
	if err != nil {
		return nil, err
	}
	if s != nil {
		return nil, errors.SessionAlreadyExistsError
	}

	return nil, apc.sessRepo.Add(model.Session{
		ID:        uuid.New(),
		Code:      c.Code,
		Name:      c.Name,
		Type:      c.Type,
		CreatedAt: time.Now(),
	})
}