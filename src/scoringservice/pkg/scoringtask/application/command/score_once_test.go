package command

import (
	"go-hackathon/src/scoringservice/pkg/scoringtask/application/adapter"
	"testing"
)

const mockScore = 5

type mockScorerFactory struct{}

func (m *mockScorerFactory) GetScorer(_ string) (adapter.Scorer, error) {
	return m, nil
}

func (m *mockScorerFactory) Score(_ string) int {
	return mockScore
}

func TestScoreOnceCommand(t *testing.T) {
	uow := &mockUnitOfWork{}
	_ = uow.Add(mockScoringTask)
	h := NewScoreOnceCommandHandler(uow, &mockScorerFactory{})
	err := h.Handle()
	if err != nil {
		t.Errorf("Score failde %v", err)
	}
	task, _ := uow.Get(mockScoringTask.ID)
	if task == nil || task.Score != mockScore || task.ScoredAt == nil {
		t.Error("score not updated")
	}
}
