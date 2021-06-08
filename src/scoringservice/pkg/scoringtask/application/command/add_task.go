package command

import (
	"github.com/google/uuid"
	"go-hackathon/src/scoringservice/pkg/scoringtask/application/errors"
	"go-hackathon/src/scoringservice/pkg/scoringtask/model"
	"time"
)

type AddTaskCommand struct {
	SolutionID uuid.UUID
	TaskType   string
	Endpoint   string
}

type addTaskCommandHandler struct {
	uow UnitOfWork
}

type AddTaskCommandHandler interface {
	Handle(command AddTaskCommand) error
}

func NewAddTaskCommandHandler(uow UnitOfWork) AddTaskCommandHandler {
	return &addTaskCommandHandler{uow: uow}
}

func (h *addTaskCommandHandler) Handle(command AddTaskCommand) error {
	return h.uow.Execute(func(r model.ScoringTaskRepository) error {
		task, err := r.GetBySolutionID(command.SolutionID)
		if err != nil {
			return err
		}

		if task != nil {
			return errors.TaskAlreadyExistError
		}

		if !model.IsTypeValid(command.TaskType) {
			return errors.InvalidTaskTypeError
		}

		return r.Add(model.ScoringTask{
			ID:         uuid.New(),
			SolutionID: command.SolutionID,
			Endpoint:   command.Endpoint,
			Type:       command.TaskType,
			CreatedAt:  time.Now(),
		})
	})
}
