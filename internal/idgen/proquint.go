// Package idgen provides identifier generation functions for Proquint and NanoID formats.
package idgen

import (
	"encoding/binary"
	"hash/fnv"
	"math/big"
	mathrand "math/rand/v2"

	"github.com/syrupyy/proquint"
)

const (
	uint32Bytes = 4
	uint64Bytes = 8
	uint32Bits  = 32
	uint64Bits  = 64
)

// fillUint64s returns length bytes taken from successive little-endian uint64 values of next.
func fillUint64s(length int, next func() uint64) []byte {
	bytes := make([]byte, length)
	var word [uint64Bytes]byte

	for i := 0; i < length; i += uint64Bytes {
		binary.LittleEndian.PutUint64(word[:], next())
		copy(bytes[i:], word[:])
	}

	return bytes
}

func generateSeededBytes(seed uint64, length int) []byte {
	// Deterministic, non-cryptographic generation is the point of a seed.
	rng := mathrand.New(mathrand.NewPCG(seed, seed)) //nolint:gosec // seeded IDs must be reproducible
	return fillUint64s(length, rng.Uint64)
}

func generateRandomBytes(length int) []byte {
	// IDs are documented as not cryptographically secure; see the provider's security notice.
	return fillUint64s(length, mathrand.Uint64)
}

func canonicalBigIntBytes(seed *big.Int) []byte {
	if seed == nil {
		return nil
	}

	var size int
	switch bits := seed.BitLen(); {
	case bits <= uint32Bits:
		size = uint32Bytes
	case bits <= uint64Bits:
		size = uint64Bytes
	default:
		return seed.Bytes()
	}

	bytes := make([]byte, size)
	seed.FillBytes(bytes)
	return bytes
}

// seededUint64 reduces a seed to the 64-bit state used by the PCG generator.
//
// Values that fit in int64 or uint64 map exactly as they did when seeds were int64,
// so existing seeded IDs are unchanged. Wider values are hashed (FNV-1a over the
// big-endian bytes, with the sign mixed in) instead of being truncated to their low
// 64 bits, which would make e.g. 2^64+1 and 2^65+1 collide.
func seededUint64(seed *big.Int) uint64 {
	switch {
	case seed == nil:
		return 0
	case seed.IsInt64():
		return uint64(seed.Int64()) //nolint:gosec // bit reinterpretation, same as the former int64 seed path
	case seed.IsUint64():
		return seed.Uint64()
	}

	h := fnv.New64a()
	if seed.Sign() < 0 {
		_, _ = h.Write([]byte{'-'}) // hash.Hash.Write never returns an error
	}
	_, _ = h.Write(seed.Bytes())
	return h.Sum64()
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

	switch {
	case seed != nil && directEncode:
		bytes = fitBytes(canonicalBigIntBytes(seed), byteLength)
	case seed != nil:
		bytes = generateSeededBytes(seededUint64(seed), byteLength)
	default:
		// Unseeded: math/rand/v2 (non-cryptographic)
		bytes = generateRandomBytes(byteLength)
	}

	return proquint.EncodeBytes(bytes, "-"), nil
}

// fitBytes adjusts big-endian bytes to length: truncation keeps the least significant
// bytes, padding prepends zeros. A length of 0 or less keeps bytes unchanged.
func fitBytes(bytes []byte, length int) []byte {
	switch {
	case length <= 0 || length == len(bytes):
		return bytes
	case length < len(bytes):
		return bytes[len(bytes)-length:]
	default:
		padded := make([]byte, length)
		copy(padded[length-len(bytes):], bytes)
		return padded
	}
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
