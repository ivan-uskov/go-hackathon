package events

import (
	"github.com/google/uuid"
	"time"
)

type StoredEvent struct {
	ID          uuid.UUID
	Type        string
	Body        string
	CreatedAt   time.Time
	PublishedAt *time.Time
}

type StoredEventRepository interface {
	Store(e StoredEvent) error
	GetNotPublishedEvents() ([]StoredEvent, error)
}

func (e *StoredEvent) Publish() {
	now := time.Now()
	e.PublishedAt = &now
}
