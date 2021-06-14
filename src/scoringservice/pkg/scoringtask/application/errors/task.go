package errors

import "errors"

var TaskAlreadyExistError = errors.New("task already exist")
var InvalidTaskTypeError = errors.New("invalid task type")
var TaskNotExistError = errors.New("task not exist")

var ScorerNotExistError = errors.New("scorer not exist")
