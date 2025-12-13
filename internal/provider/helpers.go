package provider

import (
	"encoding/binary"
	"fmt"
	"net"
	"strconv"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/diag"

	"github.com/iilei/terraform-provider-idgen/internal/idgen"
)

// applyGrouping inserts dashes between groups of characters in the ID.
// For example, with groupSize=4, "abcdefghij" becomes "abcd-efgh-ij".
func applyGrouping(id string, groupSize int) string {
	if groupSize <= 0 || groupSize >= len(id) {
		return id
	}

	// Remove any existing dashes first
	id = strings.ReplaceAll(id, "-", "")

	var grouped strings.Builder
	for i, char := range id {
		if i > 0 && i%groupSize == 0 {
			grouped.WriteRune('-')
		}
		grouped.WriteRune(char)
	}

	return grouped.String()
}

// stringToSeed converts a string to an int64 seed and returns whether it should be directly encoded.
// This is a wrapper around idgen.StringToSeed for use in the provider package.
func stringToSeed(s string) (int64, bool) {
	return idgen.StringToSeed(s)
}

// stringToCanonicalValue parses a string for canonical proquint encoding.
// Returns (value, byteSize, error) where:
//   - value: the numeric value to encode
//   - byteSize: number of bytes needed (4 for uint32, 8 for uint64)
//   - error: description if parsing failed
//
// Supports:
//   - IPv4 addresses (e.g., "127.0.0.1")~>uint32~>4 bytes~>11 chars
//   - Hexadecimal strings (e.g., "0x7f000001" or "7f000001")~>4 or 8 bytes~>11 or 23 chars
//   - uint32 integers (0-4294967295)~>4 bytes~>11 chars
//   - uint64 integers (4294967296-18446744073709551615)~>8 bytes~>23 chars
func stringToCanonicalValue(s string) (uint64, int, string) {
	// Try parsing as IPv4 address first
	if ip := net.ParseIP(s); ip != nil {
		if ipv4 := ip.To4(); ipv4 != nil {
			return uint64(binary.BigEndian.Uint32(ipv4)), 4, ""
		}
		// IPv6 not supported - the original proquint specification focuses on 32-bit values
		// Users must convert IPv6 to its 128-bit integer representation first
		return 0, 0, "IPv6 addresses are not supported for canonical encoding (convert to integer first)"
	}

	// Check if it looks like hex (starts with 0x or contains a-f/A-F)
	lowerS := strings.ToLower(s)
	isHex := strings.HasPrefix(lowerS, "0x") ||
		strings.ContainsAny(s, "abcdefABCDEF")

	if isHex {
		// Try to parse as hexadecimal (with or without 0x prefix)
		hexStr := s
		if strings.HasPrefix(lowerS, "0x") {
			hexStr = s[2:]
		}
		if val, err := strconv.ParseUint(hexStr, 16, 64); err == nil {
			if val <= 0xFFFFFFFF {
				// uint32 range: 4 bytes~>11 chars (2 proquint words)
				return val, 4, ""
			}
			// uint64 range: 8 bytes~>23 chars (4 proquint words)
			return val, 8, ""
		}
	} else {
		// Try to parse as decimal integer
		if val, err := strconv.ParseUint(s, 10, 64); err == nil {
			if val <= 0xFFFFFFFF {
				// uint32 range: 4 bytes~>11 chars (2 proquint words)
				return val, 4, ""
			}
			// uint64 range: 8 bytes~>23 chars (4 proquint words)
			return val, 8, ""
		}
	}

	// Not a valid canonical value
	return 0, 0, "value must be an IPv4 address, hexadecimal string (e.g., 0x7f000001), or unsigned integer (0-18446744073709551615)"
}

// validateLength validates the requested ID length and returns appropriate diagnostics.
// Returns true if validation passed, false otherwise.
func validateLength(length int64, diags *diag.Diagnostics) bool {
	if length < MinIDLength {
		diags.AddError(
			"Invalid length",
			fmt.Sprintf("Length must be at least %d character", MinIDLength),
		)
		return false
	}
	if length > MaxIDLength {
		diags.AddError(
			"Invalid length",
			fmt.Sprintf("Length exceeds maximum allowed value of %d characters", MaxIDLength),
		)
		return false
	}
	if length > WarnIDLength {
		diags.AddWarning(
			"Unusually long ID",
			fmt.Sprintf("The requested length exceeds %d characters. While withthin the limit up to %d, very long IDs may impact performance and readability.", WarnIDLength, MaxIDLength),
		)
	}
	return true
}
