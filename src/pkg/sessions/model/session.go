package model

import (
	"github.com/google/uuid"
	"time"
)

const SessionTypeArithmeticExpression = 1

type Session struct {
	ID        uuid.UUID
	Code      string
	Name      string
	Type      int
	CreatedAt time.Time
}

type SessionRepository interface {
	Add(order Session) error
	Get(id uuid.UUID) (*Session, error)
}
