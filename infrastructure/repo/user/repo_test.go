package userrepo_test

import (
	"context"
	"testing"

	usermodel "github.com/beka-birhanu/task_manager_final/domain/models/user"
	userrepo "github.com/beka-birhanu/task_manager_final/infrastructure/repo/user"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type UserRepositorySuite struct {
	suite.Suite
	client     *mongo.Client
	collection *mongo.Collection
	repo       *userrepo.Repo
	user       *usermodel.User
}

func (suite *UserRepositorySuite) SetupSuite() {
	// Connect to the MongoDB instance
	clientOptions := options.Client().ApplyURI("mongodb://localhost:27017")
	client, err := mongo.Connect(context.Background(), clientOptions)
	if err != nil {
		suite.T().Fatal(err)
	}

	// Set up a test database
	suite.client = client
	suite.collection = client.Database("test_db").Collection("users")
	suite.repo = userrepo.NewRepo(client, "test_db", "users")
}

func (suite *UserRepositorySuite) TearDownSuite() {
	// Disconnect the client
	err := suite.client.Disconnect(context.Background())
	if err != nil {
		suite.T().Fatal(err)
	}
}

func (suite *UserRepositorySuite) TearDownTest() {
	// Clear the collection after each test
	_, err := suite.collection.DeleteMany(context.Background(), bson.D{{}})
	if err != nil {
		suite.T().Fatal(err)
	}
}

func (suite *UserRepositorySuite) SetupTest() {
	// Create a user for reuse in tests
	user, err := usermodel.New(usermodel.Config{
		Username:      "testuser",
		PlainPassword: "testpassword",
		IsAdmin:       false,
	})
	if err != nil {
		suite.T().Fatal(err)
	}
	suite.user = user
}

func (suite *UserRepositorySuite) TestSave_Success() {
	err := suite.repo.Save(suite.user)
	assert.NoError(suite.T(), err)

	// Retrieve the user
	retrievedUser, err := suite.repo.ByUsername(suite.user.Username())
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), suite.user.Username(), retrievedUser.Username())
	assert.Equal(suite.T(), suite.user.PasswordHash(), retrievedUser.PasswordHash())
}

// Test Suite Execution
func TestUserRepositorySuite(t *testing.T) {
	suite.Run(t, new(UserRepositorySuite))
}
