package idgen

import (
	"math/rand"

	gonanoid "github.com/matoous/go-nanoid/v2"
)

var (
	// DefaultAlphabet is the alphabet used for ID characters by default.
	DefaultAlphabet = "_-0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

	// Alphanumeric is an alphabet of alpha-numerical characters (a-zA-Z0-9).
	Alphanumeric = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

	// Numeric is an alphabet of numerical characters (0-9).
	Numeric = "0123456789"

	// Readable avoids visually confusing characters (excludes 0/O, 1/l/I).
	Readable = "23456789abcdefghijkmnpqrstuvwxyzABCDEFGHJKLMNPQRSTUVWXYZ"

	// Legacy alphabets for internal use
	AlphaNum      = Alphanumeric
	Alpha         = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	AlphaLowerNum = "abcdefghijklmnopqrstuvwxyz0123456789"
	AlphaUpperNum = "ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	AlphaLower    = "abcdefghijklmnopqrstuvwxyz"
	AlphaUpper    = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
)

// GenerateNanoID generates a NanoID with the given alphabet and length.
// If seed is non-nil, it generates a deterministic (seeded) ID.
// Otherwise, it uses crypto/rand for secure random generation.
func GenerateNanoID(alphabet string, length int, seed *int64) (string, error) {
	if seed != nil {
		// Seeded mode: deterministic generation using math/rand
		return generateSeededNanoID(*seed, alphabet, length), nil
	}

	// Unseeded mode: use go-nanoid with crypto/rand
	return gonanoid.Generate(alphabet, length)
}

// generateSeededNanoID creates a deterministic NanoID using a seed.
// This is NOT cryptographically secure and should only be used for
// testing or reproducible infrastructure patterns.
func generateSeededNanoID(seed int64, alphabet string, length int) string {
	rng := rand.New(rand.NewSource(seed))
	result := make([]byte, length)
	alphabetLen := len(alphabet)

	for i := 0; i < length; i++ {
		result[i] = alphabet[rng.Intn(alphabetLen)]
	}

	return string(result)
}
