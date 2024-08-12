// Package authresult defines the structure for authentication results,
// including user details and an authentication token.
package authresult

import (
	usermodel "github.com/beka-birhanu/task_manager_final/src/domain/models/user"
	"github.com/google/uuid"
)

// Result represents the outcome of an authentication process,
// including the user's ID, username, and an authentication token.
type Result struct {
	ID       uuid.UUID // Unique identifier for the user.
	Username string    // Username of the authenticated user.
	Token    string    // Authentication token.
	IsAdmin  bool
}

// New creates a new authentication result with the given user and token.
func New(user *usermodel.User, token string) *Result {
	return &Result{
		ID:       user.ID(),
		Username: user.Username(),
		Token:    token,
		IsAdmin:  user.IsAdmin(),
	}
}
