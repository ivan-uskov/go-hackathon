package input

import (
	"errors"
	"github.com/google/uuid"
	"go-hackaton/src/pkg/sessions/application/command"
	"go-hackaton/src/pkg/sessions/model"
)

const SessionTypeExpressions = "expressions"

type AddSessionInput struct {
	Code string
	Name string
	Type string
}

func (i AddSessionInput) Command() (*command.AddSessionCommand, error) {
	if i.Type != SessionTypeExpressions {
		return nil, errors.New("invalid session type")
	}

	return &command.AddSessionCommand{
		Code: i.Code,
		Name: i.Name,
		Type: model.SessionTypeArithmeticExpression,
	}, nil
}

type CloseSessionInput struct {
	SessionID string
}

func (i CloseSessionInput) Command() (*command.CloseSessionCommand, error) {
	sessionID, err := uuid.Parse(i.SessionID)
	if err != nil {
		return nil, err
	}

	return &command.CloseSessionCommand{SessionID: sessionID}, nil
}
