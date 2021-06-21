package command

import (
	"github.com/google/uuid"
	"go-hackathon/src/scoringservice/pkg/scoringtask/application/errors"
)

type RemoveTasksCommand struct {
	SolutionIDs []uuid.UUID
}

type removeTasksCommandHandler struct {
	uow UnitOfWork
}

type RemoveTasksCommandHandler interface {
	Handle(command RemoveTasksCommand) error
}

func NewRemoveTasksCommandHandler(uow UnitOfWork) RemoveTasksCommandHandler {
	return &removeTasksCommandHandler{uow}
}

func (h *removeTasksCommandHandler) Handle(command RemoveTasksCommand) error {
	return h.uow.Execute(func(rp RepositoryProvider) error {
		r := rp.ScoringTaskRepository()

		for _, solutionId := range command.SolutionIDs {
			task, err := r.GetBySolutionID(solutionId)
			if err != nil {
				return err
			}

			if task == nil {
				return errors.TaskNotExistError
			}

			task.Delete()

			err = r.Add(*task)
			if err != nil {
				return err
			}
		}

		return nil
	})
}
