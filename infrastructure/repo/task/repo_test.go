package taskrepo_test

import (
	"context"
	"testing"
	"time"

	taskmodel "github.com/beka-birhanu/task_manager_final/domain/models/task"
	taskrepo "github.com/beka-birhanu/task_manager_final/infrastructure/repo/task"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type TaskRepositorySuite struct {
	suite.Suite
	client     *mongo.Client
	repo       *taskrepo.Repo
	collection *mongo.Collection
	task       *taskmodel.Task
}

func (suite *TaskRepositorySuite) SetupSuite() {
	clientOptions := options.Client().ApplyURI("mongodb://localhost:27017")
	client, err := mongo.Connect(context.Background(), clientOptions)
	if err != nil {
		suite.T().Fatal(err)
	}

	suite.client = client
	suite.collection = client.Database("test_db").Collection("tasks")
	suite.repo = taskrepo.New(client, "test_db", "tasks")
}

func (suite *TaskRepositorySuite) TearDownSuite() {
	err := suite.client.Disconnect(context.Background())
	if err != nil {
		suite.T().Fatal(err)
	}
}

func (suite *TaskRepositorySuite) SetupTest() {
	// Clear the collection before each test
	err := suite.collection.Drop(context.Background())
	if err != nil {
		suite.T().Fatal(err)
	}

	// Create a task for testing
	suite.task, err = taskmodel.New(taskmodel.Config{
		Title:       `Test Task`,
		Description: "Test Description",
		DueDate:     time.Now(),
		Status:      "pending",
	})
	if err != nil {
		suite.T().Fatal(err)
	}

	// Save the task to the repository
	err = suite.repo.Save(suite.task)
	if err != nil {
		suite.T().Fatal(err)
	}
}

func (suite *TaskRepositorySuite) TearDownTest() {
	// Additional cleanup if needed
	_, err := suite.collection.DeleteMany(context.Background(), bson.D{})
	if err != nil {
		suite.T().Fatal(err)
	}
}

func (suite *TaskRepositorySuite) TestSaveTask() {
	savedTask, err := suite.repo.GetSingle(suite.task.ID())
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), suite.task.Title(), savedTask.Title())
	assert.Equal(suite.T(), suite.task.Description(), savedTask.Description())
	assert.Equal(suite.T(), suite.task.Status(), savedTask.Status())
}

func (suite *TaskRepositorySuite) TestDeleteTask() {
	err := suite.repo.Delete(suite.task.ID())
	assert.NoError(suite.T(), err)

	_, err = suite.repo.GetSingle(suite.task.ID())
	assert.Error(suite.T(), err)
	assert.Equal(suite.T(), "task Not Found", err.Error())
}

func (suite *TaskRepositorySuite) TestGetAllTasks() {
	tasks, err := suite.repo.GetAll()
	assert.NoError(suite.T(), err)
	assert.Len(suite.T(), tasks, 1) // Adjusted to match the number of tasks saved in SetupTest
}

func (suite *TaskRepositorySuite) TestGetSingleTask() {
	foundTask, err := suite.repo.GetSingle(suite.task.ID())
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), suite.task.Title(), foundTask.Title())
	assert.Equal(suite.T(), suite.task.Description(), foundTask.Description())
	assert.Equal(suite.T(), suite.task.Status(), foundTask.Status())
}

func TestTaskRepositorySuite(t *testing.T) {
	suite.Run(t, new(TaskRepositorySuite))
}
