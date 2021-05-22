package adapter

import (
	"go-hackaton/src/pkg/sessions/application/adapter"
	"go-hackaton/src/pkg/tasks/api"
)

type taskAdapter struct {
	taskApi api.Api
}

func (a *taskAdapter) TranslateType(t string) (int, bool) {
	return a.taskApi.TranslateType(t)
}

func NewTaskAdapter(taskApi api.Api) adapter.TaskAdapter {
	return &taskAdapter{taskApi: taskApi}
}
