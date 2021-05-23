package errors

import "errors"

var ParticipantAlreadyExistsError = errors.New("participant already exists")
var ParticipantNameIsEmptyError = errors.New("participant name is empty")
var ParticipantEndpointIsEmptyError = errors.New("participant endpoint is empty")
var InvalidParticipantScoreError = errors.New("invalid participant score")
