package query

import (
	"database/sql"
	log "github.com/sirupsen/logrus"
	"go-hackaton/src/pkg/common/application/errors"
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
		"s.created_at " +
		"FROM `session` s " +
		"LEFT JOIN session_participant sp ON (s.session_id = sp.session_id) " +
		"GROUP BY s.session_id")

	if err != nil {
		log.Error(err)
		return nil, errors.InternalError
	}
	defer rows.Close()

	sessions := make([]data.SessionData, 0)
	for rows.Next() {
		session, err := parseSession(rows)
		if err != nil {
			log.Error(err)
			return nil, errors.InternalError
		}

		sessions = append(sessions, *session)
	}

	return sessions, nil
}

func parseSession(r *sql.Rows) (*data.SessionData, error) {
	var sessionId string
	var name string
	var participants int
	var t int
	var createdAt time.Time

	err := r.Scan(&sessionId, &name, &participants, &t, &createdAt)
	if err != nil {
		return nil, err
	}

	return &data.SessionData{
		ID:           sessionId,
		Name:         name,
		Participants: participants,
		Type:         t,
		CreatedAt:    createdAt,
	}, nil
}
