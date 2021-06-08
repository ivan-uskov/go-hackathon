package data

import (
	"time"
)

type HackathonData struct {
	ID           string
	Name         string
	Participants int
	Type         string
	CreatedAt    time.Time
	ClosedAt     *time.Time
}
