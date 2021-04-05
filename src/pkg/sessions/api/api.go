package api

import (
	"database/sql"
	"go-hackaton/src/pkg/sessions/api/input"
	"go-hackaton/src/pkg/sessions/api/output"
	"go-hackaton/src/pkg/sessions/application/command"
	"go-hackaton/src/pkg/sessions/application/query"
	queryImpl "go-hackaton/src/pkg/sessions/infrastructure/query"
	"go-hackaton/src/pkg/sessions/infrastructure/repository"
	"go-hackaton/src/pkg/sessions/model"
)

type Api interface {
	GetSessions() ([]output.SessionOutput, error)
	AddSessionParticipant(in input.AddSessionParticipantInput) error
}

type api struct {
	sqs      query.SessionQueryService
	partRepo model.ParticipantRepository
	sessRepo model.SessionRepository
}

func (a *api) GetSessions() ([]output.SessionOutput, error) {
	sessions, err := a.sqs.GetSessions()
	if err != nil {
		return nil, err
	}

	sessionsOutput := make([]output.SessionOutput, len(sessions))
	for i, session := range sessions {
		sessionsOutput[i] = output.NewSessionOutput(session)
	}

	return sessionsOutput, nil
}

func (a *api) AddSessionParticipant(in input.AddSessionParticipantInput) error {
	c, err := in.Command()
	if err != nil {
		return err
	}

	h := command.NewAddParticipantCommandHandler(a.sessRepo, a.partRepo)
	return h.Handle(*c)
}

func NewApi(db *sql.DB) Api {
	return &api{
		sqs:      queryImpl.NewSessionQueryService(db),
		partRepo: repository.NewParticipantRepository(db),
		sessRepo: repository.NewSessionRepository(db),
	}
}
