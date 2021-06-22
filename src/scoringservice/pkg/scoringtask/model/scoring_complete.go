package model

import (
	"github.com/google/uuid"
)

const ScoringCompleteEventType = "ScoringCompleteEvent"

type ScoringCompleteEvent struct {
	SolutionID uuid.UUID
	Score      int
}

func (s *ScoringCompleteEvent) GetType() string {
	return ScoringCompleteEventType
}
