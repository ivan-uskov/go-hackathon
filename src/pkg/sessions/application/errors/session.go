package errors

import "errors"

var SessionNotExistsError = errors.New("session not exists")
var InvalidSessionCodeError = errors.New("invalid session code")
