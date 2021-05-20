package api

import (
	"database/sql"
	"github.com/google/uuid"
	"go-hackaton/src/pkg/sessions/api/input"
	"go-hackaton/src/pkg/sessions/api/output"
	"go-hackaton/src/pkg/sessions/application/command"
	"go-hackaton/src/pkg/sessions/application/query"
	queryImpl "go-hackaton/src/pkg/sessions/infrastructure/query"
	"go-hackaton/src/pkg/sessions/infrastructure/repository"
	"go-hackaton/src/pkg/sessions/model"
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
	sqs      query.SessionQueryService
	pqs      query.ParticipantQueryService
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

func (a *api) GetSession(id string) (*output.SessionOutput, error) {
	session, err := a.sqs.GetSession(id)
	if err != nil {
		return nil, err
	}

	out := output.NewSessionOutput(*session)

	return &out, nil
}

func (a *api) GetSessionParticipants(sessionId string) ([]output.ParticipantOutput, error) {
	participants, err := a.pqs.GetParticipants(sessionId)
	if err != nil {
		return nil, err
	}

	participantsOutput := make([]output.ParticipantOutput, len(participants))
	for i, participant := range participants {
		participantsOutput[i] = output.NewParticipantOutput(participant)
	}

	return participantsOutput, nil
}

func (a *api) GetFirstScoredParticipantBefore(time time.Time) (*output.ParticipantOutput, error) {
	participant, err := a.pqs.GetFirstScoredParticipantBefore(time)
	if err != nil {
		return nil, err
	}

	var participantOutput *output.ParticipantOutput
	if participant != nil {
		out := output.NewParticipantOutput(*participant)
		participantOutput = &out
	}

	return participantOutput, nil
}

func (a *api) AddSession(in input.AddSessionInput) (*uuid.UUID, error) {
	c, err := in.Command()
	if err != nil {
		return nil, err
	}

	h := command.NewAddSessionCommandHandler(a.sessRepo)
	return h.Handle(*c)
}

func (a *api) CloseSession(in input.CloseSessionInput) error {
	c, err := in.Command()
	if err != nil {
		return err
	}

	h := command.NewCloseSessionCommandHandler(a.sessRepo)
	return h.Handle(*c)
}

func (a *api) AddSessionParticipant(in input.AddSessionParticipantInput) error {
	c, err := in.Command()
	if err != nil {
		return err
	}

	h := command.NewAddParticipantCommandHandler(a.sessRepo, a.partRepo)
	return h.Handle(*c)
}

func (a *api) UpdateSessionParticipantScore(in input.UpdateSessionParticipantScoreInput) error {
	c, err := in.Command()
	if err != nil {
		return err
	}

	h := command.NewUpdateParticipantScoreCommandHandler(a.partRepo)
	return h.Handle(*c)
}

func NewApi(db *sql.DB) Api {
	return &api{
		sqs:      queryImpl.NewSessionQueryService(db),
		pqs:      queryImpl.NewParticipantQueryService(db),
		partRepo: repository.NewParticipantRepository(db),
		sessRepo: repository.NewSessionRepository(db),
	}
}
