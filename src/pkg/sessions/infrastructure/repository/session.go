package repository

import (
	"database/sql"
	"github.com/google/uuid"
	"go-hackaton/src/pkg/common/infrastructure/repository"
	"go-hackaton/src/pkg/sessions/model"
	"time"
)

type sessionRepository struct {
	tx *sql.Tx
}

func (s *sessionRepository) Add(session model.Session) error {
	_, err := s.tx.Exec(
		"INSERT INTO `session` (`session_id`, `code`, `code_hash`, `name`, `type`, `created_at`, `updated_at`) VALUES (UUID_TO_BIN(?), ?, UNHEX(MD5(?)), ?, ?, ?, ?)"+
			"ON DUPLICATE KEY UPDATE `code` = ?, `code_hash` = UNHEX(MD5(?)), `name` = ?, `type` = ?, `updated_at` = NOW(), `closed_at` = ?",
		session.ID, session.Code, session.Code, session.Name, session.Type, session.CreatedAt, session.CreatedAt,
		session.Code, session.Code, session.Name, session.Type, session.ClosedAt)

	return err
}

func (s *sessionRepository) Get(id uuid.UUID) (*model.Session, error) {
	rows, err := s.tx.Query(""+
		getSelectSessionQuery()+
		"WHERE BIN_TO_UUID(s.session_id) = ? ", id)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	if rows.Next() {
		return parseSession(rows)
	}

	return nil, nil // not found
}

func (s *sessionRepository) GetBySessionCode(code string) (*model.Session, error) {
	rows, err := s.tx.Query(""+
		getSelectSessionQuery()+
		"WHERE s.code = ? ", code)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	if rows.Next() {
		return parseSession(rows)
	}

	return nil, nil // not found
}

func getSelectSessionQuery() string {
	return "" +
		"SELECT " +
		"BIN_TO_UUID(s.session_id) AS session_id, " +
		"s.code, " +
		"s.name, " +
		"s.type, " +
		"s.created_at, " +
		"s.closed_at " +
		"FROM `session` s "
}

func parseSession(r *sql.Rows) (*model.Session, error) {
	var sessionId string
	var code string
	var name string
	var t int
	var createdAt time.Time
	var closedAtNullable sql.NullTime

	err := r.Scan(&sessionId, &code, &name, &t, &createdAt, &closedAtNullable)
	if err != nil {
		return nil, err
	}

	sessionUid, err := uuid.Parse(sessionId)
	if err != nil {
		return nil, err
	}

	return &model.Session{
		ID:        sessionUid,
		Code:      code,
		Name:      name,
		Type:      t,
		CreatedAt: createdAt,
		ClosedAt:  repository.TimePointer(closedAtNullable),
	}, nil
}
