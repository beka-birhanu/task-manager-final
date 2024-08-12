// Package registercmd provides the command, handler structure and factory function
// for user registration.
package registercmd

import (
	icmd "github.com/beka-birhanu/app/common/cqrs/command"
	ijwt "github.com/beka-birhanu/app/common/i_jwt"
	irepo "github.com/beka-birhanu/app/common/i_repo"
	authresult "github.com/beka-birhanu/app/user/auth/common"
	ihash "github.com/beka-birhanu/domain/i_hash"
	usermodel "github.com/beka-birhanu/domain/models/user"
)

// Handler handles the user registration process.
type Handler struct {
	userRepo irepo.User    // Repository for user data operations.
	jwtSvc   ijwt.Service  // Service for JWT operations.
	hashSvc  ihash.Service // Service for password hashing.
}

// Ensure Handler implementes icmd.Handler
var _ icmd.IHandler[*Command, *authresult.Result] = &Handler{}

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

// Handle processes a registration command, creating and saving a new user,
// and generating a JWT for the user.
func (h *Handler) Handle(cmd *Command) (*authresult.Result, error) {
	var isAdmin bool

	// Check if there are any users to determine if the new user should be an admin.
	count, err := h.userRepo.Count()
	if err != nil {
		return nil, err
	} else if count == 0 {
		isAdmin = true
	}

	// Create a new user.
	user, err := createUser(cmd, h.hashSvc, isAdmin)
	if err != nil {
		return nil, err
	}

	// Save the new user to the repository.
	if err := h.userRepo.Save(user); err != nil {
		return nil, err
	}

	// Generate a JWT for the new user.
	token, err := h.jwtSvc.Generate(user)
	if err != nil {
		return nil, err
	}

	return authresult.New(user, token), nil
}

// createUser creates a new user with the provided command data and password hashing service.
func createUser(cmd *Command, hashSvc ihash.Service, isAdmin bool) (*usermodel.User, error) {
	cfg := usermodel.Config{
		Username:       cmd.Username,
		PlainPassword:  cmd.Password,
		IsAdmin:        isAdmin,
		PasswordHasher: hashSvc,
	}
	return usermodel.New(cfg)
}
