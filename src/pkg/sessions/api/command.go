package api

import (
	"github.com/google/uuid"
	"go-hackaton/src/pkg/sessions/api/input"
	"go-hackaton/src/pkg/sessions/application/command"
)

func (a *api) AddSession(in input.AddSessionInput) (*uuid.UUID, error) {
	c, err := in.Command()
	if err != nil {
		return nil, err
	}

	h := command.NewAddSessionCommandHandler(a.unitOfWork, a.tasks)
	return h.Handle(*c)
}

func (a *api) CloseSession(in input.CloseSessionInput) error {
	c, err := in.Command()
	if err != nil {
		return err
	}

	h := command.NewCloseSessionCommandHandler(a.unitOfWork)
	return h.Handle(*c)
}

func (a *api) AddSessionParticipant(in input.AddSessionParticipantInput) error {
	c, err := in.Command()
	if err != nil {
		return err
	}

	h := command.NewAddParticipantCommandHandler(a.unitOfWork)
	return h.Handle(*c)
}

func (a *api) UpdateSessionParticipantScore(in input.UpdateSessionParticipantScoreInput) error {
	c, err := in.Command()
	if err != nil {
		return err
	}

	h := command.NewUpdateParticipantScoreCommandHandler(a.unitOfWork)
	return h.Handle(*c)
}
