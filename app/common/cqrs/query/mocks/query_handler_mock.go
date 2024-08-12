package iquery

import "github.com/stretchr/testify/mock"

// IHandler is a mock implementation of the IHandler interface using testify's mock functionalities.
type IHandler[Query any, Result any] struct {
	mock.Mock
}

// Handle mocks the Handle method of the IHandler interface.
func (m *IHandler[Query, Result]) Handle(query Query) (Result, error) {
	args := m.Called(query)
	return args.Get(0).(Result), args.Error(1)
}
