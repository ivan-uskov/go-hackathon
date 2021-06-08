package adapter

type ScoringAdapter interface {
	AddTask(solutionID string, taskType string, endpoint string) error
	RemoveTasks(solutionIDs []string) error

	ValidateTaskType(taskType string) bool
}
