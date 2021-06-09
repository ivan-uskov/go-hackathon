package input

import (
	"github.com/google/uuid"
	"go-hackathon/src/hackathonservice/pkg/hackathon/api/errors"
	"go-hackathon/src/hackathonservice/pkg/hackathon/application/command"
)

type AddHackathonInput struct {
	Name string
	Type string
}

func (i AddHackathonInput) Command() command.AddHackathonCommand {
	return command.AddHackathonCommand{
		Name: i.Name,
		Type: i.Type,
	}
}

type CloseHackathonInput struct {
	HackathonID string
}

func (i CloseHackathonInput) Command() (*command.CloseHackathonCommand, error) {
	hackathonID, err := uuid.Parse(i.HackathonID)
	if err != nil {
		return nil, errors.InvalidHackathonIDError
	}

	return &command.CloseHackathonCommand{HackathonID: hackathonID}, nil
}
