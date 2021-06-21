package command

import (
	"github.com/google/uuid"
	"go-hackathon/src/common/application/events"
	"go-hackathon/src/scoringservice/pkg/scoringtask/model"
	"time"
)

var mockScoringTask = model.ScoringTask{
	ID:         uuid.New(),
	SolutionID: uuid.New(),
	Endpoint:   "http://localhost",
	Type:       model.TaskTypeArithmeticExpression,
	CreatedAt:  time.Now(),
}

type mockEventStore struct{}

func (m *mockEventStore) Add(_ events.Event) error {
	return nil
}

type mockUnitOfWork struct {
	tasks map[string]model.ScoringTask
}

func (m *mockUnitOfWork) Execute(f func(r RepositoryProvider) error) error {
	return f(m)
}

func (m *mockUnitOfWork) ScoringTaskRepository() model.ScoringTaskRepository {
	return m
}

func (m *mockUnitOfWork) EventStore() events.EventStore {
	return &mockEventStore{}
}

func (m *mockUnitOfWork) Add(task model.ScoringTask) error {
	if m.tasks == nil {
		m.tasks = make(map[string]model.ScoringTask)
	}

	m.tasks[task.ID.String()] = task
	return nil
}

func (m *mockUnitOfWork) Get(id uuid.UUID) (*model.ScoringTask, error) {
	if m.tasks == nil {
		return nil, nil
	}

	t, found := m.tasks[id.String()]
	if !found {
		return nil, nil
	}

	return &t, nil
}

func (m *mockUnitOfWork) GetBySolutionID(id uuid.UUID) (*model.ScoringTask, error) {
	if m.tasks == nil {
		return nil, nil
	}

	for _, task := range m.tasks {
		if task.SolutionID == id {
			return &task, nil
		}
	}

	return nil, nil
}

func (m *mockUnitOfWork) GetFirstScoringTaskBefore(time time.Time) (*model.ScoringTask, error) {
	if m.tasks == nil {
		return nil, nil
	}

	var t *model.ScoringTask
	for _, task := range m.tasks {
		if task.DeletedAt == nil && (task.ScoredAt == nil || (time.After(*task.ScoredAt) && (t == nil || (t.ScoredAt != nil && t.ScoredAt.After(*task.ScoredAt))))) {
			t = &task
		}
	}

	return t, nil
}
