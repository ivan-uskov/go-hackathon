package command

import (
	"github.com/google/uuid"
	"go-hackathon/src/hackathonservice/pkg/hackathon/application/errors"
	"go-hackathon/src/hackathonservice/pkg/hackathon/model"
	"time"
)

type AddParticipantCommand struct {
	HackathonID uuid.UUID
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
		hackRepo := rp.HackathonRepository()
		partRepo := rp.ParticipantRepository()

		hackathon, err := hackRepo.Get(command.HackathonID)
		if err != nil {
			return err
		}

		if hackathon == nil {
			return errors.HackathonNotExistsError
		}

		if hackathon.IsClosed() {
			return errors.HackathonClosedError
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
			ID:          uuid.New(),
			HackathonID: command.HackathonID,
			Endpoint:    removeSlashFromEnd(command.Endpoint),
			Name:        command.Name,
			CreatedAt:   time.Now(),
		})
	})
}

func removeSlashFromEnd(endpoint string) string {
	if endpoint[len(endpoint)-1:] == "/" {
		return endpoint[:len(endpoint)-1]
	}

	return endpoint
}
