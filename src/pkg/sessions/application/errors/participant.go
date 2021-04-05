package errors

import "errors"

var ParticipantAlreadyExistsError = errors.New("participant already exists")
var ParticipantNameIsEmptyError = errors.New("participant name is empty")
