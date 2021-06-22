package events

import (
	"encoding/json"
	"errors"
	"github.com/google/uuid"
	"go-hackathon/src/hackathonservice/pkg/hackathon/application/command"
	"go-hackathon/src/hackathonservice/pkg/hackathon/application/event/integration"
)

const scoringCompleteEventType = "ScoringCompleteEvent"

var EventHandlerNotExistError = errors.New("event handler not exist")

type EventHandler func() error

type eventHandlerFactory struct {
	uow command.UnitOfWork
}

type EventHandlerFactory interface {
	Create(event string) (EventHandler, error)
}

func NewEventHandlerFactory(uow command.UnitOfWork) EventHandlerFactory {
	return &eventHandlerFactory{uow: uow}
}

func (f *eventHandlerFactory) Create(event string) (EventHandler, error) {
	t, payload, err := parseEvent(event)
	if err != nil {
		return nil, err
	}

	switch t {
	case scoringCompleteEventType:
		return f.createScoringCompleteEventHandler(payload)
	}

	return nil, EventHandlerNotExistError
}

type scoringCompleteEvent struct {
	SolutionID string
	Score      int
}

func (f *eventHandlerFactory) createScoringCompleteEventHandler(payload string) (EventHandler, error) {
	e := scoringCompleteEvent{}
	err := json.Unmarshal([]byte(payload), &e)
	if err != nil {
		return nil, err
	}

	uid, err := uuid.Parse(e.SolutionID)
	if err != nil {
		return nil, err
	}

	return func() error {
		h := integration.NewScoringCompleteEventHandler(f.uow)
		return h.Handle(uid, e.Score)
	}, nil
}

func parseEvent(event string) (string, string, error) {
	var eventFields map[string]string
	err := json.Unmarshal([]byte(event), &eventFields)
	if err != nil {
		return "", "", err
	}

	t, typeExists := eventFields["Type"]
	if !typeExists {
		return "", "", errors.New("undefined events type")
	}

	payload, payloadExists := eventFields["Payload"]
	if !payloadExists {
		return "", "", errors.New("undefined payload")
	}

	return t, payload, nil
}
