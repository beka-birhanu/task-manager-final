// Package irepo provides interfaces for task repository operations.
package irepo

import (
	"time"

	taskmodel "github.com/beka-birhanu/domain/models/task"
	"github.com/google/uuid"
)

// Task defines methods to manage tasks in the store.
type Task interface {

	// Add adds a new task to the store.
	Add(title, description, status string, dueDate time.Time) (*taskmodel.Task, error)

	// Update modifies an existing task.
	Update(id uuid.UUID, title, description, status string, dueDate time.Time) (*taskmodel.Task, error)

	// Delete removes a task by ID.
	Delete(id uuid.UUID) error

	// GetAll retrieves all tasks.
	GetAll() ([]*taskmodel.Task, error)

	// GetSingle returns a task by ID.
	GetSingle(id uuid.UUID) (*taskmodel.Task, error)
}

