package ihash

type Service interface {
	// Hash hashes the provided plain text string and returns the hashed value or
	// An error if the hashing process fails.
	Hash(word string) (string, error)

	// Match checks if the hashed string matches the plain text string.
	// Returns:
	// - true if the hashed string matches the plain text string.
	// - false if they do not match.
	// - An error if the comparison process fails.
	Match(hashedWord, plainWord string) (bool, error)
}
