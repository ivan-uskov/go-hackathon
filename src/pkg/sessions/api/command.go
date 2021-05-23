package api

import (
	"github.com/google/uuid"
	"go-hackaton/src/pkg/sessions/api/errors"
	"go-hackaton/src/pkg/sessions/api/input"
	"go-hackaton/src/pkg/sessions/application/command"
)

func (a *api) AddSession(in input.AddSessionInput) (*uuid.UUID, error) {
	c, err := in.Command()
	if err != nil {
		return nil, errors.WrapError(err)
	}

	h := command.NewAddSessionCommandHandler(a.unitOfWork, a.tasks)
	id, err := h.Handle(*c)
	return id, errors.WrapError(err)
}

func (a *api) CloseSession(in input.CloseSessionInput) error {
	c, err := in.Command()
	if err != nil {
		return errors.WrapError(err)
	}

	h := command.NewCloseSessionCommandHandler(a.unitOfWork)
	return errors.WrapError(h.Handle(*c))
}

func (a *api) AddSessionParticipant(in input.AddSessionParticipantInput) error {
	c, err := in.Command()
	if err != nil {
		return errors.WrapError(err)
	}

	h := command.NewAddParticipantCommandHandler(a.unitOfWork)
	return errors.WrapError(h.Handle(*c))
}

func (a *api) UpdateSessionParticipantScore(in input.UpdateSessionParticipantScoreInput) error {
	c, err := in.Command()
	if err != nil {
		return errors.WrapError(err)
	}

	h := command.NewUpdateParticipantScoreCommandHandler(a.unitOfWork)
	return errors.WrapError(h.Handle(*c))
}
