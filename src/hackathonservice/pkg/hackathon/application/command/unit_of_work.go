package command

import (
	"go-hackathon/src/hackathonservice/pkg/hackathon/model"
)

type RepositoryProvider interface {
	HackathonRepository() model.HackathonRepository
	ParticipantRepository() model.ParticipantRepository
}

type UnitOfWork interface {
	Execute(func(rp RepositoryProvider) error) error
}
