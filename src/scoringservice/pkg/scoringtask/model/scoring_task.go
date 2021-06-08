package model

import (
	"github.com/google/uuid"
	"time"
)

type ScoringTask struct {
	ID         uuid.UUID
	SolutionID uuid.UUID
	Endpoint   string
	Type       string
	Score      int
	CreatedAt  time.Time
	ScoredAt   *time.Time
	DeletedAt  *time.Time
}

type ScoringTaskRepository interface {
	Add(task ScoringTask) error
	Get(id uuid.UUID) (*ScoringTask, error)
	GetBySolutionID(id uuid.UUID) (*ScoringTask, error)
}

func (st *ScoringTask) UpdateScore(score int) {
	now := time.Now()

	st.Score = score
	st.ScoredAt = &now
}

func (st *ScoringTask) Delete() {
	now := time.Now()
	st.DeletedAt = &now
}

func (st *ScoringTask) IsDeleted() bool {
	return st.DeletedAt != nil
}
