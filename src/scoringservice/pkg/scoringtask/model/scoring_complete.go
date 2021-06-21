package model

import (
	"github.com/google/uuid"
	"time"
)

const ScoringCompleteEventType = "ScoringCompleteEvent"

type ScoringCompleteEvent struct {
	SolutionID uuid.UUID
	Score      int
	ScoredAt   time.Time
}

func (s *ScoringCompleteEvent) GetType() string {
	return ScoringCompleteEventType
}
