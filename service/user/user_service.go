package usersvc

import (
	"fmt"

	"github.com/beka-birhanu/common/hash"
	ijwt "github.com/beka-birhanu/common/i_jwt"
	irepo "github.com/beka-birhanu/common/i_repo"
	usermodel "github.com/beka-birhanu/models/user"
	"github.com/google/uuid"
)

type AuthCommand struct {
	Username string
	Password string
}

type AuthResult struct {
	ID       uuid.UUID
	Username string
	Token    string
}

func NewAuthResult(id uuid.UUID, username, token string) *AuthResult {
	return &AuthResult{
		ID:       id,
		Username: username,
		Token:    token,
	}
}

func NewCommand(username, password string) (*AuthCommand, error) {
	return &AuthCommand{Username: username, Password: password}, nil
}

type Service struct {
	userRepo irepo.User
	jwtSvc   ijwt.Service
	hashSvc  hash.IService
}

type Config struct {
	UserRepo irepo.User
	JwtSvc   ijwt.Service
	HashSvc  hash.IService
}

func NewService(cfg Config) *Service {
	return &Service{
		userRepo: cfg.UserRepo,
		jwtSvc:   cfg.JwtSvc,
		hashSvc:  cfg.HashSvc,
	}
}

func (s *Service) Register(cmd *AuthCommand) (*AuthResult, error) {
	user, err := createUser(cmd, s.hashSvc)
	if err != nil {
		return nil, fmt.Errorf("creating new user failed: %w", err)
	}

	if err := s.userRepo.Save(user); err != nil {
		return nil, fmt.Errorf("saving user to repository failed: %w", err)
	}

	token, err := s.jwtSvc.Generate(user)
	if err != nil {
		return nil, fmt.Errorf("JWT generation failed: %w", err)
	}

	return NewAuthResult(user.ID(), user.Username(), token), nil
}

func (s *Service) SignIn(cmd *AuthCommand) (*AuthResult, error) {
	user, err := s.userRepo.ByUsername(cmd.Username)
	if err != nil {
		return nil, err
	}

	isPasswordCorrect, err := s.hashSvc.Match(user.PasswordHash(), cmd.Password)
	if err != nil {
		return nil, fmt.Errorf("failed to validate user password, %w", err)
	}

	if !isPasswordCorrect {
		return nil, fmt.Errorf("incorrect password")
	}

	token, err := s.jwtSvc.Generate(user)
	if err != nil {
		return nil, fmt.Errorf("failed to generate JWT for user, %w", err)
	}

	return NewAuthResult(user.ID(), user.Username(), token), nil
}

func (s *Service) Promote(username string) error {
	user, err := s.userRepo.ByUsername(username)
	if err != nil {
		return err
	}

	user.UpdateAdminStatus(true)
	err = s.userRepo.Save(user)
	return err
}

func createUser(cmd *AuthCommand, hashSvc hash.IService) (*usermodel.User, error) {
	cfg := usermodel.Config{
		Username:       cmd.Username,
		PlainPassword:  cmd.Password,
		PasswordHasher: hashSvc,
	}
	return usermodel.New(cfg)
}
