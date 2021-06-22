package command

import (
	"go-hackathon/src/hackathonservice/pkg/hackathon/model"
)

type RepositoryProvider interface {
	HackathonRepository() model.HackathonRepository
	ParticipantRepository() model.ParticipantRepository
}

type Job func(rp RepositoryProvider) error

type UnitOfWork interface {
	Execute(Job) error
	WithLock(name string, job Job) Job
}
