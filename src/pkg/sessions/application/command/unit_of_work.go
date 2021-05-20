package command

import "go-hackaton/src/pkg/sessions/model"

type RepositoryProvider interface {
	SessionRepository() model.SessionRepository
	ParticipantRepository() model.ParticipantRepository
}

type UnitOfWork interface {
	Execute(func(rp RepositoryProvider) error) error
}
