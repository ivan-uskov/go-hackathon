package command

import (
	"github.com/google/uuid"
	"go-hackaton/src/pkg/sessions/application/errors"
	"go-hackaton/src/pkg/sessions/model"
	"time"
)

type AddParticipantCommand struct {
	SessionId   uuid.UUID
	SessionCode string
	Name        string
	Endpoint    string
}

type addParticipantCommandHandler struct {
	sessRepo model.SessionRepository
	partRepo model.ParticipantRepository
}

type AddParticipantCommandHandler interface {
	Handle(command AddParticipantCommand) error
}

func NewAddParticipantCommandHandler(sessRepo model.SessionRepository, partRepo model.ParticipantRepository) AddParticipantCommandHandler {
	return &addParticipantCommandHandler{sessRepo, partRepo}
}

func (apc *addParticipantCommandHandler) Handle(command AddParticipantCommand) error {
	session, err := apc.sessRepo.Get(command.SessionId)
	if err != nil {
		return err
	}

	if session == nil {
		return errors.SessionNotExistsError
	}

	if session.Code != command.SessionCode {
		return errors.InvalidSessionCodeError
	}

	if command.Name == "" {
		return errors.ParticipantNameIsEmptyError
	}

	participant, err := apc.partRepo.GetByName(command.Name)
	if err != nil {
		return err
	}
	if participant != nil {
		return errors.ParticipantAlreadyExistsError
	}

	return apc.partRepo.Add(model.Participant{
		ID:        uuid.New(),
		SessionID: command.SessionId,
		Endpoint:  command.Endpoint,
		Name:      command.Name,
		Score:     0,
		CreatedAt: time.Now(),
		ScoredAt:  nil,
	})
}
