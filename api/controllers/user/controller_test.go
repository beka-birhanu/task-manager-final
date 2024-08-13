package usercontroller_test

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	usercontroller "github.com/beka-birhanu/task_manager_final/api/controllers/user"
	icmd_mock "github.com/beka-birhanu/task_manager_final/app/common/cqrs/command/mocks"
	promotcmd "github.com/beka-birhanu/task_manager_final/app/user/admin_status/command"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

type UserControllerTestSuite struct {
	suite.Suite
	controller        *usercontroller.Controller
	mockPromotHandler *icmd_mock.IHandler[*promotcmd.Command, bool]
	router            *gin.Engine
}

func (suite *UserControllerTestSuite) SetupTest() {
	suite.mockPromotHandler = new(icmd_mock.IHandler[*promotcmd.Command, bool])

	suite.controller = usercontroller.New(usercontroller.Config{
		PromotHandler: suite.mockPromotHandler,
	})

	suite.router = gin.Default()
	api := suite.router.Group("/api")
	suite.controller.RegisterPrivileged(api)
}

func (suite *UserControllerTestSuite) TestPromot_UsernameMissing() {
	req, _ := http.NewRequest(http.MethodPatch, "/api/users//promot", nil)
	ctx := context.WithValue(req.Context(), "userClaims", jwt.MapClaims{}) //nolint
	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, req.WithContext(ctx))

	suite.Equal(http.StatusBadRequest, w.Code)
}

func (suite *UserControllerTestSuite) TestPromot_ClaimsNotFound() {
	req, _ := http.NewRequest(http.MethodPatch, "/api/users/testuser/promot", nil)
	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, req)

	suite.Equal(http.StatusUnauthorized, w.Code)
}

func (suite *UserControllerTestSuite) TestPromot_InvalidClaims() {
	req, _ := http.NewRequest(http.MethodPatch, "/api/users/testuser/promot", nil)
	// Set invalid claims
	ctx := context.WithValue(req.Context(), "userClaims", "invalid_claims") //nolint
	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, req.WithContext(ctx))

	suite.Equal(http.StatusUnauthorized, w.Code)
}

func (suite *UserControllerTestSuite) TestPromot_InvalidUserIDClaim() {
	req, _ := http.NewRequest(http.MethodPatch, "/api/users/testuser/promot", nil)
	ctx := context.WithValue(req.Context(), "userClaims", //nolint
		jwt.MapClaims{
			"user_id": "invalid-uuid",
		})

	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, req.WithContext(ctx))

	suite.Equal(http.StatusUnauthorized, w.Code)
}

func (suite *UserControllerTestSuite) TestPromot_PromotionFailure() {
	promoterId := uuid.New()
	suite.mockPromotHandler.On("Handle", mock.AnythingOfType("*promotcmd.Command")).Return(false, errors.New("promotion failed"))

	// Create a mock JWT token with user_id claim
	req, _ := http.NewRequest(http.MethodPatch, "/api/users/testuser/promot", nil)
	ctx := context.WithValue(req.Context(), "userClaims", //nolint
		jwt.MapClaims{
			"user_id": promoterId.String(),
		})

	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, req.WithContext(ctx))

	suite.Equal(http.StatusUnauthorized, w.Code)
}

func TestUserControllerTestSuite(t *testing.T) {
	suite.Run(t, new(UserControllerTestSuite))
}
