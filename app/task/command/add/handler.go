// Package addcmd provides the logic for adding new tasks.
// It includes the command structure and the handler to process the add task command.
package addcmd

import (
	icmd "github.com/beka-birhanu/task_manager_final/app/common/cqrs/command"
	irepo "github.com/beka-birhanu/task_manager_final/app/common/i_repo"
	taskmodel "github.com/beka-birhanu/task_manager_final/domain/models/task"
)

// Handler handles the logic for adding a new task to the repository.
type Handler struct {
	repo irepo.Task
}

// Ensure Handler implements icmd.IHandler
var _ icmd.IHandler[*Command, *taskmodel.Task] = &Handler{}

// NewHandler creates a new instance of Handler with the given task repository.
func NewHandler(repo irepo.Task) *Handler {
	return &Handler{repo: repo}
}

// Handle processes the command to add a new task to the repository.
func (h *Handler) Handle(cmd *Command) (*taskmodel.Task, error) {
	task, err := taskmodel.New(taskmodel.Config{
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
