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
	unitOfWork UnitOfWork
}

type AddParticipantCommandHandler interface {
	Handle(command AddParticipantCommand) error
}

func NewAddParticipantCommandHandler(unitOfWork UnitOfWork) AddParticipantCommandHandler {
	return &addParticipantCommandHandler{unitOfWork}
}

func (h *addParticipantCommandHandler) Handle(command AddParticipantCommand) error {
	return h.unitOfWork.Execute(func(rp RepositoryProvider) error {
		sessRepo := rp.SessionRepository()
		partRepo := rp.ParticipantRepository()

		session, err := sessRepo.Get(command.SessionId)
		if err != nil {
			return err
		}

		if session == nil {
			return errors.SessionNotExistsError
		}

		if session.Code != command.SessionCode {
			return errors.InvalidSessionCodeError
		}

		if session.IsClose() {
			return errors.SessionClosedError
		}

		if command.Name == "" {
			return errors.ParticipantNameIsEmptyError
		}
		if command.Endpoint == "" {
			return errors.ParticipantEndpointIsEmptyError
		}

		participant, err := partRepo.GetByName(command.Name)
		if err != nil {
			return err
		}
		if participant != nil {
			return errors.ParticipantAlreadyExistsError
		}

		return partRepo.Add(model.Participant{
			ID:        uuid.New(),
			SessionID: command.SessionId,
			Endpoint:  removeSlashFromEnd(command.Endpoint),
			Name:      command.Name,
			Score:     0,
			CreatedAt: time.Now(),
			ScoredAt:  nil,
		})
	})
}

func removeSlashFromEnd(endpoint string) string {
	if endpoint[len(endpoint)-1:] == "/" {
		return endpoint[:len(endpoint)-1]
	}

	return endpoint
}
