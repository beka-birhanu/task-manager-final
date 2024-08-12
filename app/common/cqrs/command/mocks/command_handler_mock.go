package icmd_mock

import "github.com/stretchr/testify/mock"

// IHandler is a mock implementation of the IHandler interface using testify.Mock.
type IHandler[Command any, Result any] struct {
	mock.Mock
}

// Handle processes the command using testify's mock functionalities.
func (m *IHandler[Command, Result]) Handle(command Command) (Result, error) {
	args := m.Called(command)
	return args.Get(0).(Result), args.Error(1)
}
