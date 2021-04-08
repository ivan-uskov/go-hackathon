package repository

import (
	"context"
	"database/sql"
	"github.com/google/uuid"
	"go-hackaton/src/pkg/common/infrastructure/repository"
	"go-hackaton/src/pkg/sessions/model"
	"time"
)

func NewParticipantRepository(db *sql.DB) model.ParticipantRepository {
	return &participantRepository{repository.Database{DB: db}}
}

type participantRepository struct {
	db repository.Database
}

func (pr *participantRepository) Add(p model.Participant) error {
	return pr.db.Tx(func(tx *sql.Tx, ctx context.Context, closeTx func(error) error) error {
		_, err := tx.ExecContext(
			ctx,
			"INSERT INTO `session_participant` (`participant_id`, `session_id`, `name`, `name_hash`, `endpoint`, `score`, `created_at`, `scored_at`) VALUES (UUID_TO_BIN(?), UUID_TO_BIN(?), ?, UNHEX(MD5(?)), ?, ?, ?, ?) "+
				"ON DUPLICATE KEY UPDATE `participant_id` = UUID_TO_BIN(?), `session_id` = UUID_TO_BIN(?), `name` = ?, `name_hash` = UNHEX(MD5(?)), `endpoint` = ?, `score` = ?, `created_at` = ?, `scored_at` = ?",
			p.ID, p.SessionID, p.Name, p.Name, p.Endpoint, p.Score, p.CreatedAt, p.ScoredAt,
			p.ID, p.SessionID, p.Name, p.Name, p.Endpoint, p.Score, p.CreatedAt, p.ScoredAt)

		return closeTx(err)
	})
}

func (pr *participantRepository) Get(id uuid.UUID) (*model.Participant, error) {
	rows, err := pr.db.Query(""+
		"SELECT "+
		"BIN_TO_UUID(sp.participant_id) AS participant_id, "+
		"BIN_TO_UUID(sp.session_id) AS session_id, "+
		"sp.name, "+
		"sp.endpoint, "+
		"sp.score, "+
		"sp.created_at, "+
		"sp.scored_at "+
		"FROM `session_participant` sp "+
		"WHERE sp.participant_id = UUID_TO_BIN(?) ", id)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	if rows.Next() {
		return parseParticipant(rows)
	}

	return nil, nil // not found
}

func (pr *participantRepository) GetByName(name string) (*model.Participant, error) {
	rows, err := pr.db.Query(""+
		"SELECT "+
		"BIN_TO_UUID(sp.participant_id) AS participant_id, "+
		"BIN_TO_UUID(sp.session_id) AS session_id, "+
		"sp.name, "+
		"sp.endpoint, "+
		"sp.score, "+
		"sp.created_at, "+
		"sp.scored_at "+
		"FROM `session_participant` sp "+
		"WHERE sp.name = ? ", name)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	if rows.Next() {
		return parseParticipant(rows)
	}

	return nil, nil // not found
}

func parseParticipant(r *sql.Rows) (*model.Participant, error) {
	var participantId string
	var sessionId string
	var name string
	var endpoint string
	var score int
	var createdAt time.Time
	var scoredAtNullable sql.NullTime

	err := r.Scan(&participantId, &sessionId, &name, &endpoint, &score, &createdAt, &scoredAtNullable)
	if err != nil {
		return nil, err
	}

	participantUid, err := uuid.Parse(participantId)
	if err != nil {
		return nil, err
	}

	sessionUid, err := uuid.Parse(sessionId)
	if err != nil {
		return nil, err
	}

	var scoredAt *time.Time
	if scoredAtNullable.Valid {
		scoredAt = &scoredAtNullable.Time
	} else {
		scoredAt = nil
	}

	return &model.Participant{
		ID:        participantUid,
		SessionID: sessionUid,
		Name:      name,
		Endpoint:  endpoint,
		Score:     score,
		CreatedAt: createdAt,
		ScoredAt:  scoredAt,
	}, nil
}
