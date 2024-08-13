package authcontroller_test

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	authcontroller "github.com/beka-birhanu/task_manager_final/api/controllers/auth"
	icmd_mock "github.com/beka-birhanu/task_manager_final/app/common/cqrs/command/mocks"
	iquery_mock "github.com/beka-birhanu/task_manager_final/app/common/cqrs/query/mocks"
	registercmd "github.com/beka-birhanu/task_manager_final/app/user/auth/command"
	authresult "github.com/beka-birhanu/task_manager_final/app/user/auth/common"
	loginqry "github.com/beka-birhanu/task_manager_final/app/user/auth/query"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

type AuthControllerTestSuite struct {
	suite.Suite
	controller          *authcontroller.Controller
	mockRegisterHandler *icmd_mock.IHandler[*registercmd.Command, *authresult.Result]
	mockLoginHandler    *iquery_mock.IHandler[*loginqry.Query, *authresult.Result]
	router              *gin.Engine
}

func (suite *AuthControllerTestSuite) SetupTest() {
	suite.mockRegisterHandler = new(icmd_mock.IHandler[*registercmd.Command, *authresult.Result])
	suite.mockLoginHandler = new(iquery_mock.IHandler[*loginqry.Query, *authresult.Result])

	suite.controller = authcontroller.New(authcontroller.Config{
		RegisterHandler: suite.mockRegisterHandler,
		LoginHandler:    suite.mockLoginHandler,
	})

	suite.router = gin.Default()
	api := suite.router.Group("/api")
	suite.controller.RegisterPublic(api)
	suite.controller.RegisterProtected(api)
}

func (suite *AuthControllerTestSuite) TestRegisterUser_Success() {
	result := &authresult.Result{Token: "testtoken"}
	suite.mockRegisterHandler.On("Handle", mock.AnythingOfType("*registercmd.Command")).Return(result, nil)

	reqBody := `{"username":"testuser","password":"password123"}`
	req, _ := http.NewRequest(http.MethodPost, "/api/auth/register", strings.NewReader(reqBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, req)

	suite.Equal(http.StatusOK, w.Code)
	suite.Contains(w.Header().Get("Set-Cookie"), "accessToken=testtoken")
	suite.mockRegisterHandler.AssertExpectations(suite.T())
}

func (suite *AuthControllerTestSuite) TestLogin_Success() {
	result := &authresult.Result{Token: "testtoken"}
	suite.mockLoginHandler.On("Handle", mock.AnythingOfType("*loginqry.Query")).Return(result, nil)

	reqBody := `{"username":"testuser","password":"password123"}`
	req, _ := http.NewRequest(http.MethodPost, "/api/auth/login", strings.NewReader(reqBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, req)

	suite.Equal(http.StatusOK, w.Code)
	suite.Contains(w.Header().Get("Set-Cookie"), "accessToken=testtoken")
	suite.mockLoginHandler.AssertExpectations(suite.T())
}

func (suite *AuthControllerTestSuite) TestRegisterUser_BadRequest() {
	reqBody := `{"username":"testuser"}`
	req, _ := http.NewRequest(http.MethodPost, "/api/auth/register", strings.NewReader(reqBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, req)

	suite.Equal(http.StatusBadRequest, w.Code)
	suite.mockRegisterHandler.AssertNotCalled(suite.T(), "Handle", mock.Anything)
}

func (suite *AuthControllerTestSuite) TestLogin_BadRequest() {
	reqBody := `{"username":"testuser"}`
	req, _ := http.NewRequest(http.MethodPost, "/api/auth/login", strings.NewReader(reqBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, req)

	suite.Equal(http.StatusBadRequest, w.Code)
	suite.mockLoginHandler.AssertNotCalled(suite.T(), "Handle", mock.Anything)
}

func (suite *AuthControllerTestSuite) TestLogOut() {
	req, _ := http.NewRequest(http.MethodPost, "/api/auth/logOut", nil)
	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, req)

	suite.Equal(http.StatusNoContent, w.Code)
	suite.Contains(w.Header().Get("Set-Cookie"), "accessToken=;")
	suite.Contains(w.Header().Get("Set-Cookie"), "Max-Age=0")
}

func TestAuthControllerTestSuite(t *testing.T) {
	suite.Run(t, new(AuthControllerTestSuite))
}
