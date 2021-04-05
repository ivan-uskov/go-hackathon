package query

import "go-hackaton/src/pkg/sessions/application/query/data"

type SessionQueryService interface {
	GetSessions() ([]data.SessionData, error)
}
