package repository

import (
	"context"
	"database/sql"
	"github.com/google/uuid"
	"go-hackaton/src/pkg/common/infrastructure/repository"
	"go-hackaton/src/pkg/sessions/model"
	"time"
)

type sessionRepository struct {
	d repository.Database
}

func NewSessionRepository(db *sql.DB) model.SessionRepository {
	return &sessionRepository{d: repository.Database{DB: db}}
}

func (s *sessionRepository) Add(session model.Session) error {
	return s.d.Tx(func(tx *sql.Tx, ctx context.Context, closeTx func(error) error) error {
		_, err := tx.ExecContext(
			ctx,
			"INSERT INTO `session` (`session_id`, `code`, `code_hash`, `name`, `type`, `created_at`, `updated_at`) VALUES (UUID_TO_BIN(?), ?, UNHEX(MD5(?)), ?, ?, ?, ?)",
			session.ID, session.Code, session.Code, session.Name, session.Type, session.CreatedAt, session.CreatedAt)

		return closeTx(err)
	})
}

func (s *sessionRepository) Get(id uuid.UUID) (*model.Session, error) {
	rows, err := s.d.Query(""+
		"SELECT "+
		"BIN_TO_UUID(s.session_id) AS session_id, "+
		"s.code, "+
		"s.name, "+
		"s.type, "+
		"s.created_at "+
		"FROM `session` s "+
		"WHERE s.closed_at IS NULL AND BIN_TO_UUID(s.session_id) = ? ", id)

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
	rows, err := s.d.Query(""+
		"SELECT "+
		"BIN_TO_UUID(s.session_id) AS session_id, "+
		"s.code, "+
		"s.name, "+
		"s.type, "+
		"s.created_at "+
		"FROM `session` s "+
		"WHERE s.closed_at IS NULL AND s.code = ? ", code)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	if rows.Next() {
		return parseSession(rows)
	}

	return nil, nil // not found
}

func parseSession(r *sql.Rows) (*model.Session, error) {
	var sessionId string
	var code string
	var name string
	var t int
	var createdAt time.Time

	err := r.Scan(&sessionId, &code, &name, &t, &createdAt)
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
	}, nil
}
