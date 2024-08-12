/*
Package jwt provides functionality for creating and validating JSON Web Tokens (JWTs).
It includes methods for generating tokens with user claims and decoding tokens to
extract claims.

Key Components:
  - Service: Implements JWT operations including token generation and validation.
  - Config: Holds the configuration for creating a new JWT Service.
  - New: Creates a new Service instance with the given configuration.
  - Generate: Generates a JWT for a given user with specified claims.
  - Decode: Decodes and validates a JWT, returning the claims if valid.

Dependencies:
- github.com/dgrijalva/jwt-go: Library for working with JWTs.
- github.com/beka-birhanu/domain/models/user: User model for JWT claims.
*/
package jwt

import (
	"errors"
	"time"

	ijwt "github.com/beka-birhanu/task_manager_final/src/app/common/i_jwt"
	usermodel "github.com/beka-birhanu/task_manager_final/src/domain/models/user"
	"github.com/dgrijalva/jwt-go"
)

// Service handles JWT operations.
// Implements ijwt.Service.
type Service struct {
	secretKey string
	issuer    string
	expTime   time.Duration
}

// Ensure Service implements ijwt.Service.
var _ ijwt.Service = &Service{}

// Config holds JWT service configuration.
type Config struct {
	SecretKey string
	Issuer    string
	ExpTime   time.Duration
}

// New creates a new JWT Service with the provided configuration.
func New(config Config) *Service {
	return &Service{
		secretKey: config.SecretKey,
		issuer:    config.Issuer,
		expTime:   config.ExpTime,
	}
}

// Generate creates a JWT for the given user.
func (s *Service) Generate(user *usermodel.User) (string, error) {
	expirationTime := time.Now().UTC().Add(s.expTime).Unix()
	claims := jwt.MapClaims{
		"user_id":  user.ID().String(),
		"is_admin": user.IsAdmin(),
		"exp":      expirationTime,
		"iss":      s.issuer,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(s.secretKey))
}

// Decode parses and validates a JWT, returning the claims if valid.
func (s *Service) Decode(tokenString string) (jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenString, s.getSigningKey)
	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, errors.New("invalid token")
}

// getSigningKey returns the signing key for token validation.
func (s *Service) getSigningKey(token *jwt.Token) (interface{}, error) {
	if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
		return nil, errors.New("unexpected signing method")
	}
	return []byte(s.secretKey), nil
}
