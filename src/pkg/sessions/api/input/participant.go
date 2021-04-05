package input

import (
	"github.com/google/uuid"
	"go-hackaton/src/pkg/sessions/application/command"
)

type AddSessionParticipantInput struct {
	SessionID   string
	SessionCode string
	Name        string
	Endpoint    string
}

func (i AddSessionParticipantInput) Command() (*command.AddParticipantCommand, error) {
	sessionID, err := uuid.Parse(i.SessionID)
	if err != nil {
		return nil, err
	}

	return &command.AddParticipantCommand{
		SessionId:   sessionID,
		SessionCode: i.SessionCode,
		Name:        i.Name,
		Endpoint:    i.Endpoint,
	}, nil
}
