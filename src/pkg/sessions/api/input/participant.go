package input

import (
	"errors"
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

type UpdateSessionParticipantScoreInput struct {
	ID    string
	Score int
}

func (i UpdateSessionParticipantScoreInput) Command() (*command.UpdateParticipantScoreCommand, error) {
	sessionID, err := uuid.Parse(i.ID)
	if err != nil {
		return nil, err
	}

	if i.Score < 0 {
		return nil, errors.New("score can't be less than 0")
	}

	return &command.UpdateParticipantScoreCommand{
		ID:    sessionID,
		Score: i.Score,
	}, nil
}
