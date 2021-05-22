package command

import (
	"github.com/google/uuid"
	"go-hackaton/src/pkg/sessions/model"
	"time"
)

var mockSession = model.Session{
	ID:        uuid.New(),
	Code:      "MockSession",
	Name:      "MockSession",
	Type:      model.SessionTypeArithmeticExpression,
	CreatedAt: time.Now(),
}

var mockParticipant = model.Participant{
	ID:        uuid.New(),
	SessionID: uuid.New(),
	Name:      "MockSession",
	Endpoint:  "http://localhost",
	Score:     0,
	CreatedAt: time.Now(),
}

type mockUnitOfWork struct {
	mockSessionRepository
	mockParticipantRepository
}
type mockSessionRepository struct {
	sessions map[string]model.Session
}
type mockParticipantRepository struct {
	participants map[string]model.Participant
}

func (m *mockUnitOfWork) Execute(job func(rp RepositoryProvider) error) error {
	return job(m)
}

func (m *mockUnitOfWork) SessionRepository() model.SessionRepository {
	return &m.mockSessionRepository
}

func (m *mockUnitOfWork) ParticipantRepository() model.ParticipantRepository {
	return &m.mockParticipantRepository
}

func (m *mockSessionRepository) Add(s model.Session) error {
	if m.sessions == nil {
		m.sessions = make(map[string]model.Session)
	}

	m.sessions[s.ID.String()] = s
	return nil
}

func (m *mockSessionRepository) Get(id uuid.UUID) (*model.Session, error) {
	if m.sessions == nil {
		return nil, nil
	}

	s, found := m.sessions[id.String()]
	if !found {
		return nil, nil
	}

	return &s, nil
}

func (m *mockSessionRepository) GetBySessionCode(code string) (*model.Session, error) {
	if m.sessions == nil {
		return nil, nil
	}

	for _, sess := range m.sessions {
		if sess.Code == code {
			return &sess, nil
		}
	}

	return nil, nil
}

func (m *mockParticipantRepository) Add(p model.Participant) error {
	if m.participants == nil {
		m.participants = make(map[string]model.Participant)
	}

	m.participants[p.ID.String()] = p
	return nil
}

func (m *mockParticipantRepository) Get(id uuid.UUID) (*model.Participant, error) {
	if m.participants == nil {
		return nil, nil
	}

	p, found := m.participants[id.String()]
	if !found {
		return nil, nil
	}

	return &p, nil
}

func (m *mockParticipantRepository) GetByName(name string) (*model.Participant, error) {
	if m.participants == nil {
		return nil, nil
	}

	for _, p := range m.participants {
		if p.Name == name {
			return &p, nil
		}
	}

	return nil, nil
}
