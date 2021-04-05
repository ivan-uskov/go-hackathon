package output

import (
	"go-hackaton/src/pkg/sessions/application/query/data"
	"time"
)

type ParticipantOutput struct {
	ID        string
	Name      string
	Score     int
	CreatedAt time.Time
}

func NewParticipantOutput(data data.ParticipantData) ParticipantOutput {
	return ParticipantOutput{
		data.ID,
		data.Name,
		data.Score,
		data.CreatedAt,
	}
}
