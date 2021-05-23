package query

import (
	"database/sql"
	"go-hackaton/src/pkg/common/infrastructure"
	"go-hackaton/src/pkg/common/infrastructure/repository"
	"go-hackaton/src/pkg/sessions/application/query"
	"go-hackaton/src/pkg/sessions/application/query/data"
	"time"
)

func NewSessionQueryService(db *sql.DB) query.SessionQueryService {
	return &sessionQueryService{db: db}
}

type sessionQueryService struct {
	db *sql.DB
}

func (qs *sessionQueryService) GetSessions() ([]data.SessionData, error) {
	rows, err := qs.db.Query("" +
		"SELECT " +
		"BIN_TO_UUID(s.session_id) AS session_id, " +
		"s.name, " +
		"COUNT(DISTINCT sp.participant_id), " +
		"s.type, " +
		"s.created_at, " +
		"s.closed_at " +
		"FROM `session` s " +
		"LEFT JOIN session_participant sp ON (s.session_id = sp.session_id) " +
		"GROUP BY s.session_id")

	if err != nil {
		return nil, infrastructure.InternalError(err)
	}
	defer infrastructure.CloseRows(rows)

	sessions := make([]data.SessionData, 0)
	for rows.Next() {
		session, err := parseSession(rows)
		if err != nil {
			return nil, infrastructure.InternalError(err)
		}

		sessions = append(sessions, *session)
	}

	return sessions, nil
}

func (qs *sessionQueryService) GetSession(id string) (*data.SessionData, error) {
	rows, err := qs.db.Query(""+
		"SELECT "+
		"BIN_TO_UUID(s.session_id) AS session_id, "+
		"s.name, "+
		"COUNT(DISTINCT sp.participant_id), "+
		"s.type, "+
		"s.created_at, "+
		"s.closed_at "+
		"FROM `session` s "+
		"LEFT JOIN session_participant sp ON (s.session_id = sp.session_id) "+
		"WHERE s.session_id = UUID_TO_BIN(?)"+
		"GROUP BY s.session_id", id)

	if err != nil {
		return nil, infrastructure.InternalError(err)
	}
	defer infrastructure.CloseRows(rows)

	if rows.Next() {
		session, err := parseSession(rows)
		if err != nil {
			return nil, infrastructure.InternalError(err)
		}

		return session, nil
	}

	return nil, nil // not found
}

func parseSession(r *sql.Rows) (*data.SessionData, error) {
	var sessionId string
	var name string
	var participants int
	var t int
	var createdAt time.Time
	var closedAtNullable sql.NullTime

	err := r.Scan(&sessionId, &name, &participants, &t, &createdAt, &closedAtNullable)
	if err != nil {
		return nil, infrastructure.InternalError(err)
	}

	return &data.SessionData{
		ID:           sessionId,
		Name:         name,
		Participants: participants,
		Type:         t,
		CreatedAt:    createdAt,
		ClosedAt:     repository.TimePointer(closedAtNullable),
	}, nil
}
