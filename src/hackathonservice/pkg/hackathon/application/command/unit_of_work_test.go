package command

import (
	"github.com/google/uuid"
	"go-hackathon/src/hackathonservice/pkg/hackathon/model"
	"time"
)

var mockHackathon = model.Hackathon{
	ID:        uuid.New(),
	Name:      "MockHackathon",
	Type:      "test_type",
	CreatedAt: time.Now(),
}

var mockParticipant = model.Participant{
	ID:          uuid.New(),
	HackathonID: uuid.New(),
	Name:        "MockHackathonUser",
	Endpoint:    "http://localhost",
	Score:       0,
	CreatedAt:   time.Now(),
}

type mockUnitOfWork struct {
	mockHackathonRepository
	mockParticipantRepository
}
type mockHackathonRepository struct {
	hackathons map[string]model.Hackathon
}
type mockParticipantRepository struct {
	participants map[string]model.Participant
}

func (m *mockUnitOfWork) Execute(job func(rp RepositoryProvider) error) error {
	return job(m)
}

func (m *mockUnitOfWork) HackathonRepository() model.HackathonRepository {
	return &m.mockHackathonRepository
}

func (m *mockUnitOfWork) ParticipantRepository() model.ParticipantRepository {
	return &m.mockParticipantRepository
}

func (m *mockHackathonRepository) Add(s model.Hackathon) error {
	if m.hackathons == nil {
		m.hackathons = make(map[string]model.Hackathon)
	}

	m.hackathons[s.ID.String()] = s
	return nil
}

func (m *mockHackathonRepository) Get(id uuid.UUID) (*model.Hackathon, error) {
	if m.hackathons == nil {
		return nil, nil
	}

	s, found := m.hackathons[id.String()]
	if !found {
		return nil, nil
	}

	return &s, nil
}

func (m *mockHackathonRepository) GetByName(name string) (*model.Hackathon, error) {
	if m.hackathons == nil {
		return nil, nil
	}

	for _, h := range m.hackathons {
		if h.Name == name {
			return &h, nil
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
