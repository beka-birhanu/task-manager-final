package authmiddleware_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	authmiddleware "github.com/beka-birhanu/task_manager_final/api/middleware/auth"
	ijwt_mock "github.com/beka-birhanu/task_manager_final/app/common/i_jwt/mock"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

// AuthMiddlewareTestSuite defines the test suite for the Authoriz middleware.
type AuthMiddlewareTestSuite struct {
	mockJwtSvc *ijwt_mock.MockService
	router     *gin.Engine
}

// SetupTest sets up the test environment.
func (suite *AuthMiddlewareTestSuite) SetupTest(hasToBeAdmin bool) {
	suite.mockJwtSvc = new(ijwt_mock.MockService)

	// Set up the Gin router with the middleware.
	suite.router = gin.Default()
	suite.router.Use(authmiddleware.Authoriz(suite.mockJwtSvc, hasToBeAdmin))

	// Example endpoint to test the middleware
	suite.router.GET("/test", func(c *gin.Context) {
		claims, _ := c.Get(authmiddleware.ContextUserClaims)
		c.JSON(http.StatusOK, claims)
	})
}

// TestValidToken tests the scenario where a valid token is provided.
func TestValidToken(t *testing.T) {
	suite := new(AuthMiddlewareTestSuite)
	suite.SetupTest(false)

	// Mock the JWT service to return valid claims.
	claims := jwt.MapClaims{"username": "testuser", "is_admin": false}
	suite.mockJwtSvc.On("Decode", "valid_token").Return(claims, nil)

	// Create a new request with a valid token cookie.
	req, _ := http.NewRequest(http.MethodGet, "/test", nil)
	req.AddCookie(&http.Cookie{Name: "accessToken", Value: "valid_token"})
	w := httptest.NewRecorder()

	// Serve the request.
	suite.router.ServeHTTP(w, req)

	// Assert that the response status is OK (200) and claims are returned.
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), `"username":"testuser"`)
	assert.Contains(t, w.Body.String(), `"is_admin":false`)

	// Verify the mock expectations.
	suite.mockJwtSvc.AssertExpectations(t)
}

// TestAdminRequired_Failure tests the scenario where admin access is required but the user is not an admin.
func TestAdminRequired_Failure(t *testing.T) {
	suite := new(AuthMiddlewareTestSuite)
	suite.SetupTest(true) // Admin access is required

	// Mock the JWT service to return valid claims but the user is not an admin.
	claims := jwt.MapClaims{"username": "testuser", "is_admin": false}
	suite.mockJwtSvc.On("Decode", "valid_token").Return(claims, nil)

	// Create a new request with a valid token cookie.
	req, _ := http.NewRequest(http.MethodGet, "/test", nil)
	req.AddCookie(&http.Cookie{Name: "accessToken", Value: "valid_token"})
	w := httptest.NewRecorder()

	// Serve the request.
	suite.router.ServeHTTP(w, req)

	// Assert that the response status is Forbidden (403).
	assert.Equal(t, http.StatusForbidden, w.Code)

	// Verify the mock expectations.
	suite.mockJwtSvc.AssertExpectations(t)
}
