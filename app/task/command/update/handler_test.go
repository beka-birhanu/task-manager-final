package updatecmd

import (
	"errors"
	"testing"
	"time"

	icmd "github.com/beka-birhanu/task_manager_final/app/common/cqrs/command"
	irepo_mock "github.com/beka-birhanu/task_manager_final/app/common/i_repo/mocks"
	taskmodel "github.com/beka-birhanu/task_manager_final/domain/models/task"
	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

// HandlerTestSuite defines the test suite for the updatecmd.Handler.
type HandlerTestSuite struct {
	suite.Suite
	mockRepo   *irepo_mock.Task
	handler    icmd.IHandler[*Command, *taskmodel.Task]
	taskID     uuid.UUID
	cmdTitle   string
	cmdDesc    string
	cmdStatus  string
	cmdDueDate time.Time
}

// SetupTest sets up the test environment.
func (suite *HandlerTestSuite) SetupTest() {
	// Initialize the mock repository
	suite.mockRepo = new(irepo_mock.Task)

	// Initialize the handler with the mock repository
	suite.handler = NewHandler(suite.mockRepo)

	// Initialize command properties
	suite.taskID = uuid.New()
	suite.cmdTitle = "Updated Task"
	suite.cmdDesc = "This is an updated task"
	suite.cmdStatus = taskmodel.StatusDone
	suite.cmdDueDate = time.Now().Add(48 * time.Hour)
}

// TestHandle tests the Handle method of the updatecmd.Handler.
func (suite *HandlerTestSuite) TestHandle() {
	// Create an existing task using the New method
	existingTask, _ := taskmodel.New(taskmodel.Config{
		Title:       "Old Task",
		Description: "This is an old task",
		DueDate:     time.Now().Add(24 * time.Hour),
		Status:      taskmodel.StatusPending,
	})

	// Set up mock repository behavior
	suite.mockRepo.On("GetSingle", suite.taskID).Return(existingTask, nil)
	suite.mockRepo.On("Save", mock.AnythingOfType("*taskmodel.Task")).Return(nil)

	// Create the command using the properties stored in the suite
	cmd := NewCommand(suite.taskID, suite.cmdTitle, suite.cmdDesc, suite.cmdStatus, suite.cmdDueDate)

	// Execute the Handle method
	updatedTask, err := suite.handler.Handle(cmd)

	// Assertions
	suite.NoError(err)
	suite.NotNil(updatedTask)
	suite.Equal(suite.cmdTitle, updatedTask.Title())
	suite.Equal(suite.cmdDesc, updatedTask.Description())
	suite.Equal(suite.cmdDueDate, updatedTask.DueDate())
	suite.Equal(suite.cmdStatus, updatedTask.Status())

	// Verify that the Save method was called on the repository with the updated task
	suite.mockRepo.AssertCalled(suite.T(), "Save", mock.AnythingOfType("*taskmodel.Task"))
	suite.mockRepo.AssertExpectations(suite.T())
}

// TestHandle_TaskNotFound tests the Handle method when the task to update is not found.
func (suite *HandlerTestSuite) TestHandle_TaskNotFound() {
	// Set up mock repository behavior for a non-existent task
	suite.mockRepo.On("GetSingle", suite.taskID).Return(nil, errors.New("task not found"))

	// Create the command using the properties stored in the suite
	cmd := NewCommand(suite.taskID, suite.cmdTitle, suite.cmdDesc, suite.cmdStatus, suite.cmdDueDate)

	// Execute the Handle method
	updatedTask, err := suite.handler.Handle(cmd)

	// Assertions
	suite.Error(err)
	suite.Nil(updatedTask)
}

// TestHandle_ErrorSavingTask tests the Handle method when saving the updated task fails.
func (suite *HandlerTestSuite) TestHandle_ErrorSavingTask() {
	// Create an existing task using the New method
	existingTask, _ := taskmodel.New(taskmodel.Config{
		Title:       "Old Task",
		Description: "This is an old task",
		DueDate:     time.Now().Add(24 * time.Hour),
		Status:      taskmodel.StatusPending,
	})

	// Set up mock repository behavior
	suite.mockRepo.On("GetSingle", suite.taskID).Return(existingTask, nil)
	suite.mockRepo.On("Save", mock.AnythingOfType("*taskmodel.Task")).Return(errors.New("failed to save task"))

	// Create the command using the properties stored in the suite
	cmd := NewCommand(suite.taskID, suite.cmdTitle, suite.cmdDesc, suite.cmdStatus, suite.cmdDueDate)

	// Execute the Handle method
	updatedTask, err := suite.handler.Handle(cmd)

	// Assertions
	suite.Error(err)
	suite.Nil(updatedTask)
}

// TestHandle_ErrorRetrievingTask tests the Handle method when retrieving the task fails.
func (suite *HandlerTestSuite) TestHandle_ErrorRetrievingTask() {
	// Set up mock repository behavior for an error retrieving the task
	suite.mockRepo.On("GetSingle", suite.taskID).Return(nil, errors.New("failed to retrieve task"))

	// Create the command using the properties stored in the suite
	cmd := NewCommand(suite.taskID, suite.cmdTitle, suite.cmdDesc, suite.cmdStatus, suite.cmdDueDate)

	// Execute the Handle method
	updatedTask, err := suite.handler.Handle(cmd)

	// Assertions
	suite.Error(err)
	suite.Nil(updatedTask)
}

// Run the test suite
func TestHandlerTestSuite(t *testing.T) {
	suite.Run(t, new(HandlerTestSuite))
}
