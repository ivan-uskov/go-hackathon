package repository

import (
	"database/sql"
	"github.com/google/uuid"
	"go-hackathon/src/common/infrastructure"
	"go-hackathon/src/common/infrastructure/repository"
	"go-hackathon/src/hackathonservice/pkg/hackathon/model"
	"time"
)

type hackathonRepository struct {
	tx *sql.Tx
}

func (hr *hackathonRepository) Add(h model.Hackathon) error {
	_, err := hr.tx.Exec(
		"INSERT INTO `hackathon` (`hackathon_id`, `name`, `type`, `created_at`) VALUES (UUID_TO_BIN(?), ?, ?, ?)"+
			"ON DUPLICATE KEY UPDATE `name` = ?, `type` = ?, `created_at` = ?, `closed_at` = ?",
		h.ID, h.Name, h.Type, h.CreatedAt,
		h.Name, h.Type, h.CreatedAt, h.ClosedAt)

	if err != nil {
		err = infrastructure.InternalError(err)
	}

	return err
}

func (hr *hackathonRepository) Get(id uuid.UUID) (*model.Hackathon, error) {
	rows, err := hr.tx.Query(""+
		getSelectHackathonSQL()+
		"WHERE h.hackathon_id = UUID_TO_BIN(?) ", id)

	if err != nil {
		return nil, infrastructure.InternalError(err)
	}
	defer infrastructure.CloseRows(rows)

	if rows.Next() {
		return parseHackathon(rows)
	}

	return nil, nil // not found
}

func (hr *hackathonRepository) GetByName(name string) (*model.Hackathon, error) {
	rows, err := hr.tx.Query(""+
		getSelectHackathonSQL()+
		"WHERE h.name = ? ", name)

	if err != nil {
		return nil, infrastructure.InternalError(err)
	}
	defer infrastructure.CloseRows(rows)

	if rows.Next() {
		return parseHackathon(rows)
	}

	return nil, nil // not found
}

func getSelectHackathonSQL() string {
	return "" +
		"SELECT " +
		"BIN_TO_UUID(h.hackathon_id) AS hackathon_id, " +
		"h.name, " +
		"h.type, " +
		"h.created_at, " +
		"h.closed_at " +
		"FROM `hackathon` h "
}

func parseHackathon(r *sql.Rows) (*model.Hackathon, error) {
	var hackathonId string
	var name string
	var t string
	var createdAt time.Time
	var closedAtNullable sql.NullTime

	err := r.Scan(&hackathonId, &name, &t, &createdAt, &closedAtNullable)
	if err != nil {
		return nil, infrastructure.InternalError(err)
	}

	hackathonUid, err := uuid.Parse(hackathonId)
	if err != nil {
		return nil, infrastructure.InternalError(err)
	}

	return &model.Hackathon{
		ID:        hackathonUid,
		Name:      name,
		Type:      t,
		CreatedAt: createdAt,
		ClosedAt:  repository.TimePointer(closedAtNullable),
	}, nil
}
