package query

import (
	"database/sql"
	"go-hackathon/src/common/infrastructure"
	"go-hackathon/src/common/infrastructure/repository"
	"go-hackathon/src/hackathonservice/pkg/hackathon/application/query"
	"go-hackathon/src/hackathonservice/pkg/hackathon/application/query/data"
	"time"
)

func NewParticipantQueryService(db *sql.DB) query.ParticipantQueryService {
	return &participantQueryService{db: db}
}

type participantQueryService struct {
	db *sql.DB
}

func (qs *participantQueryService) GetParticipants(hackathonID string) ([]data.ParticipantData, error) {
	rows, err := qs.db.Query(""+
		"SELECT "+
		"BIN_TO_UUID(hp.participant_id) AS participant_id, "+
		"hp.name, "+
		"hp.score, "+
		"hp.endpoint, "+
		"hp.created_at, "+
		"hp.scored_at "+
		"FROM `hackathon_participant` hp "+
		"WHERE hp.hackathon_id = UUID_TO_BIN(?) ", hackathonID)

	if err != nil {
		return nil, infrastructure.InternalError(err)
	}
	defer infrastructure.CloseRows(rows)

	return parseParticipants(rows)
}

func parseParticipants(rows *sql.Rows) ([]data.ParticipantData, error) {
	participants := make([]data.ParticipantData, 0)
	for rows.Next() {
		participant, err := parseParticipant(rows)
		if err != nil {
			return nil, infrastructure.InternalError(err)
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
		return nil, infrastructure.InternalError(err)
	}

	return &data.ParticipantData{
		ID:        id,
		Name:      name,
		Score:     score,
		Endpoint:  endpoint,
		CreatedAt: createdAt,
		ScoredAt:  repository.TimePointer(scoredAtNullable),
	}, nil
}
