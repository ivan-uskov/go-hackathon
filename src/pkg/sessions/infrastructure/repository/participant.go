package repository

import (
	"database/sql"
	"github.com/google/uuid"
	"go-hackaton/src/pkg/common/infrastructure"
	"go-hackaton/src/pkg/common/infrastructure/repository"
	"go-hackaton/src/pkg/sessions/model"
	"time"
)

type participantRepository struct {
	tx *sql.Tx
}

func (pr *participantRepository) Add(p model.Participant) error {
	_, err := pr.tx.Exec(
		"INSERT INTO `session_participant` (`participant_id`, `session_id`, `name`, `name_hash`, `endpoint`, `score`, `created_at`, `scored_at`) VALUES (UUID_TO_BIN(?), UUID_TO_BIN(?), ?, UNHEX(MD5(?)), ?, ?, ?, ?) "+
			"ON DUPLICATE KEY UPDATE `participant_id` = UUID_TO_BIN(?), `session_id` = UUID_TO_BIN(?), `name` = ?, `name_hash` = UNHEX(MD5(?)), `endpoint` = ?, `score` = ?, `created_at` = ?, `scored_at` = ?",
		p.ID, p.SessionID, p.Name, p.Name, p.Endpoint, p.Score, p.CreatedAt, p.ScoredAt,
		p.ID, p.SessionID, p.Name, p.Name, p.Endpoint, p.Score, p.CreatedAt, p.ScoredAt)

	return infrastructure.InternalError(err)
}

func (pr *participantRepository) Get(id uuid.UUID) (*model.Participant, error) {
	rows, err := pr.tx.Query(""+
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
		return nil, infrastructure.InternalError(err)
	}
	defer infrastructure.CloseRows(rows)

	if rows.Next() {
		return parseParticipant(rows)
	}

	return nil, nil // not found
}

func (pr *participantRepository) GetByName(name string) (*model.Participant, error) {
	rows, err := pr.tx.Query(""+
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
		return nil, infrastructure.InternalError(err)
	}
	defer infrastructure.CloseRows(rows)

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
		return nil, infrastructure.InternalError(err)
	}

	participantUid, err := uuid.Parse(participantId)
	if err != nil {
		return nil, infrastructure.InternalError(err)
	}

	sessionUid, err := uuid.Parse(sessionId)
	if err != nil {
		return nil, infrastructure.InternalError(err)
	}

	return &model.Participant{
		ID:        participantUid,
		SessionID: sessionUid,
		Name:      name,
		Endpoint:  endpoint,
		Score:     score,
		CreatedAt: createdAt,
		ScoredAt:  repository.TimePointer(scoredAtNullable),
	}, nil
}
