package api

import (
	"database/sql"
	"github.com/google/uuid"
	"go-hackaton/src/pkg/sessions/api/input"
	"go-hackaton/src/pkg/sessions/api/output"
	"go-hackaton/src/pkg/sessions/application/adapter"
	"go-hackaton/src/pkg/sessions/application/command"
	"go-hackaton/src/pkg/sessions/application/query"
	adapterImpl "go-hackaton/src/pkg/sessions/infrastructure/adapter"
	queryImpl "go-hackaton/src/pkg/sessions/infrastructure/query"
	"go-hackaton/src/pkg/sessions/infrastructure/repository"
	tasks "go-hackaton/src/pkg/tasks/api"
	"time"
)

type Api interface {
	GetSessions() ([]output.SessionOutput, error)
	GetSession(id string) (*output.SessionOutput, error)
	GetSessionParticipants(sessionId string) ([]output.ParticipantOutput, error)
	GetFirstScoredParticipantBefore(time time.Time) (*output.ParticipantOutput, error)

	AddSession(in input.AddSessionInput) (*uuid.UUID, error)
	CloseSession(in input.CloseSessionInput) error
	AddSessionParticipant(in input.AddSessionParticipantInput) error
	UpdateSessionParticipantScore(in input.UpdateSessionParticipantScoreInput) error
}

type api struct {
	sqs        query.SessionQueryService
	pqs        query.ParticipantQueryService
	unitOfWork command.UnitOfWork

	tasks adapter.TaskAdapter
}

func NewApi(db *sql.DB, taskApi tasks.Api) Api {
	return &api{
		sqs:        queryImpl.NewSessionQueryService(db),
		pqs:        queryImpl.NewParticipantQueryService(db),
		unitOfWork: repository.NewUnitOfWork(db),
		tasks:      adapterImpl.NewTaskAdapter(taskApi),
	}
}
