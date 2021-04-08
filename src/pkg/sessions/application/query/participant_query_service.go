package query

import (
	"go-hackaton/src/pkg/sessions/application/query/data"
	"time"
)

type ParticipantQueryService interface {
	GetParticipants(sessionID string) ([]data.ParticipantData, error)
	GetFirstScoredParticipantBefore(time time.Time) (*data.ParticipantData, error)
}
