package input

import (
	"github.com/google/uuid"
	"go-hackaton/src/pkg/sessions/api/errors"
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
		return nil, errors.InvalidSessionIdError
	}

	return &command.AddParticipantCommand{
		SessionId:   sessionID,
		SessionCode: i.SessionCode,
		Name:        i.Name,
		Endpoint:    i.Endpoint,
	}, nil
}

type UpdateSessionParticipantScoreInput struct {
	ID    string
	Score int
}

func (i UpdateSessionParticipantScoreInput) Command() (*command.UpdateParticipantScoreCommand, error) {
	sessionID, err := uuid.Parse(i.ID)
	if err != nil {
		return nil, errors.InvalidSessionIdError
	}

	return &command.UpdateParticipantScoreCommand{
		ID:    sessionID,
		Score: i.Score,
	}, nil
}
