package model

import (
	"github.com/google/uuid"
	"time"
)

type Hackathon struct {
	ID        uuid.UUID
	Name      string
	Type      string
	CreatedAt time.Time
	ClosedAt  *time.Time
}

type HackathonRepository interface {
	Add(s Hackathon) error
	Get(id uuid.UUID) (*Hackathon, error)
	GetByName(name string) (*Hackathon, error)
}

func (s *Hackathon) Close() {
	now := time.Now()
	s.ClosedAt = &now
}

func (s *Hackathon) IsClosed() bool {
	return s.ClosedAt != nil
}
