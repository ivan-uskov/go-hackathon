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
		"sp.created_at "+
		"FROM `session_participant` sp "+
		"WHERE sp.session_id = UUID_TO_BIN(?) ", sessionID)

	if err != nil {
		log.Error(err)
		return nil, errors.InternalError
	}
	defer rows.Close()

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
	var createdAt time.Time

	err := r.Scan(&id, &name, &score, &createdAt)
	if err != nil {
		return nil, err
	}

	return &data.ParticipantData{
		ID:        id,
		Name:      name,
		Score:     score,
		CreatedAt: createdAt,
	}, nil
}
