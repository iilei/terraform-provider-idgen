// Package idgen provides identifier generation functions for Proquint and NanoID formats.
package idgen

import (
	"encoding/binary"
	"hash/fnv"
	"math/big"
	mathrand "math/rand/v2"

	"github.com/syrupyy/proquint"
)

func generateSeededBytes(seed uint64, length int) []byte {
	rng := mathrand.New(mathrand.NewPCG(seed, seed))
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

func canonicalBigIntBytes(seed *big.Int) []byte {
	if seed == nil {
		return nil
	}

	if seed.Sign() == 0 {
		return make([]byte, 4)
	}

	if seed.BitLen() <= 32 {
		bytes := make([]byte, 4)
		seed.FillBytes(bytes)
		return bytes
	}

	if seed.BitLen() <= 64 {
		bytes := make([]byte, 8)
		seed.FillBytes(bytes)
		return bytes
	}

	return seed.Bytes()
}

// seededUint64 reduces a seed to the 64-bit state used by the PCG generator.
//
// Values that fit in int64 or uint64 map exactly as they did when seeds were int64,
// so existing seeded IDs are unchanged. Wider values are hashed (FNV-1a over the
// big-endian bytes, with the sign mixed in) instead of being truncated to their low
// 64 bits, which would make e.g. 2^64+1 and 2^65+1 collide.
func seededUint64(seed *big.Int) uint64 {
	if seed == nil {
		return 0
	}

	if seed.IsInt64() {
		return uint64(seed.Int64())
	}

	if seed.IsUint64() {
		return seed.Uint64()
	}

	h := fnv.New64a()
	if seed.Sign() < 0 {
		h.Write([]byte{'-'})
	}
	h.Write(seed.Bytes())
	return h.Sum64()
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
//     If byteLength differs from the canonical or requested size, the output is adjusted (padded with zeros or truncated).
//   - If seed is non-nil and directEncode is false: generates deterministic random bytes using the seed.
//   - If seed is nil: uses math/rand/v2 for random generation (NOT cryptographically secure).
func GenerateProquint(byteLength int, seed *big.Int, directEncode bool) (string, error) {
	var bytes []byte

	if seed != nil && directEncode {
		// Direct encoding mode: use canonical big-endian encoding for uint32/uint64-sized values.
		bytes = canonicalBigIntBytes(seed)
		canonicalByteLength := len(bytes)

		// Adjust to requested byte length if different
		if byteLength > 0 && byteLength != canonicalByteLength {
			if byteLength < canonicalByteLength {
				// Truncate: take the rightmost bytes (least significant)
				bytes = bytes[canonicalByteLength-byteLength:]
			} else {
				// Pad: prepend zero bytes (most significant)
				bytes = make([]byte, byteLength)
				copy(bytes[byteLength-canonicalByteLength:], canonicalBigIntBytes(seed))
			}
		}
	} else if seed != nil {
		// Seeded random generation
		bytes = generateSeededBytes(seededUint64(seed), byteLength)
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
	seed := new(big.Int).SetUint64(value)
	return GenerateProquint(0, seed, true)
}
