package deletecmd_test

import (
	"errors"
	"testing"

	icmd "github.com/beka-birhanu/task_manager_final/app/common/cqrs/command"
	irepo_mock "github.com/beka-birhanu/task_manager_final/app/common/i_repo/mocks"
	deletecmd "github.com/beka-birhanu/task_manager_final/app/task/command/delete"
	"github.com/google/uuid"
	"github.com/stretchr/testify/suite"
)

// HandlerTestSuite defines the test suite for the deletecmd.Handler.
type HandlerTestSuite struct {
	suite.Suite
	mockRepo *irepo_mock.Task
	handler  icmd.IHandler[uuid.UUID, bool]
	taskID   uuid.UUID
}

// SetupTest sets up the test environment.
func (suite *HandlerTestSuite) SetupTest() {
	// Initialize the mock repository
	suite.mockRepo = new(irepo_mock.Task)

	// Initialize the handler with the mock repository
	suite.handler = deletecmd.New(suite.mockRepo)

	// Initialize a task ID for testing
	suite.taskID = uuid.New()
}

// TestHandle tests the Handle method of the deletecmd.Handler.
func (suite *HandlerTestSuite) TestHandle() {
	// Set up expected behavior for the mock repository
	suite.mockRepo.On("Delete", suite.taskID).Return(nil)

	// Execute the Handle method
	result, err := suite.handler.Handle(suite.taskID)

	// Assertions
	suite.NoError(err)
	suite.True(result)

	// Verify that the Delete method was called on the repository with the expected task ID
	suite.mockRepo.AssertCalled(suite.T(), "Delete", suite.taskID)
	suite.mockRepo.AssertExpectations(suite.T())
}

// TestHandle_ErrorNotFound tests the Handle method when the task to delete is not found.
func (suite *HandlerTestSuite) TestHandle_ErrorNotFound() {
	// Set up expected behavior for the mock repository to return an error indicating task not found
	suite.mockRepo.On("Delete", suite.taskID).Return(errors.New("task not found"))

	// Execute the Handle method
	result, err := suite.handler.Handle(suite.taskID)

	// Assertions
	suite.Error(err)
	suite.False(result)

	// Verify that the Delete method was called on the repository with the expected task ID
	suite.mockRepo.AssertCalled(suite.T(), "Delete", suite.taskID)
	suite.mockRepo.AssertExpectations(suite.T())
}

// Run the test suite
func TestHandlerTestSuite(t *testing.T) {
	suite.Run(t, new(HandlerTestSuite))
}
