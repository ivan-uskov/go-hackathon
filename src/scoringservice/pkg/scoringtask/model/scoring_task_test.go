package model

import (
	"testing"
	"time"
)

func TestDeleteScoringTask(t *testing.T) {
	st := ScoringTask{}

	now := time.Now()
	st.Delete()

	if !st.IsDeleted() {
		t.Error("Delete scoring task not working")
	}
	if !st.DeletedAt.After(now) {
		t.Error("DeletedAt is less then time before delete scoring task called")
	}
}

func TestUpdateScore(t *testing.T) {
	st := ScoringTask{}

	now := time.Now()
	st.UpdateScore(5)

	if st.Score != 5 {
		t.Error("Update score not working")
	}
	if !st.ScoredAt.After(now) {
		t.Error("ScoredAt is less then time before update score called")
	}
}
