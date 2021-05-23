package errors

import (
	"errors"
	log "github.com/sirupsen/logrus"
	commonErrors "go-hackaton/src/pkg/common/application/errors"
	appErrors "go-hackaton/src/pkg/sessions/application/errors"
)

var ParticipantAlreadyExistsError = errors.New("sessions: participant already exists")
var ParticipantNameIsEmptyError = errors.New("sessions: participant name is empty")
var ParticipantEndpointIsEmptyError = errors.New("sessions: participant endpoint is empty")
var SessionNotExistsError = errors.New("sessions: session not exists")
var SessionAlreadyExistsError = errors.New("sessions: session already exists")
var SessionClosedError = errors.New("sessions: session closed")
var SessionAlreadyClosedError = errors.New("sessions: session already closed")
var InvalidSessionCodeError = errors.New("sessions: invalid session code")
var InvalidSessionNameError = errors.New("sessions: invalid session name")
var InvalidTaskTypeError = errors.New("sessions: invalid task type")
var InternalError = commonErrors.InternalError

func WrapError(err error) error {
	switch err {
	case nil:
		return nil
	case appErrors.ParticipantAlreadyExistsError:
		return ParticipantAlreadyExistsError
	case appErrors.ParticipantNameIsEmptyError:
		return ParticipantNameIsEmptyError
	case appErrors.ParticipantEndpointIsEmptyError:
		return ParticipantEndpointIsEmptyError
	case appErrors.SessionNotExistsError:
		return SessionNotExistsError
	case appErrors.SessionAlreadyExistsError:
		return SessionAlreadyExistsError
	case appErrors.SessionClosedError:
		return SessionClosedError
	case appErrors.SessionAlreadyClosedError:
		return SessionAlreadyClosedError
	case appErrors.InvalidSessionCodeError:
		return InvalidSessionCodeError
	case appErrors.InvalidSessionNameError:
		return InvalidSessionNameError
	case appErrors.InvalidTaskTypeError:
		return InvalidTaskTypeError
	case commonErrors.InternalError:
		return InternalError
	default:
		log.Error(err)
		return InternalError
	}
}
