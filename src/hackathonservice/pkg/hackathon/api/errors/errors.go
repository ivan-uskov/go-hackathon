package errors

import (
	"errors"
	log "github.com/sirupsen/logrus"
	commonErrors "go-hackathon/src/common/application/errors"
	appErrors "go-hackathon/src/hackathonservice/pkg/hackathon/application/errors"
)

var InvalidHackathonIDError = errors.New("hackathon: invalid hackathon id")
var ParticipantAlreadyExistsError = errors.New("hackathon: participant already exists")
var ParticipantNameIsEmptyError = errors.New("hackathon: participant name is empty")
var ParticipantEndpointIsEmptyError = errors.New("hackathon: participant endpoint is empty")
var InvalidParticipantScoreError = errors.New("hackathon: invalid participant score")
var HackathonNotExistsError = errors.New("hackathon: hackathon not exists")
var HackathonAlreadyExistsError = errors.New("hackathon: hackathon already exists")
var HackathonClosedError = errors.New("hackathon: hackathon closed")
var HackathonAlreadyClosedError = errors.New("hackathon: hackathon already closed")
var InvalidHackathonNameError = errors.New("hackathon: invalid hackathon name")
var InvalidHackathonTypeError = errors.New("hackathon: invalid hackathon type")
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
	case appErrors.InvalidParticipantScoreError:
		return InvalidParticipantScoreError
	case appErrors.HackathonNotExistsError:
		return HackathonNotExistsError
	case appErrors.HackathonAlreadyExistsError:
		return HackathonAlreadyExistsError
	case appErrors.HackathonClosedError:
		return HackathonClosedError
	case appErrors.HackathonAlreadyClosedError:
		return HackathonAlreadyClosedError
	case appErrors.InvalidHackathonNameError:
		return InvalidHackathonNameError
	case appErrors.InvalidHackathonTypeError:
		return InvalidHackathonTypeError
	case commonErrors.InternalError:
		return InternalError
	default:
		log.Error(err)
		return InternalError
	}
}
