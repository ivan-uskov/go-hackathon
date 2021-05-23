package input

import (
	"github.com/google/uuid"
	"go-hackaton/src/pkg/sessions/api/errors"
	"go-hackaton/src/pkg/sessions/application/command"
)

type AddSessionInput struct {
	Code string
	Name string
	Type string
}

func (i AddSessionInput) Command() command.AddSessionCommand {
	return command.AddSessionCommand{
		Code: i.Code,
		Name: i.Name,
		Type: i.Type,
	}
}

type CloseSessionInput struct {
	SessionID string
}

func (i CloseSessionInput) Command() (*command.CloseSessionCommand, error) {
	sessionID, err := uuid.Parse(i.SessionID)
	if err != nil {
		return nil, errors.InvalidSessionIdError
	}

	return &command.CloseSessionCommand{SessionID: sessionID}, nil
}
