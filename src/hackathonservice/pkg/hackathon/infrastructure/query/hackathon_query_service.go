package query

import (
	"database/sql"
	"go-hackathon/src/common/infrastructure"
	"go-hackathon/src/common/infrastructure/repository"
	"go-hackathon/src/hackathonservice/pkg/hackathon/application/query"
	"go-hackathon/src/hackathonservice/pkg/hackathon/application/query/data"
	"time"
)

func NewHackathonQueryService(db *sql.DB) query.HackathonQueryService {
	return &hackatonQueryService{db: db}
}

type hackatonQueryService struct {
	db *sql.DB
}

func (qs *hackatonQueryService) GetHackathons() ([]data.HackathonData, error) {
	rows, err := qs.db.Query("" +
		getSelectHackathonSQL() +
		"GROUP BY h.hackathon_id")

	if err != nil {
		return nil, infrastructure.InternalError(err)
	}
	defer infrastructure.CloseRows(rows)

	hackathons := make([]data.HackathonData, 0)
	for rows.Next() {
		hackathon, err := parseHackathon(rows)
		if err != nil {
			return nil, infrastructure.InternalError(err)
		}

		hackathons = append(hackathons, *hackathon)
	}

	return hackathons, nil
}

func (qs *hackatonQueryService) GetHackathon(id string) (*data.HackathonData, error) {
	rows, err := qs.db.Query(""+
		getSelectHackathonSQL()+
		"WHERE h.hackathon_id = UUID_TO_BIN(?)"+
		"GROUP BY h.hackathon_id", id)

	if err != nil {
		return nil, infrastructure.InternalError(err)
	}
	defer infrastructure.CloseRows(rows)

	if rows.Next() {
		hackathon, err := parseHackathon(rows)
		if err != nil {
			return nil, infrastructure.InternalError(err)
		}

		return hackathon, nil
	}

	return nil, nil // not found
}

func getSelectHackathonSQL() string {
	return "" +
		"SELECT " +
		"BIN_TO_UUID(h.hackathon_id) AS hackathon_id, " +
		"h.name, " +
		"COUNT(DISTINCT hp.participant_id), " +
		"h.type, " +
		"h.created_at, " +
		"h.closed_at " +
		"FROM `hackathon` h " +
		"LEFT JOIN hackathon_participant hp ON (h.hackathon_id = hp.hackathon_id) "
}

func parseHackathon(r *sql.Rows) (*data.HackathonData, error) {
	var hackathonId string
	var name string
	var participants int
	var t int
	var createdAt time.Time
	var closedAtNullable sql.NullTime

	err := r.Scan(&hackathonId, &name, &participants, &t, &createdAt, &closedAtNullable)
	if err != nil {
		return nil, infrastructure.InternalError(err)
	}

	return &data.HackathonData{
		ID:           hackathonId,
		Name:         name,
		Participants: participants,
		Type:         t,
		CreatedAt:    createdAt,
		ClosedAt:     repository.TimePointer(closedAtNullable),
	}, nil
}
