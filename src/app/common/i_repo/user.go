// Package irepo provides interfaces for user repository operations.
package irepo

import (
	usermodel "github.com/beka-birhanu/domain/models/user"
	"github.com/google/uuid"
)

// User defines methods to manage users in the store.
type User interface {
	Save(user *usermodel.User) error
	ById(id uuid.UUID) (*usermodel.User, error)
	ByUsername(username string) (*usermodel.User, error)
	Count() (int64, error)
}

