package irepo_mock

import (
	taskmodel "github.com/beka-birhanu/task_manager_final/domain/models/task"
	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
)

// Task is a mock implementation of the Task interface using testify.
type Task struct {
	mock.Mock
}

// Save mocks the Save method of the Task interface.
func (m *Task) Save(task *taskmodel.Task) error {
	args := m.Called(task)
	return args.Error(0)
}

// Delete mocks the Delete method of the Task interface.
func (m *Task) Delete(id uuid.UUID) error {
	args := m.Called(id)
	return args.Error(0)
}

// GetAll mocks the GetAll method of the Task interface.
func (m *Task) GetAll() ([]*taskmodel.Task, error) {
	args := m.Called()
	// Ensure to return a slice of *taskmodel.Task and error
	if tasks, ok := args.Get(0).([]*taskmodel.Task); ok {
		return tasks, args.Error(1)
	}
	return nil, args.Error(1)
}

// GetSingle mocks the GetSingle method of the Task interface.
func (m *Task) GetSingle(id uuid.UUID) (*taskmodel.Task, error) {
	args := m.Called(id)
	// Ensure to return *taskmodel.Task and error
	if task, ok := args.Get(0).(*taskmodel.Task); ok {
		return task, args.Error(1)
	}
	return nil, args.Error(1)
}
