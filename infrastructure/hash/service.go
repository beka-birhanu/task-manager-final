/*
Package hash provides a service for securely hashing and matching passwords
using PBKDF2 with SHA-256. It includes a singleton service that ensures only one
instance of the hashing service is created.

Key Components:
  - Service: Provides methods for hashing and matching passwords.
  - SingletonService: Returns the singleton instance of the Service.
  - Hash: Generates a secure hash for a given word.
  - Match: Compares a plain text word to a hashed word to verify a match.

Dependencies:
- golang.org/x/crypto/pbkdf2: Used for key derivation.
- crypto/sha256: Used as the hash function in PBKDF2.
*/
package hash

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"sync"

	"golang.org/x/crypto/pbkdf2"
)

const (
	iterations = 10000
	saltSize   = 16
	keySize    = 32
)

// Service implements the hash.IService interface.
type Service struct{}

var (
	instance *Service
	once     sync.Once
)

// SingletonService returns a singleton instance of the hash Service.
func SingletonService() *Service {
	once.Do(func() {
		instance = &Service{}
	})
	return instance
}

// Hash generates a hashed representation of the given word.
func (hs *Service) Hash(word string) (string, error) {
	salt := make([]byte, saltSize)
	if _, err := rand.Read(salt); err != nil {
		return "", err
	}

	hash := pbkdf2.Key([]byte(word), salt, iterations, keySize, sha256.New)
	result := append(salt, hash...)

	return base64.StdEncoding.EncodeToString(result), nil
}

// Match compares a plain text word to a hashed word to determine if they match.
func (hs *Service) Match(hashedWord, plainWord string) (bool, error) {
	hashedWordBytes, err := base64.StdEncoding.DecodeString(hashedWord)
	if err != nil {
		return false, err
	}

	if len(hashedWordBytes) != saltSize+keySize {
		return false, errors.New("invalid hashed word length")
	}

	salt := hashedWordBytes[:saltSize]
	expectedHash := hashedWordBytes[saltSize:]

	hash := pbkdf2.Key([]byte(plainWord), salt, iterations, keySize, sha256.New)
	for i := 0; i < keySize; i++ {
		if expectedHash[i] != hash[i] {
			return false, nil
		}
	}

	return true, nil
}

