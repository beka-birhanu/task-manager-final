package taskcontroller_test

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	taskcontroller "github.com/beka-birhanu/task_manager_final/api/controllers/task"
	icmd_mock "github.com/beka-birhanu/task_manager_final/app/common/cqrs/command/mocks"
	addcmd "github.com/beka-birhanu/task_manager_final/app/task/command/add"
	updatecmd "github.com/beka-birhanu/task_manager_final/app/task/command/update"
	taskmodel "github.com/beka-birhanu/task_manager_final/domain/models/task"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

type TaskControllerTestSuite struct {
	suite.Suite
	controller        *taskcontroller.Controller
	mockAddHandler    *icmd_mock.IHandler[*addcmd.Command, *taskmodel.Task]
	mockUpdateHandler *icmd_mock.IHandler[*updatecmd.Command, *taskmodel.Task]
	mockDeleteHandler *icmd_mock.IHandler[uuid.UUID, bool]
	mockGetAllHandler *icmd_mock.IHandler[struct{}, []*taskmodel.Task]
	mockGetHandler    *icmd_mock.IHandler[uuid.UUID, *taskmodel.Task]
	router            *gin.Engine
	testTask          *taskmodel.Task
}

func (suite *TaskControllerTestSuite) SetupTest() {
	suite.mockAddHandler = new(icmd_mock.IHandler[*addcmd.Command, *taskmodel.Task])
	suite.mockUpdateHandler = new(icmd_mock.IHandler[*updatecmd.Command, *taskmodel.Task])
	suite.mockDeleteHandler = new(icmd_mock.IHandler[uuid.UUID, bool])
	suite.mockGetAllHandler = new(icmd_mock.IHandler[struct{}, []*taskmodel.Task])
	suite.mockGetHandler = new(icmd_mock.IHandler[uuid.UUID, *taskmodel.Task])

	suite.controller = taskcontroller.New(taskcontroller.Config{
		AddHandler:    suite.mockAddHandler,
		UpdateHandler: suite.mockUpdateHandler,
		DeleteHandler: suite.mockDeleteHandler,
		GetAllHandler: suite.mockGetAllHandler,
		GetHandler:    suite.mockGetHandler,
	})

	suite.router = gin.Default()
	api := suite.router.Group("/api")
	suite.controller.RegisterProtected(api)
	suite.controller.RegisterPrivileged(api)

	suite.testTask, _ = taskmodel.New(
		taskmodel.Config{
			Title:       "Test Task",
			Description: "This is a test task.",
			DueDate:     time.Now(),
			Status:      "pending",
		})
}

func (suite *TaskControllerTestSuite) TestAddTask_Success() {
	suite.mockAddHandler.On("Handle", mock.AnythingOfType("*addcmd.Command")).Return(suite.testTask, nil)

	reqBody := `{
		"title": "Test Task",
		"description": "This is a test task.",
		"status": "pending",
		"dueDate": "2024-08-30T00:00:00Z"
	}`
	req, _ := http.NewRequest(http.MethodPost, "/api/tasks", strings.NewReader(reqBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, req)

	suite.Equal(http.StatusCreated, w.Code)
	suite.mockAddHandler.AssertExpectations(suite.T())
}

func (suite *TaskControllerTestSuite) TestUpdateTask_Success() {
	id := suite.testTask.ID()
	suite.mockUpdateHandler.On("Handle", mock.AnythingOfType("*updatecmd.Command")).Return(suite.testTask, nil)

	reqBody := `{
		"title": "Test Task",
		"description": "This is a test task.",
		"status": "pending",
		"dueDate": "2024-08-30T00:00:00Z"
	}`
	req, _ := http.NewRequest(http.MethodPut, "/api/tasks/"+id.String(), strings.NewReader(reqBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, req)

	suite.Equal(http.StatusOK, w.Code)
	suite.mockUpdateHandler.AssertExpectations(suite.T())
}

func (suite *TaskControllerTestSuite) TestDeleteTask_Success() {
	id := suite.testTask.ID()
	suite.mockDeleteHandler.On("Handle", id).Return(true, nil)

	req, _ := http.NewRequest(http.MethodDelete, "/api/tasks/"+id.String(), nil)
	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, req)

	suite.Equal(http.StatusOK, w.Code)
	suite.mockDeleteHandler.AssertExpectations(suite.T())
}

func (suite *TaskControllerTestSuite) TestGetAllTasks_Success() {
	suite.mockGetAllHandler.On("Handle", mock.Anything).Return([]*taskmodel.Task{suite.testTask}, nil)

	req, _ := http.NewRequest(http.MethodGet, "/api/tasks", nil)
	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, req)

	suite.Equal(http.StatusOK, w.Code)
	suite.mockGetAllHandler.AssertExpectations(suite.T())
}

func (suite *TaskControllerTestSuite) TestGetTask_Success() {
	id := suite.testTask.ID()
	suite.mockGetHandler.On("Handle", id).Return(suite.testTask, nil)

	req, _ := http.NewRequest(http.MethodGet, "/api/tasks/"+id.String(), nil)
	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, req)

	suite.Equal(http.StatusOK, w.Code)
	suite.mockGetHandler.AssertExpectations(suite.T())
}

func TestTaskControllerTestSuite(t *testing.T) {
	suite.Run(t, new(TaskControllerTestSuite))
}
