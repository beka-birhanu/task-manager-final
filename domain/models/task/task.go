/*
Package taskmodel provides the `Task` aggregate, which represents a task with
a title, description, due date, and status. The package includes functionality
for creating, updating, and converting tasks to and from BSON format for MongoDB operations.

Key Components:
  - Task: Represents a task with an ID, title, description, due date, and status.
  - TaskConfig: Holds parameters for creating or updating a Task.
  - New: Creates a new Task with validation and generates a unique ID.
  - TaskBSON: Represents the BSON format of a Task for MongoDB operations.
  - ToBSON: Converts a Task to its BSON representation.
  - FromBSON: Converts a BSON representation back to a Task.

Dependencies:
- github.com/google/uuid: For generating unique task IDs.
- github.com/beka-birhanu/errors: For handling validation errors.
*/
package taskmodel

import (
	"time"

	errdmn "github.com/beka-birhanu/domain/errors"
	"github.com/google/uuid"
)

const (
	StatusDone       = "done"
	StatusInProgress = "inprogress"
	StatusPending    = "pending"
)

// Task represents a task with an ID, title, description, due date, and status.
type Task struct {
	id          uuid.UUID
	title       string
	description string
	dueDate     time.Time
	status      string
}

// TaskBSON represents the BSON format of a Task for MongoDB operations.
type TaskBSON struct {
	ID          uuid.UUID `bson:"_id"`
	Title       string    `bson:"title"`
	Description string    `bson:"description"`
	DueDate     time.Time `bson:"dueDate"`
	Status      string    `bson:"status"`
	UpdatedAt   time.Time `bson:"updatedAt"`
}

// ToBSON converts a Task to a TaskBSON.
func (t *Task) ToBSON() *TaskBSON {
	return &TaskBSON{
		ID:          t.ID(),
		Title:       t.Title(),
		Description: t.Description(),
		DueDate:     t.DueDate(),
		Status:      t.Status(),
		UpdatedAt:   time.Now(),
	}
}

// FromBSON converts a TaskBSON to a Task.
func FromBSON(bson *TaskBSON) *Task {
	return &Task{
		id:          bson.ID,
		title:       bson.Title,
		description: bson.Description,
		dueDate:     bson.DueDate,
		status:      bson.Status,
	}
}

// TaskConfig represents the configuration for creating or updating a Task.
type TaskConfig struct {
	Title       string
	Description string
	DueDate     time.Time
	Status      string
}

// New creates a new Task with the given configuration, validates its properties, and generates an ID.
func New(config TaskConfig) (*Task, error) {
	if err := validateTaskConfig(config); err != nil {
		return nil, err
	}

	return &Task{
		id:          uuid.New(),
		title:       config.Title,
		description: config.Description,
		dueDate:     config.DueDate,
		status:      config.Status,
	}, nil
}

// validateTaskConfig checks if the provided task configuration is valid.
func validateTaskConfig(config TaskConfig) error {
	if config.Title == "" {
		return errdmn.TitleEmpty
	}
	if config.Description == "" {
		return errdmn.DescriptionEmpty
	}
	if config.DueDate.IsZero() {
		return errdmn.DueDateZero
	}
	if !isValidStatus(config.Status) {
		return errdmn.InvalidStatus
	}
	return nil
}

// isValidStatus checks if the given status is one of the allowed statuses.
func isValidStatus(status string) bool {
	switch status {
	case StatusDone, StatusInProgress, StatusPending:
		return true
	default:
		return false
	}
}

// ID returns the task's ID.
func (t *Task) ID() uuid.UUID {
	return t.id
}

// Title returns the task's title.
func (t *Task) Title() string {
	return t.title
}

// Description returns the task's description.
func (t *Task) Description() string {
	return t.description
}

// DueDate returns the task's due date.
func (t *Task) DueDate() time.Time {
	return t.dueDate
}

// Status returns the task's status.
func (t *Task) Status() string {
	return t.status
}

// Update updates the task's fields with the provided configuration after validating the data.
func (t *Task) Update(config TaskConfig) error {
	if err := validateTaskConfig(config); err != nil {
		return err
	}

	t.title = config.Title
	t.description = config.Description
	t.dueDate = config.DueDate
	t.status = config.Status
	return nil
}

