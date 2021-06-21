package events

import (
	"database/sql"
	"encoding/json"
	"github.com/google/uuid"
	"go-hackathon/src/common/application/events"
	"go-hackathon/src/common/infrastructure"
	eventsModel "go-hackathon/src/common/model/events"
	"time"
)

type eventStore struct {
	tx *sql.Tx
}

func NewEventStore(tx *sql.Tx) events.EventStore {
	return &eventStore{tx}
}

func (s *eventStore) Add(e events.Event) error {
	b, err := json.Marshal(e)
	if err != nil {
		return err
	}

	return s.Store(eventsModel.StoredEvent{
		ID:        uuid.New(),
		Type:      e.GetType(),
		Body:      string(b),
		CreatedAt: time.Now(),
	})
}

func (s *eventStore) Store(e eventsModel.StoredEvent) error {
	_, err := s.tx.Exec(
		"INSERT INTO `stored_event` (`event_id`, `type`, `body`, `created_at`) VALUES (UUID_TO_BIN(?), ?, ?, ?)"+
			"ON DUPLICATE KEY UPDATE `type` = ?, `body` = ?, `created_at` = ?, `published_at` = ?",
		uuid.New().String(), e.Type, e.Body, e.CreatedAt,
		e.Type, e.Body, e.CreatedAt, e.PublishedAt,
	)

	if err != nil {
		return infrastructure.InternalError(err)
	}

	return err
}
