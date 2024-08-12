// Package getallqry provides the logic to retrieve all tasks from the repository.
// It includes a handler that processes the GetAll query and returns a list of tasks.
package getallqry

import (
	icmd "github.com/beka-birhanu/task_manager_final/app/common/cqrs/command"
	irepo "github.com/beka-birhanu/task_manager_final/app/common/i_repo"
	taskmodel "github.com/beka-birhanu/task_manager_final/domain/models/task"
)

// Handler is responsible for handling the GetAll tasks query.
type Handler struct {
	repo irepo.Task
}

// Ensure Handler implements the IHandler interface
var _ icmd.IHandler[struct{}, []*taskmodel.Task] = &Handler{}

// New creates a new instance of Handler with the provided task repository.
func New(taskRepo irepo.Task) *Handler {
	return &Handler{repo: taskRepo}
}

// Handle processes the GetAll command and returns a list of tasks.
func (h *Handler) Handle(_ struct{}) ([]*taskmodel.Task, error) {
	return h.repo.GetAll()
}
