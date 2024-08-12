// Package irepo provides interfaces for task repository operations.
package irepo

import (
	taskmodel "github.com/beka-birhanu/task_manager_final/domain/models/task"
	"github.com/google/uuid"
)

// Task defines methods to manage tasks in the store.
type Task interface {

	// Save adds a new task if it doesnot exist else updates the existing one.
	Save(task *taskmodel.Task) error

	// Delete removes a task by ID.
	Delete(id uuid.UUID) error

	// GetAll retrieves all tasks.
	GetAll() ([]*taskmodel.Task, error)

	// GetSingle returns a task by ID.
	GetSingle(id uuid.UUID) (*taskmodel.Task, error)
}
