// Package loginqry handles the login query, including user authentication
// and JWT generation.
package loginqry

import (
	"fmt"

	iquery "github.com/beka-birhanu/task_manager_final/src/app/common/cqrs/query"
	ijwt "github.com/beka-birhanu/task_manager_final/src/app/common/i_jwt"
	irepo "github.com/beka-birhanu/task_manager_final/src/app/common/i_repo"
	authresult "github.com/beka-birhanu/task_manager_final/src/app/user/auth/common"
	errdmn "github.com/beka-birhanu/task_manager_final/src/domain/errors"
	ihash "github.com/beka-birhanu/task_manager_final/src/domain/i_hash"
)

// Handler processes login queries.
type Handler struct {
	userRepo irepo.User    // Repository for user data operations.
	jwtSvc   ijwt.Service  // Service for JWT operations.
	hashSvc  ihash.Service // Service for password hashing.
}

// Ensure Handler implements iquery.Handler
var _ iquery.IHandler[*Query, *authresult.Result] = &Handler{}

// Config holds the dependencies for creating a new Handler.
type Config struct {
	UserRepo irepo.User    // Repository for user data operations.
	JwtSvc   ijwt.Service  // Service for JWT operations.
	HashSvc  ihash.Service // Service for password hashing.
}

// NewHandler creates a new Handler with the given configuration.
func NewHandler(cfg Config) *Handler {
	return &Handler{
		userRepo: cfg.UserRepo,
		jwtSvc:   cfg.JwtSvc,
		hashSvc:  cfg.HashSvc,
	}
}

// Handle processes a login query, authenticates the user, and generates a JWT.
func (s *Handler) Handle(qry *Query) (*authresult.Result, error) {
	// Retrieve user by username.
	user, err := s.userRepo.ByUsername(qry.Username)
	if err != nil {
		return nil, errdmn.NewUnauthorized(err.Error())
	}

	// Validate the provided password.
	isPasswordCorrect, err := s.hashSvc.Match(user.PasswordHash(), qry.Password)
	if err != nil {
		errMessage := fmt.Sprintf("failed to validate user password, %v", err)
		return nil, errdmn.NewUnexpected(errMessage)
	}

	if !isPasswordCorrect {
		return nil, errdmn.NewUnauthorized("incorrect password")
	}

	// Generate JWT for the authenticated user.
	token, err := s.jwtSvc.Generate(user)
	if err != nil {
		errMessage := fmt.Sprintf("failed to generate JWT for user, %v", err)
		return nil, errdmn.NewUnexpected(errMessage)
	}

	return authresult.New(user, token), nil
}
