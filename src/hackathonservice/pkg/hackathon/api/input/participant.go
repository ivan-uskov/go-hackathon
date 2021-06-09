package input

import (
	"github.com/google/uuid"
	"go-hackathon/src/hackathonservice/pkg/hackathon/api/errors"
	"go-hackathon/src/hackathonservice/pkg/hackathon/application/command"
)

type AddHackathonParticipantInput struct {
	HackathonID string
	Name        string
	Endpoint    string
}

func (i AddHackathonParticipantInput) Command() (*command.AddParticipantCommand, error) {
	hackathonID, err := uuid.Parse(i.HackathonID)
	if err != nil {
		return nil, errors.InvalidHackathonIDError
	}

	return &command.AddParticipantCommand{
		HackathonID: hackathonID,
		Name:        i.Name,
		Endpoint:    i.Endpoint,
	}, nil
}
