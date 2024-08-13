package promotcmd_test

import (
	"errors"
	"testing"

	irepo_mock "github.com/beka-birhanu/task_manager_final/app/common/i_repo/mocks"
	"github.com/beka-birhanu/task_manager_final/app/user/admin_status/command"
	ihash_mocks "github.com/beka-birhanu/task_manager_final/domain/i_hash/mocks"
	usermodel "github.com/beka-birhanu/task_manager_final/domain/models/user"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

// PromoteCommandHandlerTestSuite defines the test suite for the PromoteCommand handler.
type PromoteCommandHandlerTestSuite struct {
	suite.Suite
	mockUserRepo   *irepo_mock.User
	handler        *promotcmd.Handler
	passwordHasher *ihash_mocks.Service
	admin          *usermodel.User
	user           *usermodel.User
}

// SetupTest sets up the test environment.
func (suite *PromoteCommandHandlerTestSuite) SetupTest() {
	suite.mockUserRepo = new(irepo_mock.User)
	suite.handler = promotcmd.New(suite.mockUserRepo)
	suite.passwordHasher = new(ihash_mocks.Service)

	suite.passwordHasher.On("Hash", mock.AnythingOfType("string")).Return("", nil)

	// Create the user and admin objects
	var err error
	suite.user, err = usermodel.New(usermodel.Config{
		Username:       "user1",
		PlainPassword:  "&&^_str0ngp@ssw0rd!@d$",
		IsAdmin:        false,
		PasswordHasher: suite.passwordHasher,
	})
	suite.Require().NoError(err)

	suite.admin, err = usermodel.New(usermodel.Config{
		Username:       "admin1",
		PlainPassword:  "&&^_str0ngp@ssw0rd!@d$",
		IsAdmin:        true,
		PasswordHasher: suite.passwordHasher,
	})
	suite.Require().NoError(err)

}

// TestHandle_Success tests the successful promotion of a user.
func (suite *PromoteCommandHandlerTestSuite) TestHandle_Success() {

	cmd := &promotcmd.Command{
		Username:   suite.user.Username(),
		PromoterID: suite.admin.ID(),
	}

	// Setup the user repository to return the correct users.
	suite.mockUserRepo.On("ByUsername", suite.user.Username()).Return(suite.user, nil)
	suite.mockUserRepo.On("ById", suite.admin.ID()).Return(suite.admin, nil)
	suite.mockUserRepo.On("Save", suite.user).Return(nil)

	// Execute the Handle method with the command.
	result, err := suite.handler.Handle(cmd)

	// Verify the results.
	suite.Assert().True(result)
	suite.Assert().NoError(err)

	// Verify that the mocks were called as expected.
	suite.mockUserRepo.AssertExpectations(suite.T())
}

// TestHandle_UserNotFound tests the scenario where the user to be promoted is not found.
func (suite *PromoteCommandHandlerTestSuite) TestHandle_UserNotFound() {
	cmd := &promotcmd.Command{
		Username:   suite.user.Username(),
		PromoterID: suite.admin.ID(),
	}

	// Setup the user repository to return an error when fetching the user.
	suite.mockUserRepo.On("ByUsername", suite.user.Username()).Return(nil, errors.New("user not found"))

	// Execute the Handle method with the command.
	result, err := suite.handler.Handle(cmd)

	// Verify the results.
	suite.Assert().False(result)
	suite.Assert().EqualError(err, "user not found")

	// Verify that the mocks were called as expected.
	suite.mockUserRepo.AssertExpectations(suite.T())
}

// TestHandle_AdminNotFound tests the scenario where the promoter (admin) is not found.
func (suite *PromoteCommandHandlerTestSuite) TestHandle_AdminNotFound() {
	cmd := &promotcmd.Command{
		Username:   suite.user.Username(),
		PromoterID: suite.admin.ID(),
	}
	// Setup the user repository to return the user and an error for the admin.
	suite.mockUserRepo.On("ByUsername", suite.user.Username()).Return(suite.user, nil)
	suite.mockUserRepo.On("ById", suite.admin.ID()).Return(nil, errors.New("admin not found"))

	// Execute the Handle method with the command.
	result, err := suite.handler.Handle(cmd)

	// Verify the results.
	suite.Assert().False(result)
	suite.Assert().EqualError(err, "admin not found")

	// Verify that the mocks were called as expected.
	suite.mockUserRepo.AssertExpectations(suite.T())
}

// TestHandle_SaveError tests the scenario where there is an error saving the updated user.
func (suite *PromoteCommandHandlerTestSuite) TestHandle_SaveError() {
	cmd := &promotcmd.Command{
		Username:   suite.user.Username(),
		PromoterID: suite.admin.ID(),
	}

	// Setup the user repository to return the correct users and an error when saving.
	suite.mockUserRepo.On("ByUsername", suite.user.Username()).Return(suite.user, nil)
	suite.mockUserRepo.On("ById", suite.admin.ID()).Return(suite.admin, nil)
	suite.mockUserRepo.On("Save", suite.user).Return(errors.New("save error"))

	// Execute the Handle method with the command.
	result, err := suite.handler.Handle(cmd)

	// Verify the results.
	suite.Assert().False(result)
	suite.Assert().EqualError(err, "save error")

	// Verify that the mocks were called as expected.
	suite.mockUserRepo.AssertExpectations(suite.T())
}

// Run the test suite
func TestPromoteCommandHandlerTestSuite(t *testing.T) {
	suite.Run(t, new(PromoteCommandHandlerTestSuite))
}
