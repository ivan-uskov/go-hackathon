package api

import "go-hackaton/src/pkg/tasks/model"

type Api interface {
	TranslateType(t string) (int, bool)
}

type api struct{}

func (a *api) TranslateType(t string) (int, bool) {
	return model.TranslateType(t)
}

func NewApi() Api {
	return &api{}
}
