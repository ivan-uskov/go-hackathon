package model

import (
	"github.com/google/uuid"
	"time"
)

const HackathonTypeArithmeticExpression = 1

type Hackathon struct {
	ID        uuid.UUID
	Code      string
	Name      string
	Type      int
	CreatedAt time.Time
	ClosedAt  *time.Time
}

type HackathonRepository interface {
	Add(s Hackathon) error
	Get(id uuid.UUID) (*Hackathon, error)
	GetByCode(code string) (*Hackathon, error)
}

func (s *Hackathon) Close() {
	now := time.Now()
	s.ClosedAt = &now
}

func (s *Hackathon) IsClosed() bool {
	return s.ClosedAt != nil
}
