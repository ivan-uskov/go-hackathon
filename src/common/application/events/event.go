package events

type Event interface {
	GetType() string
}

type EventStore interface {
	Add(e Event) error
}
