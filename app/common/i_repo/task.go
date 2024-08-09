package irepo

import (
	"time"

	taskmodel "github.com/beka-birhanu/domain/models/task"
	taskrepo "github.com/beka-birhanu/infrastructure/repo/task"
	"github.com/google/uuid"
)

type Task interface {

	// Add adds a new task to the store. Returns an error if there is an ID conflict.
	Add(title, description, status string, dueDate time.Time) (*taskrepo.Repo, error)

	// Update updates an existing task. Returns an error if the task is not found.
	Update(id uuid.UUID, title, description, status string, dueDate time.Time) (*taskmodel.Task, error)

	// Delete removes a task by ID. Returns an error if the task is not found.
	Delete(id uuid.UUID) error

	// GetAll retrieves all tasks from the MongoDB collection.
	//
	// Returns a slice of pointers to `taskmodel.Task` and an error if there is a connection or query issue with database.
	GetAll() ([]*taskmodel.Task, error)

	// GetSingle returns a task by ID. Returns an error if the task is not found.
	GetSingle(id uuid.UUID) (*taskmodel.Task, error)
}
