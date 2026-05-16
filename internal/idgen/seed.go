package idgen

import (
	"encoding/binary"
	"hash/fnv"
	"math/big"
	"net"
	"strconv"
)

// StringToSeed converts a string to an int64 seed and returns whether it should be directly encoded.
// Returns (seed, shouldDirectEncode).
// Handles three cases:
// 1. IPv4 address (e.g., "127.0.0.1") - converts to uint32, directly encoded as bytes
// 2. Decimal integer in uint32 range (0-4294967295) - directly encoded as bytes (canonical proquint)
// 3. Arbitrary string or large integer - hashed deterministically, used as random seed
func StringToSeed(s string) (int64, bool) {
	// Try parsing as IPv4 address first - these should be encoded directly
	if ip := net.ParseIP(s); ip != nil {
		if ipv4 := ip.To4(); ipv4 != nil {
			return int64(binary.BigEndian.Uint32(ipv4)), true
		}
	}

	// Try to parse as integer
	if val, err := strconv.ParseInt(s, 10, 64); err == nil {
		// If it fits in uint32 range (0-4294967295), encode directly (canonical proquint behavior)
		if val >= 0 && val <= 0xFFFFFFFF {
			return val, true
		}
		// Large integers are used as random seed
		return val, false
	}

	// Hash the string to get a deterministic seed - used as random seed
	h := fnv.New64a()
	h.Write([]byte(s))
	return int64(h.Sum64()), false
}

// StringToSeedBigInt converts a string to a big.Int seed and returns whether it should be directly encoded.
// Returns (seed, shouldDirectEncode).
// Handles three cases:
// 1. IPv4 address (e.g., "127.0.0.1") - converts to uint32, directly encoded as bytes
// 2. Decimal integer in uint32 range (0-4294967295) - directly encoded as bytes (canonical proquint)
// 3. Arbitrary string or large integer - hashed deterministically, used as random seed
func StringToSeedBigInt(s string) (*big.Int, bool) {
	// Try parsing as IPv4 address first - these should be encoded directly
	if ip := net.ParseIP(s); ip != nil {
		if ipv4 := ip.To4(); ipv4 != nil {
			return new(big.Int).SetUint64(uint64(binary.BigEndian.Uint32(ipv4))), true
		}
	}

	// Try to parse as an arbitrarily large integer
	if value, ok := new(big.Int).SetString(s, 10); ok {
		if value.Sign() >= 0 && value.Cmp(new(big.Int).SetUint64(0xFFFFFFFF)) <= 0 {
			return value, true
		}
		return value, false
	}

	// Hash the string to get a deterministic seed - used as random seed
	h := fnv.New64a()
	h.Write([]byte(s))
	return new(big.Int).SetUint64(h.Sum64()), false
}
