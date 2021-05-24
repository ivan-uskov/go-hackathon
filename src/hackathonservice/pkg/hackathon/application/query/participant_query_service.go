package query

import (
	"go-hackathon/src/hackathonservice/pkg/hackathon/application/query/data"
)

type ParticipantQueryService interface {
	GetParticipants(hackathonID string) ([]data.ParticipantData, error)
}
