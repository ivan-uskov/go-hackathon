package query

import (
	"database/sql"
	log "github.com/sirupsen/logrus"
	"go-hackaton/src/pkg/common/application/errors"
	"go-hackaton/src/pkg/sessions/application/query"
	"go-hackaton/src/pkg/sessions/application/query/data"
	"time"
)

func NewParticipantQueryService(db *sql.DB) query.ParticipantQueryService {
	return &participantQueryService{db: db}
}

type participantQueryService struct {
	db *sql.DB
}

func (qs *participantQueryService) GetParticipants(sessionID string) ([]data.ParticipantData, error) {
	rows, err := qs.db.Query(""+
		"SELECT "+
		"BIN_TO_UUID(sp.participant_id) AS participant_id, "+
		"sp.name, "+
		"sp.score, "+
		"sp.endpoint, "+
		"sp.created_at, "+
		"sp.scored_at "+
		"FROM `session_participant` sp "+
		"WHERE sp.session_id = UUID_TO_BIN(?) ", sessionID)

	if err != nil {
		log.Error(err)
		return nil, errors.InternalError
	}
	defer rows.Close()

	return parseParticipants(rows)
}

func (qs *participantQueryService) GetFirstScoredParticipantBefore(time time.Time) (*data.ParticipantData, error) {
	rows, err := qs.db.Query(""+
		"SELECT "+
		"BIN_TO_UUID(sp.participant_id) AS participant_id, "+
		"sp.name, "+
		"sp.score, "+
		"sp.endpoint, "+
		"sp.created_at, "+
		"sp.scored_at "+
		"FROM `session_participant` sp "+
		"WHERE scored_at IS NULL OR scored_at < ? "+
		"ORDER BY scored_at ASC "+
		"LIMIT 1", time)

	if err != nil {
		log.Error(err)
		return nil, errors.InternalError
	}
	defer rows.Close()

	if rows.Next() {
		return parseParticipant(rows)
	}

	return nil, nil
}

func parseParticipants(rows *sql.Rows) ([]data.ParticipantData, error) {
	participants := make([]data.ParticipantData, 0)
	for rows.Next() {
		participant, err := parseParticipant(rows)
		if err != nil {
			log.Error(err)
			return nil, errors.InternalError
		}

		participants = append(participants, *participant)
	}

	return participants, nil
}

func parseParticipant(r *sql.Rows) (*data.ParticipantData, error) {
	var id string
	var name string
	var score int
	var endpoint string
	var createdAt time.Time
	var scoredAtNullable sql.NullTime

	err := r.Scan(&id, &name, &score, &endpoint, &createdAt, &scoredAtNullable)
	if err != nil {
		return nil, err
	}

	var scoredAt *time.Time
	if scoredAtNullable.Valid {
		scoredAt = &scoredAtNullable.Time
	} else {
		scoredAt = nil
	}

	return &data.ParticipantData{
		ID:        id,
		Name:      name,
		Score:     score,
		CreatedAt: createdAt,
		ScoredAt:  scoredAt,
	}, nil
}
