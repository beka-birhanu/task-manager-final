package getqry_test

import (
	"errors"
	"testing"
	"time"

	icmd "github.com/beka-birhanu/task_manager_final/app/common/cqrs/command"
	irepo_mock "github.com/beka-birhanu/task_manager_final/app/common/i_repo/mocks"
	getqry "github.com/beka-birhanu/task_manager_final/app/task/query/get"
	taskmodel "github.com/beka-birhanu/task_manager_final/domain/models/task"
	"github.com/google/uuid"
	"github.com/stretchr/testify/suite"
)

// HandlerTestSuite defines the test suite for the getqry.Handler.
type HandlerTestSuite struct {
	suite.Suite
	mockRepo *irepo_mock.Task
	handler  icmd.IHandler[uuid.UUID, *taskmodel.Task]
	taskID   uuid.UUID
}

// SetupTest sets up the test environment.
func (suite *HandlerTestSuite) SetupTest() {
	// Initialize the mock repository
	suite.mockRepo = new(irepo_mock.Task)

	// Initialize the handler with the mock repository
	suite.handler = getqry.New(suite.mockRepo)

	// Initialize a task ID for testing
	suite.taskID = uuid.New()
}

// TestHandle tests the Handle method of the getqry.Handler for successful retrieval.
func (suite *HandlerTestSuite) TestHandle_Success() {
	// Create a mock task
	expectedTask, _ := taskmodel.New(taskmodel.Config{
		Title:       "Mock Task",
		Description: "This is a mock task",
		DueDate:     time.Now().Add(24 * time.Hour),
		Status:      taskmodel.StatusPending,
	})

	// Set up expected behavior for the mock repository
	suite.mockRepo.On("GetSingle", suite.taskID).Return(expectedTask, nil)

	// Execute the Handle method
	task, err := suite.handler.Handle(suite.taskID)

	// Assertions
	suite.NoError(err)
	suite.NotNil(task)
	suite.Equal(expectedTask.Title(), task.Title())
	suite.Equal(expectedTask.Description(), task.Description())
	suite.Equal(expectedTask.DueDate(), task.DueDate())
	suite.Equal(expectedTask.Status(), task.Status())

	// Verify that the GetSingle method was called on the repository with the expected ID
	suite.mockRepo.AssertCalled(suite.T(), "GetSingle", suite.taskID)
	suite.mockRepo.AssertExpectations(suite.T())
}

// TestHandle_ErrorRetrievingTask tests the Handle method when retrieving the task fails.
func (suite *HandlerTestSuite) TestHandle_ErrorRetrievingTask() {
	// Set up expected behavior for the mock repository
	suite.mockRepo.On("GetSingle", suite.taskID).Return(nil, errors.New("failed to retrieve task"))

	// Execute the Handle method
	task, err := suite.handler.Handle(suite.taskID)

	// Assertions
	suite.Error(err)
	suite.Nil(task)

	// Verify that the GetSingle method was called on the repository with the expected ID
	suite.mockRepo.AssertCalled(suite.T(), "GetSingle", suite.taskID)
	suite.mockRepo.AssertExpectations(suite.T())
}

// TestHandle_TaskNotFound tests the Handle method when the task is not found in the repository.
func (suite *HandlerTestSuite) TestHandle_TaskNotFound() {
	// Set up expected behavior for the mock repository
	suite.mockRepo.On("GetSingle", suite.taskID).Return(nil, nil)

	// Execute the Handle method
	task, err := suite.handler.Handle(suite.taskID)

	// Assertions
	suite.NoError(err)
	suite.Nil(task)

	// Verify that the GetSingle method was called on the repository with the expected ID
	suite.mockRepo.AssertCalled(suite.T(), "GetSingle", suite.taskID)
	suite.mockRepo.AssertExpectations(suite.T())
}

// Run the test suite
func TestHandlerTestSuite(t *testing.T) {
	suite.Run(t, new(HandlerTestSuite))
}
