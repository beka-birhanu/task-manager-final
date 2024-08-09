package ijwt

import (
	usermodel "github.com/beka-birhanu/models/user"
	"github.com/dgrijalva/jwt-go"
)

type Service interface {
	// Generate creates a JWT for the specified user.
	Generate(user *usermodel.User) (string, error)

	// Decode parses the provided JWT and returns the claims or an error.
	Decode(token string) (jwt.MapClaims, error)
}

