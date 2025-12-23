// Package idgen provides identifier generation functions for Proquint and NanoID formats.
package idgen

import (
	"encoding/binary"
	mathrand "math/rand/v2"

	"github.com/syrupyy/proquint"
)

func generateSeededBytes(seed int64, length int) []byte {
	rng := mathrand.New(mathrand.NewPCG(uint64(seed), uint64(seed)))
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

func generateRandomBytes(length int) []byte {
	// Use global math/rand/v2 for non-cryptographic random generation
	bytes := make([]byte, length)
	for i := 0; i < length; i += 8 {
		val := mathrand.Uint64()
		remaining := length - i
		if remaining >= 8 {
			binary.LittleEndian.PutUint64(bytes[i:], val)
		} else {
			temp := make([]byte, 8)
			binary.LittleEndian.PutUint64(temp, val)
			copy(bytes[i:], temp[:remaining])
		}
	}
	return bytes
}

// GenerateProquint generates a Proquint ID with the given byte length.
//
// Behavior:
//   - If seed is non-nil and directEncode is true: encodes the seed value directly as bytes.
//     If byteLength differs from canonical size, the output is adjusted (padded with zeros or truncated).
//   - If seed is non-nil and directEncode is false: generates deterministic random bytes using the seed.
//   - If seed is nil: uses math/rand/v2 for random generation (NOT cryptographically secure).
func GenerateProquint(byteLength int, seed *int64, directEncode bool) (string, error) {
	var bytes []byte

	if seed != nil && directEncode {
		// Direct encoding mode: use canonical encoding
		value := uint64(*seed)

		// Determine canonical byte size
		var canonicalBytes []byte
		if value > 0xFFFFFFFF {
			// 64-bit canonical encoding
			canonicalBytes = make([]byte, 8)
			binary.BigEndian.PutUint64(canonicalBytes, value)
		} else {
			// 32-bit canonical encoding
			canonicalBytes = make([]byte, 4)
			binary.BigEndian.PutUint32(canonicalBytes, uint32(value))
		}

		canonicalByteLength := len(canonicalBytes)

		// Adjust to requested byte length if different
		if byteLength > 0 && byteLength != canonicalByteLength {
			if byteLength < canonicalByteLength {
				// Truncate: take the rightmost bytes (least significant)
				bytes = canonicalBytes[canonicalByteLength-byteLength:]
			} else {
				// Pad: prepend zero bytes (most significant)
				bytes = make([]byte, byteLength)
				copy(bytes[byteLength-canonicalByteLength:], canonicalBytes)
			}
		} else {
			bytes = canonicalBytes
		}
	} else if seed != nil {
		// Seeded random generation
		bytes = generateSeededBytes(*seed, byteLength)
	} else {
		// Unseeded: math/rand/v2 (non-cryptographic)
		bytes = generateRandomBytes(byteLength)
	}

	return proquint.EncodeBytes(bytes, "-"), nil
}

// GenerateCanonicalProquint generates a canonical Proquint from a uint64 value.
// The output length is automatically determined by the value:
//   - Values 0-4294967295 (uint32 range): 4 bytes~>11 characters (2 proquint words)
//   - Values 4294967296+ (uint64 range): 8 bytes~>23 characters (4 proquint words)
//
// This implements the canonical encoding described in the original proquint specification.
func GenerateCanonicalProquint(value uint64) (string, error) {
	var bytes []byte

	if value > 0xFFFFFFFF {
		// 64-bit encoding: 8 bytes~>4 words~>23 chars
		bytes = make([]byte, 8)
		binary.BigEndian.PutUint64(bytes, value)
	} else {
		// 32-bit encoding: 4 bytes~>2 words~>11 chars
		bytes = make([]byte, 4)
		binary.BigEndian.PutUint32(bytes, uint32(value))
	}

	return proquint.EncodeBytes(bytes, "-"), nil
}
