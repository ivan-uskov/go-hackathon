package data

import "time"

type ParticipantData struct {
	ID        string
	Name      string
	Score     int
	CreatedAt time.Time
}
