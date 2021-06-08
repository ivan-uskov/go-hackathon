package api

import (
	"github.com/google/uuid"
	"go-hackathon/src/hackathonservice/pkg/hackathon/api/errors"
	"go-hackathon/src/hackathonservice/pkg/hackathon/api/input"
	"go-hackathon/src/hackathonservice/pkg/hackathon/application/command"
)

func (a *api) AddHackathon(in input.AddHackathonInput) (*uuid.UUID, error) {
	h := command.NewAddHackathonCommandHandler(a.unitOfWork, a.scoring)
	id, err := h.Handle(in.Command())
	return id, errors.WrapError(err)
}

func (a *api) CloseHackathon(in input.CloseHackathonInput) error {
	c, err := in.Command()
	if err != nil {
		return err
	}

	h := command.NewCloseHackathonCommandHandler(a.unitOfWork)
	return errors.WrapError(h.Handle(*c))
}

func (a *api) AddHackathonParticipant(in input.AddHackathonParticipantInput) error {
	c, err := in.Command()
	if err != nil {
		return err
	}

	h := command.NewAddParticipantCommandHandler(a.unitOfWork)
	return errors.WrapError(h.Handle(*c))
}
