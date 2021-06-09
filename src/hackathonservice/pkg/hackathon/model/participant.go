package model

import (
	"github.com/google/uuid"
	"time"
)

type Participant struct {
	ID          uuid.UUID
	HackathonID uuid.UUID
	Name        string
	Endpoint    string
	Score       int
	CreatedAt   time.Time
	ScoredAt    *time.Time
}

type ParticipantRepository interface {
	Add(p Participant) error
	Get(id uuid.UUID) (*Participant, error)
	GetByName(name string) (*Participant, error)
	GetByHackathonID(hackathonID uuid.UUID) ([]Participant, error)
}

func (p *Participant) UpdateScore(s int) {
	now := time.Now()
	p.Score = s
	p.ScoredAt = &now
}
