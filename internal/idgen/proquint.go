package idgen

import (
	"crypto/rand"
	"encoding/binary"
	mathrand "math/rand"

	"github.com/syrupyy/proquint"
)

// GenerateProquint generates a Proquint ID with the given byte length.
// If seed is non-nil, it generates a deterministic (seeded) ID.
// Otherwise, it uses crypto/rand for secure random generation.
// The separator is always "-".
func GenerateProquint(byteLength int, seed *int64) (string, error) {
	var bytes []byte
	var err error

	if seed != nil {
		// Seeded mode: deterministic generation using math/rand
		bytes = generateSeededBytes(*seed, byteLength)
	} else {
		// Unseeded mode: use crypto/rand
		bytes = make([]byte, byteLength)
		_, err = rand.Read(bytes)
		if err != nil {
			return "", err
		}
	}

	// Encode the bytes as proquint with "-" separator
	return proquint.EncodeBytes(bytes, "-"), nil
}

// generateSeededBytes creates deterministic bytes using a seed.
// This is NOT cryptographically secure and should only be used for
// testing or reproducible infrastructure patterns.
func generateSeededBytes(seed int64, length int) []byte {
	rng := mathrand.New(mathrand.NewSource(seed))
	bytes := make([]byte, length)

	for i := 0; i < length; i += 8 {
		val := rng.Uint64()
		remaining := length - i
		if remaining >= 8 {
			binary.LittleEndian.PutUint64(bytes[i:], val)
		} else {
			// Handle remaining bytes
			temp := make([]byte, 8)
			binary.LittleEndian.PutUint64(temp, val)
			copy(bytes[i:], temp[:remaining])
		}
	}

	return bytes
}
