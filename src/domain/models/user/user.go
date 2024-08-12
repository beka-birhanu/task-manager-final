/*
Package usermodel defines the `User` aggregate, representing an individual user with methods
for creation and management. It handles user creation, username and password validation,
and user-expense associations.

Key Components:
  - User: Represents a user with details like username, password hash, and role.
  - Config: Holds parameters required to create a new User.
  - New: Creates a new User instance using the provided configuration.
  - ConfigBSON: Holds parameters for creating a User with an existing password hash.
  - ToBSON: Creates a User instance with a pre-hashed password.

Dependencies:
- github.com/google/uuid: For generating unique IDs.
- github.com/nbutton23/zxcvbn-go: For evaluating password strength.
*/
package usermodel

import (
	"regexp"

	errdmn "github.com/beka-birhanu/domain/errors"
	"github.com/beka-birhanu/domain/i_hash"
	"github.com/google/uuid"
	"github.com/nbutton23/zxcvbn-go"
)

const (
	minPasswordStrengthScore = 3

	usernamePattern   = `^[a-zA-Z0-9_]+$` // Alphanumeric with underscores
	minUsernameLength = 3
	maxUsernameLength = 20
)

var (
	usernameRegex = regexp.MustCompile(usernamePattern)
)

// User represents the aggregate user with private fields.
type User struct {
	id           uuid.UUID
	username     string
	passwordHash string
	isAdmin      bool
}

// UserBSON represents the BSON version of the User for database storage.
type UserBSON struct {
	ID           uuid.UUID `bson:"_id"`
	Username     string    `bson:"username"`
	PasswordHash string    `bson:"passwordHash"`
	IsAdmin      bool      `bson:"isAdmin"`
}

// Config holds parameters for creating a new User.
type Config struct {
	Username       string
	PlainPassword  string
	IsAdmin        bool
	PasswordHasher ihash.Service
}

// ConfigBSON holds parameters for creating a User with an existing password hash.
type ConfigBSON struct {
	ID           uuid.UUID
	Username     string
	PasswordHash string
	IsAdmin      bool
}

// New creates a new User with the provided configuration.
func New(config Config) (*User, error) {
	if err := validateUsername(config.Username); err != nil {
		return nil, err
	}

	if err := validatePassword(config.PlainPassword); err != nil {
		return nil, err
	}

	passwordHash, err := config.PasswordHasher.Hash(config.PlainPassword)
	if err != nil {
		return nil, err
	}

	return &User{
		id:           uuid.New(), // New ID for the user
		username:     config.Username,
		passwordHash: passwordHash,
		isAdmin:      config.IsAdmin,
	}, nil
}

// ToBSON creates a new User with a pre-hashed password.
func ToBSON(config ConfigBSON) (*User, error) {
	if err := validateUsername(config.Username); err != nil {
		return nil, err
	}

	return &User{
		id:           config.ID,
		username:     config.Username,
		passwordHash: config.PasswordHash,
		isAdmin:      config.IsAdmin,
	}, nil
}

// FromBSON creates a User from a BSON representation.
func FromBSON(bsonUser *UserBSON) *User {
	return &User{
		id:           bsonUser.ID,
		username:     bsonUser.Username,
		passwordHash: bsonUser.PasswordHash,
		isAdmin:      bsonUser.IsAdmin,
	}
}

// validateUsername validates the username.
func validateUsername(username string) error {
	if len(username) < minUsernameLength {
		return errdmn.UsernameTooShort
	}
	if len(username) > maxUsernameLength {
		return errdmn.UsernameTooLong
	}
	if !usernameRegex.MatchString(username) {
		return errdmn.UsernameInvalidFormat
	}
	return nil
}

// validatePassword checks the strength of the password.
func validatePassword(password string) error {
	result := zxcvbn.PasswordStrength(password, nil)
	if result.Score < minPasswordStrengthScore {
		return errdmn.WeakPassword
	}
	return nil
}

// ID returns the user's ID.
func (u *User) ID() uuid.UUID {
	return u.id
}

// Username returns the user's username.
func (u *User) Username() string {
	return u.username
}

// PasswordHash returns the user's password hash.
func (u *User) PasswordHash() string {
	return u.passwordHash
}

// IsAdmin returns whether the user is an admin.
func (u *User) IsAdmin() bool {
	return u.isAdmin
}

// UpdateUsername updates the user's username after validation.
func (u *User) UpdateUsername(newUsername string) error {
	if err := validateUsername(newUsername); err != nil {
		return err
	}
	u.username = newUsername
	return nil
}

// UpdatePassword updates the user's password after validation.
func (u *User) UpdatePassword(newPassword string, passwordHasher ihash.Service) error {
	if err := validatePassword(newPassword); err != nil {
		return err
	}

	hashedPassword, err := passwordHasher.Hash(newPassword)
	if err != nil {
		return err
	}

	u.passwordHash = hashedPassword
	return nil
}

// UpdateAdminStatus updates the user's admin status.
func (u *User) UpdateAdminStatus(isAdmin bool) {
	u.isAdmin = isAdmin
}

