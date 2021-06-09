package errors

import "errors"

var HackathonNotExistsError = errors.New("hackathon not exists")
var HackathonAlreadyExistsError = errors.New("hackathon already exists")
var HackathonClosedError = errors.New("hackathon closed")
var HackathonAlreadyClosedError = errors.New("hackathon already closed")
var InvalidHackathonNameError = errors.New("invalid hackathon name")
var InvalidHackathonTypeError = errors.New("invalid hackathon type")
