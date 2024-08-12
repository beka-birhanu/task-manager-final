package taskmodel_test

import (
	"testing"
	"time"

	errdmn "github.com/beka-birhanu/task_manager_final/domain/errors"
	taskmodel "github.com/beka-birhanu/task_manager_final/domain/models/task"
	"github.com/google/uuid"
	"github.com/stretchr/testify/suite"
)

type TaskModelSuite struct {
	suite.Suite
	validConfig taskmodel.Config
	task        *taskmodel.Task
}

func (suite *TaskModelSuite) SetupTest() {
	suite.validConfig = taskmodel.Config{
		Title:       "Test Task",
		Description: "This is a test task.",
		DueDate:     time.Now().Add(24 * time.Hour),
		Status:      taskmodel.StatusPending,
	}
	var err error
	suite.task, err = taskmodel.New(suite.validConfig)
	suite.NoError(err)
}

func (suite *TaskModelSuite) TestNewTask() {
	suite.Run("should create a new task with valid config", func() {
		task, err := taskmodel.New(suite.validConfig)
		suite.NoError(err)
		suite.NotNil(task)
		suite.Equal(suite.validConfig.Title, task.Title())
		suite.Equal(suite.validConfig.Description, task.Description())
		suite.Equal(suite.validConfig.DueDate, task.DueDate())
		suite.Equal(suite.validConfig.Status, task.Status())
		suite.NotEqual(uuid.Nil, task.ID())
	})

	suite.Run("should return error if title is empty", func() {
		invalidConfig := suite.validConfig
		invalidConfig.Title = ""
		task, err := taskmodel.New(invalidConfig)
		suite.Nil(task)
		suite.Equal(errdmn.TitleEmpty, err)
	})

	suite.Run("should return error if description is empty", func() {
		invalidConfig := suite.validConfig
		invalidConfig.Description = ""
		task, err := taskmodel.New(invalidConfig)
		suite.Nil(task)
		suite.Equal(errdmn.DescriptionEmpty, err)
	})

	suite.Run("should return error if due date is zero", func() {
		invalidConfig := suite.validConfig
		invalidConfig.DueDate = time.Time{}
		task, err := taskmodel.New(invalidConfig)
		suite.Nil(task)
		suite.Equal(errdmn.DueDateZero, err)
	})

	suite.Run("should return error if status is invalid", func() {
		invalidConfig := suite.validConfig
		invalidConfig.Status = "invalid-status"
		task, err := taskmodel.New(invalidConfig)
		suite.Nil(task)
		suite.Equal(errdmn.InvalidStatus, err)
	})
}

func (suite *TaskModelSuite) TestTask_Update() {
	suite.Run("should update task fields with valid config", func() {
		updateConfig := taskmodel.Config{
			Title:       "Updated Task",
			Description: "This is an updated test task.",
			DueDate:     time.Now().Add(48 * time.Hour),
			Status:      taskmodel.StatusInProgress,
		}
		err := suite.task.Update(updateConfig)
		suite.NoError(err)
		suite.Equal(updateConfig.Title, suite.task.Title())
		suite.Equal(updateConfig.Description, suite.task.Description())
		suite.Equal(updateConfig.DueDate, suite.task.DueDate())
		suite.Equal(updateConfig.Status, suite.task.Status())
	})

	suite.Run("should return error when updating with invalid config", func() {
		invalidConfig := taskmodel.Config{
			Title:       "",
			Description: "Invalid update",
			DueDate:     time.Now().Add(48 * time.Hour),
			Status:      taskmodel.StatusInProgress,
		}
		err := suite.task.Update(invalidConfig)
		suite.Equal(errdmn.TitleEmpty, err)
	})
}

func (suite *TaskModelSuite) TestTask_ToBSON() {
	suite.Run("should convert task to BSON", func() {
		bson := suite.task.ToBSON()
		suite.Equal(suite.task.ID(), bson.ID)
		suite.Equal(suite.task.Title(), bson.Title)
		suite.Equal(suite.task.Description(), bson.Description)
		suite.Equal(suite.task.DueDate(), bson.DueDate)
		suite.Equal(suite.task.Status(), bson.Status)
		suite.NotZero(bson.UpdatedAt)
	})
}

func (suite *TaskModelSuite) TestFromBSON() {
	taskBSON := &taskmodel.TaskBSON{
		ID:          uuid.New(),
		Title:       "BSON Task",
		Description: "This is a BSON task.",
		DueDate:     time.Now().Add(24 * time.Hour),
		Status:      taskmodel.StatusDone,
		UpdatedAt:   time.Now(),
	}

	suite.Run("should convert BSON to task", func() {
		task := taskmodel.FromBSON(taskBSON)
		suite.Equal(taskBSON.ID, task.ID())
		suite.Equal(taskBSON.Title, task.Title())
		suite.Equal(taskBSON.Description, task.Description())
		suite.Equal(taskBSON.DueDate, task.DueDate())
		suite.Equal(taskBSON.Status, task.Status())
	})
}

func TestTaskModelSuite(t *testing.T) {
	suite.Run(t, new(TaskModelSuite))
}
