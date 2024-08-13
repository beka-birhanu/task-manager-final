package registercmd_test

import (
	"errors"
	"testing"

	ijwt_mock "github.com/beka-birhanu/task_manager_final/app/common/i_jwt/mock"
	irepo_mock "github.com/beka-birhanu/task_manager_final/app/common/i_repo/mocks"
	registercmd "github.com/beka-birhanu/task_manager_final/app/user/auth/command"
	ihash_mocks "github.com/beka-birhanu/task_manager_final/domain/i_hash/mocks"
	usermodel "github.com/beka-birhanu/task_manager_final/domain/models/user"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

// RegisterCommandHandlerTestSuite defines the test suite for the RegisterCommand handler.
type RegisterCommandHandlerTestSuite struct {
	suite.Suite
	mockUserRepo *irepo_mock.User
	mockJwtSvc   *ijwt_mock.MockService
	mockHashSvc  *ihash_mocks.Service
	handler      *registercmd.Handler
	adminUser    *usermodel.User
	normalUser   *usermodel.User
	cmd          *registercmd.Command
}

// SetupTest sets up the test environment.
func (suite *RegisterCommandHandlerTestSuite) SetupTest() {
	suite.mockUserRepo = new(irepo_mock.User)
	suite.mockJwtSvc = new(ijwt_mock.MockService)
	suite.mockHashSvc = new(ihash_mocks.Service)

	suite.handler = registercmd.NewHandler(registercmd.Config{
		UserRepo: suite.mockUserRepo,
		JwtSvc:   suite.mockJwtSvc,
		HashSvc:  suite.mockHashSvc,
	})

	suite.mockHashSvc.On("Hash", mock.AnythingOfType("string")).Return("hashed_password", nil)

	// Initialize admin user object
	var err error
	suite.adminUser, err = usermodel.New(usermodel.Config{
		Username:       "adminuser",
		PlainPassword:  "&&^_str0ngp@ssw0rd!@d$",
		IsAdmin:        true,
		PasswordHasher: suite.mockHashSvc,
	})
	suite.Require().NoError(err)

	// Initialize normal user object
	suite.normalUser, err = usermodel.New(usermodel.Config{
		Username:       "normaluser",
		PlainPassword:  "&&^_str0ngp@ssw0rd!@d$",
		IsAdmin:        false,
		PasswordHasher: suite.mockHashSvc,
	})
	suite.Require().NoError(err)

	// Initialize the command object
	suite.cmd = &registercmd.Command{
		Username: "normaluser",
		Password: "&&^_str0ngp@ssw0rd!@d$",
	}
}

// TestHandle_Success tests the successful registration of a user.
func (suite *RegisterCommandHandlerTestSuite) TestHandle_Success() {
	// Mock expected behavior
	suite.mockHashSvc.On("Hash", suite.cmd.Password).Return("hashed_password", nil)
	suite.mockUserRepo.On("Count").Return(int64(0), nil)
	suite.mockUserRepo.On("Save", mock.AnythingOfType("*usermodel.User")).Return(nil)
	suite.mockJwtSvc.On("Generate", mock.AnythingOfType("*usermodel.User")).Return("jwt_token", nil)

	// Execute the Handle method with the command.
	result, err := suite.handler.Handle(suite.cmd)

	// Verify the results.
	suite.Assert().NotNil(result)
	suite.Assert().Equal("jwt_token", result.Token)
	suite.Assert().NoError(err)

	// Verify that the mocks were called as expected.
	suite.mockUserRepo.AssertExpectations(suite.T())
	suite.mockJwtSvc.AssertExpectations(suite.T())
	suite.mockHashSvc.AssertExpectations(suite.T())
}

// TestHandle_UserRepoCountError tests the scenario where counting users in the repository fails.
func (suite *RegisterCommandHandlerTestSuite) TestHandle_UserRepoCountError() {
	suite.mockUserRepo.On("Count").Return(int64(0), errors.New("repo error"))

	// Execute the Handle method with the command.
	result, err := suite.handler.Handle(suite.cmd)

	// Verify the results.
	suite.Assert().Nil(result)
	suite.Assert().EqualError(err, "repo error")

	// Verify that the mocks were called as expected.
	suite.mockUserRepo.AssertExpectations(suite.T())
	suite.mockJwtSvc.AssertExpectations(suite.T())
	suite.mockHashSvc.AssertExpectations(suite.T())
}

// TestHandle_UserCreationError tests the scenario where user creation fails.
func (suite *RegisterCommandHandlerTestSuite) TestHandle_UserCreationError() {
	suite.mockHashSvc.On("Hash", suite.cmd.Password).Return("", errors.New("hash error"))
	suite.mockUserRepo.On("Count").Return(int64(0), nil)

	// Execute the Handle method with the command.
	cmd := suite.cmd
	cmd.Password = ""
	result, err := suite.handler.Handle(cmd)

	// Verify the results.
	suite.Assert().Nil(result)
	suite.Assert().EqualError(err, "Validation: password is too weak.")

	// Verify that the mocks were called as expected.
	suite.mockUserRepo.AssertExpectations(suite.T())
	suite.mockJwtSvc.AssertExpectations(suite.T())
	suite.mockHashSvc.AssertExpectations(suite.T())
}

// TestHandle_JwtGenerationError tests the scenario where JWT generation fails.
func (suite *RegisterCommandHandlerTestSuite) TestHandle_JwtGenerationError() {
	suite.mockHashSvc.On("Hash", suite.cmd.Password).Return("hashed_password", nil)
	suite.mockUserRepo.On("Count").Return(int64(0), nil)
	suite.mockUserRepo.On("Save", mock.AnythingOfType("*usermodel.User")).Return(nil)
	suite.mockJwtSvc.On("Generate", mock.AnythingOfType("*usermodel.User")).Return("", errors.New("jwt error"))

	// Execute the Handle method with the command.
	result, err := suite.handler.Handle(suite.cmd)

	// Verify the results.
	suite.Assert().Nil(result)
	suite.Assert().EqualError(err, "jwt error")

	// Verify that the mocks were called as expected.
	suite.mockUserRepo.AssertExpectations(suite.T())
	suite.mockJwtSvc.AssertExpectations(suite.T())
	suite.mockHashSvc.AssertExpectations(suite.T())
}

// TestHandle_NormalUser_Success tests the successful registration of a normal user.
func (suite *RegisterCommandHandlerTestSuite) TestHandle_NormalUser_Success() {
	// Mock expected behavior for a normal user
	suite.mockHashSvc.On("Hash", suite.cmd.Password).Return("hashed_password", nil)
	suite.mockUserRepo.On("Count").Return(int64(0), nil)
	suite.mockUserRepo.On("Save", mock.AnythingOfType("*usermodel.User")).Return(nil)
	suite.mockJwtSvc.On("Generate", mock.AnythingOfType("*usermodel.User")).Return("jwt_token", nil)

	// Execute the Handle method with the command.
	result, err := suite.handler.Handle(suite.cmd)

	// Verify the results.
	suite.Assert().NotNil(result)
	suite.Assert().Equal("jwt_token", result.Token)
	suite.Assert().NoError(err)

	// Verify that the mocks were called as expected.
	suite.mockUserRepo.AssertExpectations(suite.T())
	suite.mockJwtSvc.AssertExpectations(suite.T())
	suite.mockHashSvc.AssertExpectations(suite.T())
}

// Run the test suite
func TestRegisterCommandHandlerTestSuite(t *testing.T) {
	suite.Run(t, new(RegisterCommandHandlerTestSuite))
}
