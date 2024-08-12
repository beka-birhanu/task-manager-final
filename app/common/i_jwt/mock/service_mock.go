package ijwt_mock

import (
	usermodel "github.com/beka-birhanu/task_manager_final/domain/models/user"
	"github.com/dgrijalva/jwt-go"
	"github.com/stretchr/testify/mock"
)

// MockService is a mock implementation of the Service interface using testify.
type MockService struct {
	mock.Mock
}

// Generate mocks the Generate method of the Service interface.
func (m *MockService) Generate(user *usermodel.User) (string, error) {
	args := m.Called(user)
	return args.String(0), args.Error(1)
}

// Decode mocks the Decode method of the Service interface.
func (m *MockService) Decode(token string) (jwt.MapClaims, error) {
	args := m.Called(token)
	return args.Get(0).(jwt.MapClaims), args.Error(1)
}
