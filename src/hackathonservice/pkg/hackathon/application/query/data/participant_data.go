package data

import "time"

type ParticipantData struct {
	ID        string
	Name      string
	Score     int
	Endpoint  string
	CreatedAt time.Time
	ScoredAt  *time.Time
}
