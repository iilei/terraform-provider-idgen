package idgen

import (
	"math"
	"math/rand/v2"
	"strings"

	gonanoid "github.com/matoous/go-nanoid/v2"
)

var (
	// Alphanumeric is an alphabet of alpha-numerical characters (a-zA-Z0-9).
	Alphanumeric = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

	// Numeric is an alphabet of numerical characters (0-9).
	Numeric = "0123456789"

	// Readable avoids visually confusing characters (excludes 0/O, 1/l/I).
	Readable = "23456789abcdefghkmnpqrstwxyzABCDEFGHJKLMNPQRSTWXYZ"
)

// GenerateNanoID generates a NanoID with the given alphabet and length.
// If groupSize > 0, the length parameter represents the final length including
// separators, and this function automatically applies grouping to the generated ID.
// If seed is non-nil, it generates a deterministic (seeded) ID.
// Otherwise, it uses crypto/rand for secure random generation.
func GenerateNanoID(alphabet string, length int, seed *int64, groupSize int) (string, error) {
	// Calculate internal length if grouping is enabled
	internalLength := length
	if groupSize > 0 {
		internalLength = int(math.Ceil(float64(length*groupSize+1) / float64(groupSize+1)))
	}

	var id string
	var err error

	if seed != nil {
		// Seeded mode: deterministic generation using math/rand
		id = generateSeededNanoID(*seed, alphabet, internalLength)
	} else {
		// Unseeded mode: use go-nanoid with crypto/rand
		id, err = gonanoid.Generate(alphabet, internalLength)
		if err != nil {
			return "", err
		}
	}

	// Apply grouping if requested
	if groupSize > 0 {
		id = ApplyGrouping(id, groupSize)
	}

	return id, nil
}

// ApplyGrouping inserts dashes between groups of characters in the ID.
// For example, with groupSize=4, "abcdefghij" becomes "abcd-efgh-ij".
// Note: The caller is responsible for removing any existing dashes before calling this function.
func ApplyGrouping(id string, groupSize int) string {
	if groupSize <= 0 || groupSize >= len(id) {
		return id
	}

	var grouped strings.Builder
	for i, char := range id {
		if i > 0 && i%groupSize == 0 {
			grouped.WriteRune('-')
		}
		grouped.WriteRune(char)
	}

	return grouped.String()
}

// generateSeededNanoID creates a deterministic NanoID using a seed.
// This is NOT cryptographically secure and should only be used for
// testing or reproducible infrastructure patterns.
func generateSeededNanoID(seed int64, alphabet string, length int) string {
	rng := rand.New(rand.NewPCG(uint64(seed), uint64(seed)))
	result := make([]byte, length)
	alphabetLen := len(alphabet)

	for i := 0; i < length; i++ {
		result[i] = alphabet[rng.IntN(alphabetLen)]
	}

	return string(result)
}
