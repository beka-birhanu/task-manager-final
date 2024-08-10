// Package updatecmd provides the logic to update an existing task.
// It includes a command structure and a handler to process the update command.
package updatecmd

import (
	icmd "github.com/beka-birhanu/app/common/cqrs/command"
	irepo "github.com/beka-birhanu/app/common/i_repo"
	taskmodel "github.com/beka-birhanu/domain/models/task"
)

type Handler struct {
	repo irepo.Task
}

// Ensure Handler implements icmd.IHandler
var _ icmd.IHandler[*Command, *taskmodel.Task] = &Handler{}

// NewHandler creates a new instance of Handler with the provided task repository.
func NewHandler(taskRepo irepo.Task) *Handler {
	return &Handler{repo: taskRepo}
}

// HandleUpdate handles updating an existing task.
func (h *Handler) Handle(cmd *Command) (*taskmodel.Task, error) {
	task, err := h.repo.GetSingle(cmd.id)
	if err != nil {
		return nil, err
	}

	err = task.Update(taskmodel.Config{
		Title:       cmd.title,
		Description: cmd.description,
		DueDate:     cmd.dueDate,
		Status:      cmd.status,
	})
	if err != nil {
		return nil, err
	}

	err = h.repo.Save(task)
	if err != nil {
		return nil, err
	}
	return task, nil
}
