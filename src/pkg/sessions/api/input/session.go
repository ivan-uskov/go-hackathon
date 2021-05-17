package input

import (
	"errors"
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
