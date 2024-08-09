// Package err provides a simple way to create and handle custom domain errors.
// It includes predefined error types like Validation, Conflict, Unexpected, NotFound, and Unauthorized,
// which can be used to categorize errors in a consistent manner.
package err

import (
	"fmt"
)

const (
	// Validation represents a validation error.
	Validation = "Validation"

	// Conflict represents a conflict error.
	Conflict = "Conflict"

	// Unexpected represents an unexpected server error.
	Unexpected = "ServerError"

	// NotFound represents a resource not found error.
	NotFound = "NotFound"

	// Unauthorized represents an unauthorized error.
	Unauthorized = "Unauthorized"
)

// Error represents a custom domain error with a type and message.
type Error struct {
	kind    string
	Message string
}

// new creates a new Error with the given type and message.
func new(errType string, message string) *Error {
	return &Error{kind: errType, Message: message}
}

// Error returns the string representation of the Error.
func (e Error) Error() string {
	return fmt.Sprintf("%s: %s", e.kind, e.Message)
}

// Type returns the type of the Error.
func (e Error) Type() string {
	return e.kind
}

// NewValidation creates a new validation error with the given message.
func NewValidation(message string) *Error {
	return new(Validation, message)
}

// NewConflict creates a new conflict error with the given message.
func NewConflict(message string) *Error {
	return new(Conflict, message)
}

// NewUnexpected creates a new unexpected server error with the given message.
func NewUnexpected(message string) *Error {
	return new(Unexpected, message)
}

// NewNotFound creates a new not found error with the given message.
func NewNotFound(message string) *Error {
	return new(NotFound, message)
}

// NewUnauthorized creates a new unauthorized error with the given message.
func NewUnauthorized(message string) *Error {
	return new(Unauthorized, message)
}

