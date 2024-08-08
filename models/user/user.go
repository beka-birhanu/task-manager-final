/*
Package usermodel defines the `User` aggregate, which represents an individual user,
and includes methods for creating and managing users. It handles user creation,
validation of usernames and passwords, and association with expenses.

Key Components:
  - User: Represents a user with details such as username, password hash, and
    associated role.
  - Config: Holds the mandatory parameters required to create a new User.
  - New: Creates a new User instance based on the provided configuration.
  - ConfigBSON: Holds parameters for creating a User with an existing password hash.
  - ToBSON: Creates a new User instance with the provided configuration where the
    password is already hashed.

Dependencies:
- github.com/google/uuid: Used for generating unique IDs.
- github.com/nbutton23/zxcvbn-go: Used for password strength evaluation.
*/
package usermodel

import (
	"regexp"

	"github.com/beka-birhanu/common/hash"
	err "github.com/beka-birhanu/errors"
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

// User represents the aggregate user with private fields.
type UserBSON struct {
	ID           uuid.UUID `bson:"_id"`
	Username     string    `bson:"username"`
	PasswordHash string    `bson:"passwordHash"`
	IsAdmin      bool      `bson:"isAdmin"`
}

// Config holds all mandatory parameters for creating a new User.
type Config struct {
	Username       string
	PlainPassword  string
	IsAdmin        bool
	PasswordHasher hash.IService
}

// ConfigBSON holds all parameters for creating a User with an existing password hash.
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

// ToBSON creates a new User with the provided configuration, where the password is already hashed.
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

// validateUsername validates the username according to the defined rules.
func validateUsername(username string) error {
	if len(username) < minUsernameLength {
		return err.UsernameTooShort
	}
	if len(username) > maxUsernameLength {
		return err.UsernameTooLong
	}
	if !usernameRegex.MatchString(username) {
		return err.UsernameInvalidFormat
	}
	return nil
}

// validatePassword checks the strength of the password.
func validatePassword(password string) error {
	result := zxcvbn.PasswordStrength(password, nil)
	if result.Score < minPasswordStrengthScore {
		return err.WeakPassword
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

// IsAdmin returns whether the user has administrative privileges.
func (u *User) IsAdmin() bool {
	return u.isAdmin
}

