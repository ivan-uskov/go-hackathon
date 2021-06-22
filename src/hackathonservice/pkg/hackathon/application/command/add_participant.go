package command

import (
	"github.com/google/uuid"
	"go-hackathon/src/hackathonservice/pkg/hackathon/application/adapter"
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
	scoring    adapter.ScoringAdapter
}

type AddParticipantCommandHandler interface {
	Handle(command AddParticipantCommand) error
}

func NewAddParticipantCommandHandler(unitOfWork UnitOfWork, scoring adapter.ScoringAdapter) AddParticipantCommandHandler {
	return &addParticipantCommandHandler{unitOfWork, scoring}
}

func (h *addParticipantCommandHandler) Handle(command AddParticipantCommand) error {
	err := h.validateCommand(command)
	if err != nil {
		return err
	}

	hackathon, participant, err := h.addParticipant(command)
	if err != nil {
		return err
	}

	return h.scoring.AddTask(participant.ID.String(), hackathon.Type, participant.Endpoint)
}

func (h *addParticipantCommandHandler) validateCommand(command AddParticipantCommand) error {
	if command.Name == "" {
		return errors.ParticipantNameIsEmptyError
	}
	if command.Endpoint == "" {
		return errors.ParticipantEndpointIsEmptyError
	}

	return nil
}

func (h *addParticipantCommandHandler) addParticipant(command AddParticipantCommand) (*model.Hackathon, *model.Participant, error) {
	var hackathon *model.Hackathon
	var participant *model.Participant
	job := func(rp RepositoryProvider) error {
		hackRepo := rp.HackathonRepository()
		partRepo := rp.ParticipantRepository()

		var err error
		hackathon, err = hackRepo.Get(command.HackathonID)
		if err != nil {
			return err
		}

		if hackathon == nil {
			return errors.HackathonNotExistsError
		}

		if hackathon.IsClosed() {
			return errors.HackathonClosedError
		}

		participant, err = partRepo.GetByNameAndHackathonID(command.Name, hackathon.ID)
		if err != nil {
			return err
		}
		if participant != nil {
			return errors.ParticipantAlreadyExistsError
		}
		participant = &model.Participant{
			ID:          uuid.New(),
			HackathonID: command.HackathonID,
			Endpoint:    removeSlashFromEnd(command.Endpoint),
			Name:        command.Name,
			CreatedAt:   time.Now(),
		}

		return partRepo.Add(*participant)
	}

	job = h.unitOfWork.WithLock(getParticipantNameLock(command.Name), job)
	job = h.unitOfWork.WithLock(getHackathonIDLock(command.HackathonID), job)
	err := h.unitOfWork.Execute(job)

	return hackathon, participant, err
}

func removeSlashFromEnd(endpoint string) string {
	if endpoint[len(endpoint)-1:] == "/" {
		return endpoint[:len(endpoint)-1]
	}

	return endpoint
}
