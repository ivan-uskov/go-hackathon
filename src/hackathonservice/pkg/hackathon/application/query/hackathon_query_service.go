package query

import (
	"go-hackathon/src/hackathonservice/pkg/hackathon/application/query/data"
)

type HackathonQueryService interface {
	GetHackathons() ([]data.HackathonData, error)
	GetHackathon(id string) (*data.HackathonData, error)
}
