package irepo_mock

import (
	usermodel "github.com/beka-birhanu/task_manager_final/domain/models/user"
	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
)

// User is a mock implementation of the User interface using testify.
type User struct {
	mock.Mock
}

// Save mocks the Save method of the User interface.
func (m *User) Save(user *usermodel.User) error {
	args := m.Called(user)
	return args.Error(0)
}

// ById mocks the ById method of the User interface.
func (m *User) ById(id uuid.UUID) (*usermodel.User, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*usermodel.User), args.Error(1)
}

// ByUsername mocks the ByUsername method of the User interface.
func (m *User) ByUsername(username string) (*usermodel.User, error) {
	args := m.Called(username)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*usermodel.User), args.Error(1)
}

// Count mocks the Count method of the User interface.
func (m *User) Count() (int64, error) {
	args := m.Called()
	return args.Get(0).(int64), args.Error(1)
}
