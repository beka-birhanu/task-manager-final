package addcmd_test

import (
	"errors"
	"testing"
	"time"

	icmd "github.com/beka-birhanu/task_manager_final/app/common/cqrs/command"
	"github.com/beka-birhanu/task_manager_final/app/common/i_repo/mocks"
	addcmd "github.com/beka-birhanu/task_manager_final/app/task/command/add"
	taskmodel "github.com/beka-birhanu/task_manager_final/domain/models/task"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

// HandlerTestSuite defines the test suite for the addcmd.Handler.
type HandlerTestSuite struct {
	suite.Suite
	mockRepo   *irepo_mock.Task
	handler    icmd.IHandler[*addcmd.Command, *taskmodel.Task]
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
	suite.handler = addcmd.NewHandler(suite.mockRepo)

	// Initialize the command properties
	suite.cmdTitle = "Test Task"
	suite.cmdDesc = "This is a test task"
	suite.cmdStatus = taskmodel.StatusPending
	suite.cmdDueDate = time.Now().Add(24 * time.Hour)
}

// TestHandle tests the Handle method of the addcmd.Handler.
func (suite *HandlerTestSuite) TestHandle() {
	// Create the command using the properties stored in the suite
	cmd := addcmd.NewCommand(suite.cmdTitle, suite.cmdDesc, suite.cmdStatus, suite.cmdDueDate)

	// Set up expected behavior for the mock repository
	suite.mockRepo.On("Save", mock.AnythingOfType("*taskmodel.Task")).Return(nil)

	// Execute the Handle method
	result, err := suite.handler.Handle(cmd)

	// Assertions
	suite.NoError(err)
	suite.NotNil(result)
	suite.Equal(suite.cmdTitle, result.Title())
	suite.Equal(suite.cmdDesc, result.Description())
	suite.Equal(suite.cmdDueDate, result.DueDate())
	suite.Equal(suite.cmdStatus, result.Status())

	// Verify that the Save method was called on the repository with the expected task
	suite.mockRepo.AssertCalled(suite.T(), "Save", mock.AnythingOfType("*taskmodel.Task"))
	suite.mockRepo.AssertExpectations(suite.T())
}

// TestHandle_ErrorCreatingTask tests the Handle method when creating a task fails.
func (suite *HandlerTestSuite) TestHandle_ErrorCreatingTask() {
	// Create a command with properties
	cmd := addcmd.NewCommand("", suite.cmdDesc, suite.cmdStatus, suite.cmdDueDate)

	// Execute the Handle method
	result, err := suite.handler.Handle(cmd)

	// Assertions
	suite.Error(err)
	suite.Nil(result)
}

// TestHandle_ErrorSavingTask tests the Handle method when saving a task fails.
func (suite *HandlerTestSuite) TestHandle_ErrorSavingTask() {
	// Create the command using the properties stored in the suite
	cmd := addcmd.NewCommand(suite.cmdTitle, suite.cmdDesc, suite.cmdStatus, suite.cmdDueDate)

	suite.mockRepo.On("Save", mock.AnythingOfType("*taskmodel.Task")).Return(errors.New("failed to save task"))
	// Execute the Handle method
	result, err := suite.handler.Handle(cmd)

	// Assertions
	suite.Error(err)
	suite.Nil(result)
}

// Run the test suite
func TestHandlerTestSuite(t *testing.T) {
	suite.Run(t, new(HandlerTestSuite))
}
