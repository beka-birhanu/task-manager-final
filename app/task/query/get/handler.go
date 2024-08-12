// Package getqry provides the logic to retrieve a specific task by its ID from the repository.
// It includes a handler that processes the Get query and returns the corresponding task.
package getqry

import (
	icmd "github.com/beka-birhanu/task_manager_final/app/common/cqrs/command"
	irepo "github.com/beka-birhanu/task_manager_final/app/common/i_repo"
	taskmodel "github.com/beka-birhanu/task_manager_final/domain/models/task"
	"github.com/google/uuid"
)

// Handler is responsible for handling the Get task query by its ID.
type Handler struct {
	repo irepo.Task // Repository for task-related operations.
}

// Ensure Handler implements the IHandler interface
var _ icmd.IHandler[uuid.UUID, *taskmodel.Task] = &Handler{}

// New creates a new instance of Handler with the provided task repository.
func New(taskRepo irepo.Task) *Handler {
	return &Handler{repo: taskRepo}
}

// Handle processes the Get query by its ID and returns the corresponding task.
func (h *Handler) Handle(id uuid.UUID) (*taskmodel.Task, error) {
	return h.repo.GetSingle(id)
}
