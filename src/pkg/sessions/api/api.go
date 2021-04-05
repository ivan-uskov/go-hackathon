package api

import (
	"database/sql"
	"go-hackaton/src/pkg/sessions/api/output"
	"go-hackaton/src/pkg/sessions/application/query"
	queryImpl "go-hackaton/src/pkg/sessions/infrastructure/query"
)

type Api interface {
	GetSessions() ([]output.SessionOutput, error)
}

type api struct {
	sqs query.SessionQueryService
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

func NewApi(db *sql.DB) Api {
	return &api{sqs: queryImpl.NewSessionQueryService(db)}
}
