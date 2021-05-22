package errors

import "errors"

var SessionNotExistsError = errors.New("session not exists")
var SessionAlreadyExistsError = errors.New("session already exists")
var SessionClosedError = errors.New("session closed")
var InvalidSessionCodeError = errors.New("invalid session code")
var InvalidSessionNameError = errors.New("invalid session name")
