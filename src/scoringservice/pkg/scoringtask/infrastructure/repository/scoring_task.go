package repository

import (
	"database/sql"
	"github.com/google/uuid"
	"go-hackathon/src/common/infrastructure"
	"go-hackathon/src/common/infrastructure/repository"
	"go-hackathon/src/scoringservice/pkg/scoringtask/model"
	"time"
)

type scoringTaskRepository struct {
	tx *sql.Tx
}

func (s *scoringTaskRepository) Add(task model.ScoringTask) error {
	_, err := s.tx.Exec(
		"INSERT INTO `scoring_task` (`task_id`, `solution_id`, `endpoint`, `type`, `score`, `created_at`, `scored_at`, `deleted_at`) VALUES (UUID_TO_BIN(?), UUID_TO_BIN(?), ?, ?, ?, ?, ?, ?)"+
			"ON DUPLICATE KEY UPDATE `solution_id` = UUID_TO_BIN(?), `endpoint` = ?, `type` = ?, `score` = ?, `scored_at` = ?, `deleted_at` = ?",
		task.ID, task.SolutionID, task.Endpoint, task.Type, task.Score, task.CreatedAt, task.ScoredAt, task.DeletedAt,
		task.SolutionID, task.Endpoint, task.Type, task.Score, task.ScoredAt, task.DeletedAt)

	if err != nil {
		err = infrastructure.InternalError(err)
	}

	return err
}

func (s *scoringTaskRepository) Get(id uuid.UUID) (*model.ScoringTask, error) {
	rows, err := s.tx.Query(""+
		selectScoringTaskSql()+
		"WHERE st.task_id = UUID_TO_BIN(?) ", id)

	if err != nil {
		return nil, infrastructure.InternalError(err)
	}
	defer infrastructure.CloseRows(rows)

	if rows.Next() {
		return parseScoringTask(rows)
	}

	return nil, nil // not found
}

func (s *scoringTaskRepository) GetBySolutionID(id uuid.UUID) (*model.ScoringTask, error) {
	rows, err := s.tx.Query(""+
		selectScoringTaskSql()+
		"WHERE st.solution_id = UUID_TO_BIN(?) ", id)

	if err != nil {
		return nil, infrastructure.InternalError(err)
	}
	defer infrastructure.CloseRows(rows)

	if rows.Next() {
		return parseScoringTask(rows)
	}

	return nil, nil // not found
}

func selectScoringTaskSql() string {
	return "" +
		"SELECT " +
		"BIN_TO_UUID(st.task_id) AS task_id, " +
		"BIN_TO_UUID(st.solution_id) AS solution_id, " +
		"st.endpoint, " +
		"st.type, " +
		"st.score, " +
		"st.created_at, " +
		"st.scored_at, " +
		"st.deleted_at " +
		"FROM `scoring_task` st "
}

func parseScoringTask(rows *sql.Rows) (*model.ScoringTask, error) {
	var id string
	var solutionID string
	var endpoint string
	var t string
	var score int
	var createdAt time.Time
	var scoredAtNullable sql.NullTime
	var deletedAtNullable sql.NullTime

	err := rows.Scan(&id, &solutionID, &endpoint, &t, &score, &createdAt, &scoredAtNullable, &deletedAtNullable)
	if err != nil {
		return nil, infrastructure.InternalError(err)
	}

	uid, err := uuid.Parse(id)
	if err != nil {
		return nil, infrastructure.InternalError(err)
	}

	solutionUid, err := uuid.Parse(solutionID)
	if err != nil {
		return nil, infrastructure.InternalError(err)
	}

	return &model.ScoringTask{
		ID:         uid,
		SolutionID: solutionUid,
		Endpoint:   endpoint,
		Type:       t,
		Score:      score,
		CreatedAt:  createdAt,
		ScoredAt:   repository.TimePointer(scoredAtNullable),
		DeletedAt:  repository.TimePointer(deletedAtNullable),
	}, nil
}
