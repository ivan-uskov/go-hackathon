package events

import (
	"database/sql"
	"encoding/json"
	"github.com/google/uuid"
	"go-hackathon/src/common/application/events"
	"go-hackathon/src/common/infrastructure"
	"go-hackathon/src/common/infrastructure/repository"
	eventsModel "go-hackathon/src/common/model/events"
	"time"
)

type storedEventUnitOfWork struct {
	db repository.Database
}

func NewStoredEventUnitOfWork(db *sql.DB) events.StoredEventUnitOfWork {
	return &storedEventUnitOfWork{db: repository.Database{DB: db}}
}

func (s *storedEventUnitOfWork) Execute(job func(eventsModel.StoredEventRepository) error) error {
	return s.db.Tx(func(tx *sql.Tx) error {
		return job(&eventStore{tx})
	})
}

func (s *storedEventUnitOfWork) Lock(name string, job func() error) error {
	return s.db.Lock(name, job)
}

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
		e.ID, e.Type, e.Body, e.CreatedAt,
		e.Type, e.Body, e.CreatedAt, e.PublishedAt,
	)

	if err != nil {
		return infrastructure.InternalError(err)
	}

	return err
}

func (s *eventStore) GetNotPublishedEvents() ([]eventsModel.StoredEvent, error) {
	rows, err := s.tx.Query("" +
		selectStoredEventSql() +
		"WHERE se.published_at IS NULL")

	if err != nil {
		return nil, infrastructure.InternalError(err)
	}
	defer infrastructure.Close(rows)

	var ee []eventsModel.StoredEvent
	for rows.Next() {
		e, err := parseStoredEvent(rows)
		if err != nil {
			return nil, err
		}

		ee = append(ee, *e)
	}

	return ee, nil
}

func parseStoredEvent(rows *sql.Rows) (*eventsModel.StoredEvent, error) {
	var id string
	var t string
	var body string
	var createdAt time.Time
	var publishedAtNullable sql.NullTime

	err := rows.Scan(&id, &t, &body, &createdAt, &publishedAtNullable)
	if err != nil {
		return nil, infrastructure.InternalError(err)
	}

	uid, err := uuid.Parse(id)
	if err != nil {
		return nil, infrastructure.InternalError(err)
	}

	return &eventsModel.StoredEvent{
		ID:          uid,
		Type:        t,
		Body:        body,
		CreatedAt:   createdAt,
		PublishedAt: repository.TimePointer(publishedAtNullable),
	}, nil
}

func selectStoredEventSql() string {
	return "" +
		"SELECT " +
		"BIN_TO_UUID(se.event_id) AS event_id, " +
		"se.type, " +
		"se.body, " +
		"se.created_at, " +
		"se.published_at " +
		"FROM `stored_event` se "
}
