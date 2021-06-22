package command

import (
	"fmt"
	"github.com/google/uuid"
)

func getSolutionIDLock(id uuid.UUID) string {
	return fmt.Sprintf("scoring_task_solution_id_%s", id)
}
