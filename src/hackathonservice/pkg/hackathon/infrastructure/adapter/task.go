package adapter

import (
	"go-hackathon/src/hackathonservice/pkg/hackathon/application/adapter"
)

type taskAdapter struct{}

func (a *taskAdapter) TranslateType(t string) (int, bool) {
	return 1, true //TODO : add scoring service call
}

func NewTaskAdapter() adapter.TaskAdapter {
	return &taskAdapter{}
}
