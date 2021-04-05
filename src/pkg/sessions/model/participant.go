package model

import "github.com/google/uuid"

type Participant struct {
	ID        uuid.UUID
	Endpoint  string
	Score     int
	CreatedAt uuid.Time
}

type ParticipantRepository interface {
	Add(order Participant) error
	Get(id uuid.UUID) (*Participant, error)
}
