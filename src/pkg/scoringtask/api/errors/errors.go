package errors

import (
	"errors"
	log "github.com/sirupsen/logrus"
	commonErrors "go-hackaton/src/pkg/common/application/errors"
	appErrors "go-hackaton/src/pkg/scoringtask/application/errors"
)

var InvalidSolutionIdError = errors.New("scoring task: invalid solution id")
var TaskAlreadyExistError = errors.New("scoring task: task already exist")
var InvalidTaskTypeError = errors.New("scoring task: invalid task type")
var TaskNotExistError = errors.New("scoring task: task not exist")
var InternalError = errors.New("scoring task: internal error")

func WrapError(err error) error {
	switch err {
	case nil:
		return nil
	case appErrors.TaskAlreadyExistError:
		return TaskAlreadyExistError
	case appErrors.InvalidTaskTypeError:
		return InvalidTaskTypeError
	case appErrors.TaskNotExistError:
		return TaskNotExistError
	case commonErrors.InternalError:
		return InternalError
	default:
		log.Error(err)
		return InternalError
	}
}
