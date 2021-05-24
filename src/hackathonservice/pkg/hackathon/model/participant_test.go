package model

import (
	"testing"
	"time"
)

func TestUpdateParticipantScore(t *testing.T) {
	p := Participant{
		Score: 0,
	}
	now := time.Now()
	p.UpdateScore(5)
	if p.Score != 5 {
		t.Error("Score not updated")
	}
	if p.ScoredAt == nil {
		t.Error("ScoredAt not set")
	}
	if !p.ScoredAt.After(now) {
		t.Error("ScoredAt is less than time before update score")
	}
}
