package command

import (
	"github.com/google/uuid"
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

type mockUnitOfWork struct {
	tasks map[string]model.ScoringTask
}

func (m *mockUnitOfWork) Execute(f func(r model.ScoringTaskRepository) error) error {
	return f(m)
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
