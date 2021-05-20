package data

import (
	"time"
)

type SessionData struct {
	ID           string
	Name         string
	Participants int
	Type         int
	CreatedAt    time.Time
	ClosedAt     *time.Time
}
