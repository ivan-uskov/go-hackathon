package model

import (
	"github.com/google/uuid"
	"time"
)

type Participant struct {
	ID        uuid.UUID
	SessionID uuid.UUID
	Name      string
	Endpoint  string
	Score     int
	CreatedAt time.Time
	ScoredAt  *time.Time
}

type ParticipantRepository interface {
	Add(p Participant) error
	Get(id uuid.UUID) (*Participant, error)
	GetByName(name string) (*Participant, error)
}
