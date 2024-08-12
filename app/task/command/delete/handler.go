// Package deletecmd provides the logic to delete a task.
// It includes the handler to process the delete command.
package deletecmd

import (
	icmd "github.com/beka-birhanu/task_manager_final/app/common/cqrs/command"
	irepo "github.com/beka-birhanu/task_manager_final/app/common/i_repo"
	"github.com/google/uuid"
)

// Handler is responsible for handling the delete task command.
type Handler struct {
	repo irepo.Task // Repository for task-related operations.
}

// Ensure Handler implements the IHandler interface
var _ icmd.IHandler[uuid.UUID, bool] = &Handler{}

// New creates a new instance of Handler with the provided task repository.
func New(taskRepo irepo.Task) *Handler {
	return &Handler{repo: taskRepo}
}

// Handle processes the delete command and removes the task from the repository.
func (h *Handler) Handle(id uuid.UUID) (bool, error) {
	err := h.repo.Delete(id)
	if err != nil {
		return false, err
	}

	return true, nil
}
