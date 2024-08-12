package getallqry_test

import (
	"errors"
	"testing"
	"time"

	icmd "github.com/beka-birhanu/task_manager_final/app/common/cqrs/command"
	irepo_mock "github.com/beka-birhanu/task_manager_final/app/common/i_repo/mocks"
	getallqry "github.com/beka-birhanu/task_manager_final/app/task/query/get_all"
	taskmodel "github.com/beka-birhanu/task_manager_final/domain/models/task"
	"github.com/stretchr/testify/suite"
)

// HandlerTestSuite defines the test suite for the getallqry.Handler.
type HandlerTestSuite struct {
	suite.Suite
	mockRepo *irepo_mock.Task
	handler  icmd.IHandler[struct{}, []*taskmodel.Task]
}

// SetupTest sets up the test environment.
func (suite *HandlerTestSuite) SetupTest() {
	// Initialize the mock repository
	suite.mockRepo = new(irepo_mock.Task)

	// Initialize the handler with the mock repository
	suite.handler = getallqry.New(suite.mockRepo)
}

// TestHandle_Success tests the Handle method of the getallqry.Handler for successful retrieval.
func (suite *HandlerTestSuite) TestHandle_Success() {
	// Create mock tasks
	task1, _ := taskmodel.New(taskmodel.Config{
		Title:       "Task 1",
		Description: "First task",
		DueDate:     time.Now().Add(24 * time.Hour),
		Status:      taskmodel.StatusPending,
	})
	task2, _ := taskmodel.New(taskmodel.Config{
		Title:       "Task 2",
		Description: "Second task",
		DueDate:     time.Now().Add(48 * time.Hour),
		Status:      taskmodel.StatusInProgress,
	})

	// Set up expected behavior for the mock repository
	suite.mockRepo.On("GetAll").Return([]*taskmodel.Task{task1, task2}, nil)

	// Execute the Handle method
	tasks, err := suite.handler.Handle(struct{}{})

	// Assertions
	suite.NoError(err)
	suite.Len(tasks, 2)
	suite.Equal(task1.Title(), tasks[0].Title())
	suite.Equal(task2.Title(), tasks[1].Title())

	// Verify that the GetAll method was called on the repository
	suite.mockRepo.AssertCalled(suite.T(), "GetAll")
	suite.mockRepo.AssertExpectations(suite.T())
}

// TestHandle_ErrorRetrievingTasks tests the Handle method when retrieving tasks fails.
func (suite *HandlerTestSuite) TestHandle_ErrorRetrievingTasks() {
	// Set up expected behavior for the mock repository
	suite.mockRepo.On("GetAll").Return(nil, errors.New("failed to retrieve tasks"))

	// Execute the Handle method
	tasks, err := suite.handler.Handle(struct{}{})

	// Assertions
	suite.Error(err)
	suite.Nil(tasks)

	// Verify that the GetAll method was called on the repository
	suite.mockRepo.AssertCalled(suite.T(), "GetAll")
	suite.mockRepo.AssertExpectations(suite.T())
}

// Run the test suite
func TestHandlerTestSuite(t *testing.T) {
	suite.Run(t, new(HandlerTestSuite))
}
