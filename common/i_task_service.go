package common

import (
	"time"

	"github.com/beka-birhanu/models/taskmodel"
	"github.com/google/uuid"
)

type ITaskService interface {

	// Add adds a new task to the store. Returns an error if there is an ID conflict.
	Add(title, description string, dueDate time.Time, status taskmodel.Status) (*taskmodel.Task, error)

	// Update updates an existing task. Returns an error if the task is not found.
	Update(id uuid.UUID, title, description string, dueDate time.Time, status taskmodel.Status) (*taskmodel.Task, error)

	// Delete removes a task by ID. Returns an error if the task is not found.
	Delete(id uuid.UUID) error

	// GetAll retrieves all tasks from the MongoDB collection.
	//
	// Returns a slice of pointers to `taskmodel.Task` and an error if there is a connection or query issue with database.
	GetAll() ([]*taskmodel.Task, error)

	// GetSingle returns a task by ID. Returns an error if the task is not found.
	GetSingle(id uuid.UUID) (*taskmodel.Task, error)
}
