// Package ijwt provides JWT generation and decoding services.
package ijwt

import (
	usermodel "github.com/beka-birhanu/task_manager_final/domain/models/user"
	"github.com/dgrijalva/jwt-go"
)

// Service defines methods to generate and decode JWTs.
type Service interface {
	// Generate creates a JWT for a user.
	Generate(user *usermodel.User) (string, error)

	// Decode parses a JWT and returns claims.
	Decode(token string) (jwt.MapClaims, error)
}
