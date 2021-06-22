package repository

import (
	"database/sql"
	"github.com/google/uuid"
	"go-hackathon/src/common/infrastructure"
	"go-hackathon/src/common/infrastructure/repository"
	"go-hackathon/src/hackathonservice/pkg/hackathon/model"
	"time"
)

type participantRepository struct {
	tx *sql.Tx
}

func (pr *participantRepository) Add(p model.Participant) error {
	_, err := pr.tx.Exec(
		"INSERT INTO `hackathon_participant` (`participant_id`, `hackathon_id`, `name`, `endpoint`, `score`, `created_at`, `scored_at`) VALUES (UUID_TO_BIN(?), UUID_TO_BIN(?), ?, ?, ?, ?, ?) "+
			"ON DUPLICATE KEY UPDATE `hackathon_id` = UUID_TO_BIN(?), `name` = ?, `endpoint` = ?, `score` = ?, `created_at` = ?, `scored_at` = ?",
		p.ID, p.HackathonID, p.Name, p.Endpoint, p.Score, p.CreatedAt, p.ScoredAt,
		p.HackathonID, p.Name, p.Endpoint, p.Score, p.CreatedAt, p.ScoredAt)

	if err != nil {
		err = infrastructure.InternalError(err)
	}

	return err
}

func (pr *participantRepository) Get(id uuid.UUID) (*model.Participant, error) {
	rows, err := pr.tx.Query(""+
		getSelectParticipantSQL()+
		"WHERE hp.participant_id = UUID_TO_BIN(?) ", id)

	if err != nil {
		return nil, infrastructure.InternalError(err)
	}
	defer infrastructure.Close(rows)

	if rows.Next() {
		return parseParticipant(rows)
	}

	return nil, nil // not found
}

func (pr *participantRepository) GetByNameAndHackathonID(name string, hackathonID uuid.UUID) (*model.Participant, error) {
	rows, err := pr.tx.Query(""+
		getSelectParticipantSQL()+
		"WHERE hp.name = ? AND hp.hackathon_id = UUID_TO_BIN(?)", name, hackathonID)

	if err != nil {
		return nil, infrastructure.InternalError(err)
	}
	defer infrastructure.Close(rows)

	if rows.Next() {
		return parseParticipant(rows)
	}

	return nil, nil // not found
}

func (pr *participantRepository) GetByHackathonID(hackathonID uuid.UUID) ([]model.Participant, error) {
	rows, err := pr.tx.Query(""+
		getSelectParticipantSQL()+
		"WHERE hp.hackathon_id = UUID_TO_BIN(?)", hackathonID.String())

	if err != nil {
		return nil, infrastructure.InternalError(err)
	}
	defer infrastructure.Close(rows)

	var pp []model.Participant
	for rows.Next() {
		part, err := parseParticipant(rows)
		if err != nil {
			return nil, err
		}

		pp = append(pp, *part)
	}

	return pp, nil // not found
}

func getSelectParticipantSQL() string {
	return "" +
		"SELECT " +
		"BIN_TO_UUID(hp.participant_id) AS participant_id, " +
		"BIN_TO_UUID(hp.hackathon_id) AS hackathon_id, " +
		"hp.name, " +
		"hp.endpoint, " +
		"hp.score, " +
		"hp.created_at, " +
		"hp.scored_at " +
		"FROM `hackathon_participant` hp "
}

func parseParticipant(r *sql.Rows) (*model.Participant, error) {
	var participantId string
	var hackathonId string
	var name string
	var endpoint string
	var score int
	var createdAt time.Time
	var scoredAtNullable sql.NullTime

	err := r.Scan(&participantId, &hackathonId, &name, &endpoint, &score, &createdAt, &scoredAtNullable)
	if err != nil {
		return nil, infrastructure.InternalError(err)
	}

	participantUid, err := uuid.Parse(participantId)
	if err != nil {
		return nil, infrastructure.InternalError(err)
	}

	hackathonUid, err := uuid.Parse(hackathonId)
	if err != nil {
		return nil, infrastructure.InternalError(err)
	}

	return &model.Participant{
		ID:          participantUid,
		HackathonID: hackathonUid,
		Name:        name,
		Endpoint:    endpoint,
		Score:       score,
		CreatedAt:   createdAt,
		ScoredAt:    repository.TimePointer(scoredAtNullable),
	}, nil
}
