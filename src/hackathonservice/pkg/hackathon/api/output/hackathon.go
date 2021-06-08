package output

import (
	"go-hackathon/src/hackathonservice/pkg/hackathon/application/query/data"
	"time"
)

type HackathonOutput struct {
	ID           string
	Name         string
	Participants int
	Type         string
	CreatedAt    time.Time
	ClosedAt     *time.Time
}

func NewHackathonOutput(data data.HackathonData) HackathonOutput {
	return HackathonOutput{
		data.ID,
		data.Name,
		data.Participants,
		data.Type,
		data.CreatedAt,
		data.ClosedAt,
	}
}
