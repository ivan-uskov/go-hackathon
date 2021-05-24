package api

import (
	"database/sql"
	"github.com/google/uuid"
	"go-hackathon/src/hackathonservice/pkg/hackathon/api/input"
	"go-hackathon/src/hackathonservice/pkg/hackathon/api/output"
	"go-hackathon/src/hackathonservice/pkg/hackathon/application/adapter"
	"go-hackathon/src/hackathonservice/pkg/hackathon/application/command"
	"go-hackathon/src/hackathonservice/pkg/hackathon/application/query"
	adapterImpl "go-hackathon/src/hackathonservice/pkg/hackathon/infrastructure/adapter"
	queryImpl "go-hackathon/src/hackathonservice/pkg/hackathon/infrastructure/query"
	"go-hackathon/src/hackathonservice/pkg/hackathon/infrastructure/repository"
)

type Api interface {
	GetHackathons() ([]output.HackathonOutput, error)
	GetHackathon(id string) (*output.HackathonOutput, error)
	GetHackathonParticipants(hackathonID string) ([]output.ParticipantOutput, error)

	AddHackathon(in input.AddHackathonInput) (*uuid.UUID, error)
	CloseHackathon(in input.CloseHackathonInput) error
	AddHackathonParticipant(in input.AddHackathonParticipantInput) error
}

type api struct {
	sqs        query.HackathonQueryService
	pqs        query.ParticipantQueryService
	unitOfWork command.UnitOfWork

	tasks adapter.TaskAdapter
}

func NewApi(db *sql.DB) Api {
	return &api{
		sqs:        queryImpl.NewHackathonQueryService(db),
		pqs:        queryImpl.NewParticipantQueryService(db),
		unitOfWork: repository.NewUnitOfWork(db),
		tasks:      adapterImpl.NewTaskAdapter(),
	}
}
