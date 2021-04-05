package query

import "go-hackaton/src/pkg/sessions/application/query/data"

type ParticipantQueryService interface {
	GetParticipants(sessionID string) ([]data.ParticipantData, error)
}
