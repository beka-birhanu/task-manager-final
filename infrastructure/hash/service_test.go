package hash_test

import (
	"testing"

	"github.com/beka-birhanu/task_manager_final/infrastructure/hash"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestService_HashAndMatch(t *testing.T) {
	service := hash.SingletonService()

	t.Run("Hash generates a valid hash", func(t *testing.T) {
		word := "password123"
		hashedWord, err := service.Hash(word)
		require.NoError(t, err)
		assert.NotEmpty(t, hashedWord, "Hashed word should not be empty")
	})

	t.Run("Match returns true for correct password", func(t *testing.T) {
		word := "password123"
		hashedWord, err := service.Hash(word)
		require.NoError(t, err)

		matches, err := service.Match(hashedWord, word)
		require.NoError(t, err)
		assert.True(t, matches, "Match should return true for correct password")
	})

	t.Run("Match returns false for incorrect password", func(t *testing.T) {
		word := "password123"
		hashedWord, err := service.Hash(word)
		require.NoError(t, err)

		matches, err := service.Match(hashedWord, "wrongpassword")
		require.NoError(t, err)
		assert.False(t, matches, "Match should return false for incorrect password")
	})

	t.Run("Match returns error for invalid hash format", func(t *testing.T) {
		_, err := service.Match("invalidHash", "password123")
		assert.Error(t, err, "Match should return an error for an invalid hash format")
	})
}
