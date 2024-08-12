package usermodel_test

import (
	"testing"

	errdmn "github.com/beka-birhanu/task_manager_final/domain/errors"
	ihash_mock "github.com/beka-birhanu/task_manager_final/domain/i_hash/mocks"
	usermodel "github.com/beka-birhanu/task_manager_final/domain/models/user"
	"github.com/google/uuid"
	"github.com/stretchr/testify/suite"
)

type UserModelSuite struct {
	suite.Suite
	validConfig     usermodel.Config
	mockHasher      *ihash_mock.Service
	weakPassword    string
	strongPassword  string
	validUsername   string
	invalidUsername string
	user            *usermodel.User
}

func (suite *UserModelSuite) SetupTest() {
	suite.mockHasher = new(ihash_mock.Service)
	suite.weakPassword = "1234"
	suite.strongPassword = "&&^_Str0ngP@ssw0rd!@D$"
	suite.validUsername = "valid_user"
	suite.invalidUsername = "invalid user!" // invalid due to space and special character

	suite.mockHasher.On("Hash", suite.strongPassword).Return("hashedStrongPassword", nil)

	suite.validConfig = usermodel.Config{
		Username:       suite.validUsername,
		PlainPassword:  suite.strongPassword,
		IsAdmin:        false,
		PasswordHasher: suite.mockHasher,
	}

	var err error
	suite.user, err = usermodel.New(suite.validConfig)
	suite.NoError(err)
	suite.NotNil(suite.user)
}

func (suite *UserModelSuite) TestNewUser() {
	suite.Run("should create a new user with valid config", func() {
		user, err := usermodel.New(suite.validConfig)
		suite.NoError(err)
		suite.NotNil(user)
		suite.Equal(suite.validConfig.Username, user.Username())
		suite.Equal("hashedStrongPassword", user.PasswordHash())
		suite.False(user.IsAdmin())
		suite.NotEqual(uuid.Nil, user.ID())
	})

	suite.Run("should return error if username is invalid", func() {
		invalidConfig := suite.validConfig
		invalidConfig.Username = suite.invalidUsername
		user, err := usermodel.New(invalidConfig)
		suite.Nil(user)
		suite.Equal(errdmn.UsernameInvalidFormat, err)
	})

	suite.Run("should return error if password is weak", func() {
		invalidConfig := suite.validConfig
		invalidConfig.PlainPassword = suite.weakPassword
		user, err := usermodel.New(invalidConfig)
		suite.Nil(user)
		suite.Equal(errdmn.WeakPassword, err)
	})
}

func (suite *UserModelSuite) TestUser_UpdateUsername() {
	suite.Run("should update the username if valid", func() {
		newUsername := "new_valid_user"
		err := suite.user.UpdateUsername(newUsername)
		suite.NoError(err)
		suite.Equal(newUsername, suite.user.Username())
	})

	suite.Run("should return error if new username is invalid", func() {
		err := suite.user.UpdateUsername(suite.invalidUsername)
		suite.Equal(errdmn.UsernameInvalidFormat, err)
	})
}

func (suite *UserModelSuite) TestUser_UpdatePassword() {
	suite.Run("should update the password if valid", func() {
		newPassword := "&&)(&&^%)N3wStr0ngP@ss"
		suite.mockHasher.On("Hash", newPassword).Return("hashedNewPassword", nil)
		err := suite.user.UpdatePassword(newPassword, suite.mockHasher)
		suite.NoError(err)
		suite.Equal("hashedNewPassword", suite.user.PasswordHash())
	})

	suite.Run("should return error if new password is weak", func() {
		err := suite.user.UpdatePassword(suite.weakPassword, suite.mockHasher)
		suite.Equal(errdmn.WeakPassword, err)
	})
}

func (suite *UserModelSuite) TestUser_UpdateAdminStatus() {
	suite.Run("should update the admin status", func() {
		suite.user.UpdateAdminStatus(true)
		suite.True(suite.user.IsAdmin())

		suite.user.UpdateAdminStatus(false)
		suite.False(suite.user.IsAdmin())
	})
}

func (suite *UserModelSuite) TestFromBSON() {
	bsonUser := &usermodel.UserBSON{
		ID:           uuid.New(),
		Username:     "bson_user",
		PasswordHash: "hashedPassword",
		IsAdmin:      false,
	}

	suite.Run("should convert BSON to user", func() {
		user := usermodel.FromBSON(bsonUser)
		suite.Equal(bsonUser.ID, user.ID())
		suite.Equal(bsonUser.Username, user.Username())
		suite.Equal(bsonUser.PasswordHash, user.PasswordHash())
		suite.Equal(bsonUser.IsAdmin, user.IsAdmin())
	})
}

func TestUserModelSuite(t *testing.T) {
	suite.Run(t, new(UserModelSuite))
}
