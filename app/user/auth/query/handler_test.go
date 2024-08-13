package loginqry_test

import (
	"errors"
	"log"
	"testing"

	ijwt_mock "github.com/beka-birhanu/task_manager_final/app/common/i_jwt/mock"
	irepo_mock "github.com/beka-birhanu/task_manager_final/app/common/i_repo/mocks"
	loginqry "github.com/beka-birhanu/task_manager_final/app/user/auth/query"
	ihash_mocks "github.com/beka-birhanu/task_manager_final/domain/i_hash/mocks"
	usermodel "github.com/beka-birhanu/task_manager_final/domain/models/user"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

// LoginQueryHandlerTestSuite defines the test suite for the LoginQuery handler.
type LoginQueryHandlerTestSuite struct {
	suite.Suite
	mockUserRepo *irepo_mock.User
	mockJwtSvc   *ijwt_mock.MockService
	mockHashSvc  *ihash_mocks.Service
	handler      *loginqry.Handler
	existingUser *usermodel.User
	query        *loginqry.Query
}

// SetupTest sets up the test environment.
func (suite *LoginQueryHandlerTestSuite) SetupTest() {
	suite.mockUserRepo = new(irepo_mock.User)
	suite.mockJwtSvc = new(ijwt_mock.MockService)
	suite.mockHashSvc = new(ihash_mocks.Service)

	suite.handler = loginqry.NewHandler(loginqry.Config{
		UserRepo: suite.mockUserRepo,
		JwtSvc:   suite.mockJwtSvc,
		HashSvc:  suite.mockHashSvc,
	})

	suite.mockHashSvc.On("Hash", mock.AnythingOfType("string")).Return("hashed_password", nil)
	// Initialize existing user object
	var err error
	suite.existingUser, err = usermodel.New(usermodel.Config{
		Username:       "existinguser",
		PlainPassword:  "&&^_str0ngp@ssw0rd!@d$",
		IsAdmin:        false,
		PasswordHasher: suite.mockHashSvc,
	})
	suite.Require().NoError(err)

	suite.query = &loginqry.Query{
		Username: "existinguser",
		Password: "&&^_str0ngp@ssw0rd!@d$",
	}
}

// TestHandle_Success tests the successful login and JWT generation.
func (suite *LoginQueryHandlerTestSuite) TestHandle_Success() {
	// Mock expected behavior
	suite.mockUserRepo.On("ByUsername", suite.query.Username).Return(suite.existingUser, nil)
	suite.mockHashSvc.On("Match", suite.existingUser.PasswordHash(), suite.query.Password).Return(true, nil)
	suite.mockJwtSvc.On("Generate", suite.existingUser).Return("jwt_token", nil)

	// Execute the Handle method with the query.
	result, err := suite.handler.Handle(suite.query)

	// Verify the results.
	suite.Assert().NotNil(result)
	suite.Assert().Equal("jwt_token", result.Token)
	suite.Assert().NoError(err)

	// Verify that the mocks were called as expected.
	suite.mockUserRepo.AssertExpectations(suite.T())
	suite.mockJwtSvc.AssertExpectations(suite.T())
	suite.mockHashSvc.AssertExpectations(suite.T())
}

// TestHandle_UserRepoError tests the scenario where retrieving the user from the repository fails.
func (suite *LoginQueryHandlerTestSuite) TestHandle_UserRepoError() {
	suite.mockUserRepo.On("ByUsername", suite.query.Username).Return(nil, errors.New("repo error"))

	// Execute the Handle method with the query.
	result, err := suite.handler.Handle(suite.query)

	// Verify the results.
	suite.Assert().Nil(result)
	suite.Assert().EqualError(err, "Unauthorized: repo error")

	// Verify that the mocks were called as expected.
	suite.mockUserRepo.AssertExpectations(suite.T())
	suite.mockJwtSvc.AssertExpectations(suite.T())
	suite.mockHashSvc.AssertExpectations(suite.T())
}

// TestHandle_InvalidPassword tests the scenario where the password does not match.
func (suite *LoginQueryHandlerTestSuite) TestHandle_InvalidPassword() {
	log.Println(suite.existingUser.PasswordHash(), suite.query.Password, "test")
	suite.mockUserRepo.On("ByUsername", suite.query.Username).Return(suite.existingUser, nil)
	suite.mockHashSvc.On("Match", suite.existingUser.PasswordHash(), suite.query.Password).Return(false, nil)

	// Execute the Handle method with the query.
	result, err := suite.handler.Handle(suite.query)

	// Verify the results.
	suite.Assert().Nil(result)
	suite.Assert().EqualError(err, "Unauthorized: incorrect password")

	// Verify that the mocks were called as expected.
	suite.mockUserRepo.AssertExpectations(suite.T())
	suite.mockJwtSvc.AssertExpectations(suite.T())
	suite.mockHashSvc.AssertExpectations(suite.T())
}

// TestHandle_PasswordHashError tests the scenario where password hash comparison fails.
func (suite *LoginQueryHandlerTestSuite) TestHandle_PasswordHashError() {
	suite.mockUserRepo.On("ByUsername", suite.query.Username).Return(suite.existingUser, nil)
	suite.mockHashSvc.On("Match", suite.existingUser.PasswordHash(), suite.query.Password).Return(false, errors.New("hash error"))

	// Execute the Handle method with the query.
	result, err := suite.handler.Handle(suite.query)

	// Verify the results.
	suite.Assert().Nil(result)
	suite.Assert().EqualError(err, "ServerError: failed to validate user password, hash error")

	// Verify that the mocks were called as expected.
	suite.mockUserRepo.AssertExpectations(suite.T())
	suite.mockJwtSvc.AssertExpectations(suite.T())
	suite.mockHashSvc.AssertExpectations(suite.T())
}

// TestHandle_JwtGenerationError tests the scenario where JWT generation fails.
func (suite *LoginQueryHandlerTestSuite) TestHandle_JwtGenerationError() {
	suite.mockUserRepo.On("ByUsername", suite.query.Username).Return(suite.existingUser, nil)
	suite.mockHashSvc.On("Match", suite.existingUser.PasswordHash(), suite.query.Password).Return(true, nil)
	suite.mockJwtSvc.On("Generate", suite.existingUser).Return("", errors.New("jwt error"))

	// Execute the Handle method with the query.
	result, err := suite.handler.Handle(suite.query)

	// Verify the results.
	suite.Assert().Nil(result)
	suite.Assert().EqualError(err, "ServerError: failed to generate JWT for user, jwt error")

	// Verify that the mocks were called as expected.
	suite.mockUserRepo.AssertExpectations(suite.T())
	suite.mockJwtSvc.AssertExpectations(suite.T())
	suite.mockHashSvc.AssertExpectations(suite.T())
}

// Run the test suite
func TestLoginQueryHandlerTestSuite(t *testing.T) {
	suite.Run(t, new(LoginQueryHandlerTestSuite))
}
