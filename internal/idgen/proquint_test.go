package idgen

import (
	"encoding/binary"
	"net"
	"testing"
)

// TestProquintCanonicalExamples tests the canonical IP-to-proquint mappings
// from the original proquint specification: https://arxiv.org/html/0901.4016
func TestProquintCanonicalExamples(t *testing.T) {
	tests := []struct {
		ip       string
		expected string
	}{
		{"127.0.0.1", "lusab-babad"},
		{"63.84.220.193", "gutih-tugad"},
		{"63.118.7.35", "gutuk-bisog"},
		{"140.98.193.141", "mudof-sakat"},
		{"64.255.6.200", "haguz-biram"},
		{"128.30.52.45", "mabiv-gibot"},
		{"147.67.119.2", "natag-lisaf"},
		{"212.58.253.68", "tibup-zujah"},
		{"216.35.68.215", "tobog-higil"},
		{"216.68.232.21", "todah-vobij"},
		{"198.81.129.136", "sinid-makam"},
		{"12.110.110.204", "budov-kuras"},
	}

	for _, tt := range tests {
		t.Run(tt.ip, func(t *testing.T) {
			// Convert IP to uint32
			ip := net.ParseIP(tt.ip)
			if ip == nil {
				t.Fatalf("Failed to parse IP: %s", tt.ip)
			}
			ipv4 := ip.To4()
			if ipv4 == nil {
				t.Fatalf("Not a valid IPv4 address: %s", tt.ip)
			}

			seed := int64(binary.BigEndian.Uint32(ipv4))

			// Generate proquint with direct encoding
			result, err := GenerateProquint(0, &seed, true)
			if err != nil {
				t.Fatalf("GenerateProquint failed: %v", err)
			}

			if result != tt.expected {
				t.Errorf("IP %s: expected %s, got %s", tt.ip, tt.expected, result)
			}
		})
	}
}

func TestProquintDirectEncoding(t *testing.T) {
	tests := []struct {
		name     string
		value    int64
		expected string
	}{
		{"localhost", 2130706433, "lusab-babad"}, // 127.0.0.1 as uint32
		{"zero", 0, "babab-babab"},
		{"max uint32", 4294967295, "zuzuz-zuzuz"}, // 255.255.255.255
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := GenerateProquint(0, &tt.value, true)
			if err != nil {
				t.Fatalf("GenerateProquint failed: %v", err)
			}

			if result != tt.expected {
				t.Errorf("Value %d: expected %s, got %s", tt.value, tt.expected, result)
			}
		})
	}
}

// TestProquintSeededGeneration tests that seeded random generation is deterministic
func TestProquintSeededGeneration(t *testing.T) {
	seed := int64(42)
	length := 6 // 3 proquint words

	// Generate twice with the same seed
	result1, err1 := GenerateProquint(length, &seed, false)
	if err1 != nil {
		t.Fatalf("First generation failed: %v", err1)
	}

	result2, err2 := GenerateProquint(length, &seed, false)
	if err2 != nil {
		t.Fatalf("Second generation failed: %v", err2)
	}

	if result1 != result2 {
		t.Errorf("Seeded generation not deterministic: %s != %s", result1, result2)
	}

	// Different seed should produce different result
	differentSeed := int64(43)
	result3, err3 := GenerateProquint(length, &differentSeed, false)
	if err3 != nil {
		t.Fatalf("Third generation failed: %v", err3)
	}

	if result1 == result3 {
		t.Errorf("Different seeds produced same result: %s", result1)
	}
}

// TestProquintUnseeded tests that unseeded generation produces valid output
func TestProquintUnseeded(t *testing.T) {
	length := 6 // 3 proquint words

	result, err := GenerateProquint(length, nil, false)
	if err != nil {
		t.Fatalf("Unseeded generation failed: %v", err)
	}

	// Check that we got a non-empty result
	if result == "" {
		t.Error("Unseeded generation returned empty string")
	}

	// Check approximate length (with dashes)
	// 6 bytes = 3 words * 5 chars = 15 chars + 2 dashes = 17 chars
	expectedLength := (length / 2 * 5) + (length/2 - 1)
	if len(result) != expectedLength {
		t.Errorf("Unseeded generation length: expected %d, got %d", expectedLength, len(result))
	}
}

// TestDirectEncodingWithNonCanonicalLength tests padding and truncation in direct encoding mode
func TestDirectEncodingWithNonCanonicalLength(t *testing.T) {
	seed := int64(2130706433) // 127.0.0.1~>canonical: "lusab-babad" (4 bytes, 11 chars)

	tests := []struct {
		name       string
		byteLength int
		expected   string
		desc       string
	}{
		{"canonical_4bytes", 4, "lusab-babad", "Exact canonical length"},
		{"padded_8bytes", 8, "babab-babab-lusab-babad", "Zero-padded to 8 bytes"},
		{"truncated_2bytes", 2, "babad", "Truncated to 2 bytes (rightmost)"},
		{"padded_6bytes", 6, "babab-lusab-babad", "Zero-padded to 6 bytes"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := GenerateProquint(tt.byteLength, &seed, true)
			if err != nil {
				t.Fatalf("Generation with byteLength=%d failed: %v", tt.byteLength, err)
			}

			if result != tt.expected {
				t.Errorf("%s: expected %s, got %s", tt.desc, tt.expected, result)
			}

			// Verify length matches byte length
			// Each 2 bytes = 1 word (5 chars), words separated by dash
			expectedLen := (tt.byteLength / 2 * 5) + (tt.byteLength/2 - 1)
			if len(result) != expectedLen {
				t.Errorf("Length mismatch: expected %d chars, got %d", expectedLen, len(result))
			}
		})
	}
}

// TestGenerateCanonicalProquint tests the canonical proquint generation for various bit sizes
func TestGenerateCanonicalProquint(t *testing.T) {
	tests := []struct {
		name     string
		value    uint64
		byteSize int
		expected string
	}{
		// 32-bit values (4 bytes~>2 words~>11 chars)
		{"zero_32bit", 0, 4, "babab-babab"},
		{"localhost_32bit", 2130706433, 4, "lusab-babad"},
		{"max_uint32", 4294967295, 4, "zuzuz-zuzuz"},

		// 64-bit values (8 bytes~>4 words~>23 chars)
		{"just_above_uint32", 4294967296, 8, "babab-babad-babab-babab"},
		{"max_int64", 9223372036854775807, 8, "luzuz-zuzuz-zuzuz-zuzuz"},
		{"max_uint64", 18446744073709551615, 8, "zuzuz-zuzuz-zuzuz-zuzuz"},

		// Boundary testing
		{"uint32_max_in_64bit", 4294967295, 8, "zuzuz-zuzuz"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := GenerateCanonicalProquint(tt.value)
			if err != nil {
				t.Fatalf("GenerateCanonicalProquint failed: %v", err)
			}

			if result != tt.expected {
				t.Errorf("Value %d: expected %s, got %s",
					tt.value, tt.expected, result)
			}
		})
	}
}

// TestCanonicalProquintLengths verifies output lengths for different bit sizes
func TestCanonicalProquintLengths(t *testing.T) {
	tests := []struct {
		byteSize       int
		expectedLength int
		description    string
	}{
		{4, 11, "4 bytes (32 bits)~>2 words~>11 chars"},
		{8, 23, "8 bytes (64 bits)~>4 words~>23 chars"},
	}

	for _, tt := range tests {
		t.Run(tt.description, func(t *testing.T) {
			// Use a value that fits in the expected byte size
			var testValue uint64
			if tt.byteSize == 4 {
				testValue = 12345 // uint32 range
			} else {
				testValue = 0x100000000 // Just above uint32 range
			}
			result, err := GenerateCanonicalProquint(testValue)
			if err != nil {
				t.Fatalf("GenerateCanonicalProquint failed: %v", err)
			}

			if len(result) != tt.expectedLength {
				t.Errorf("Value %d: expected length %d, got %d (result: %s)",
					testValue, tt.expectedLength, len(result), result)
			}
		})
	}
}

func TestGenerateSeededBytes_NonMultipleOf8(t *testing.T) {
	// Test the missing branch in generateSeededBytes for when remaining < 8
	testCases := []struct {
		name   string
		length int
		seed   int64
	}{
		{"1 byte", 1, 12345},
		{"3 bytes", 3, 12345},
		{"5 bytes", 5, 12345},
		{"7 bytes", 7, 12345},
		{"9 bytes (8+1)", 9, 12345},
		{"11 bytes (8+3)", 11, 12345},
		{"15 bytes (8+7)", 15, 12345},
		{"17 bytes (16+1)", 17, 12345},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Test via GenerateProquint which calls generateSeededBytes
			result, err := GenerateProquint(tc.length, &tc.seed, false)
			if err != nil {
				t.Fatalf("GenerateProquint failed: %v", err)
			}

			// Result should be non-empty and properly formed proquint
			if result == "" {
				t.Error("Expected non-empty result")
			}

			// Test determinism - same seed should produce same result
			result2, err2 := GenerateProquint(tc.length, &tc.seed, false)
			if err2 != nil {
				t.Fatalf("Second GenerateProquint failed: %v", err2)
			}

			if result != result2 {
				t.Errorf("Non-deterministic result: got %q and %q", result, result2)
			}

			// Different seeds should produce different results
			differentSeed := tc.seed + 1
			result3, err3 := GenerateProquint(tc.length, &differentSeed, false)
			if err3 != nil {
				t.Fatalf("Third GenerateProquint failed: %v", err3)
			}

			if result == result3 {
				t.Errorf("Same result for different seeds: %q", result)
			}
		})
	}
}

func TestGenerateSeededBytes_EdgeCases(t *testing.T) {
	t.Run("zero length", func(t *testing.T) {
		result, err := GenerateProquint(0, nil, false)
		if err != nil {
			t.Fatalf("GenerateProquint with zero length failed: %v", err)
		}
		// Should produce empty result or minimal proquint
		if len(result) > 10 { // Arbitrary reasonable limit for empty case
			t.Errorf("Expected short result for zero length, got: %q", result)
		}
	})

	t.Run("large non-multiple of 8", func(t *testing.T) {
		// Test a larger number that's not a multiple of 8
		seed := int64(42)
		result, err := GenerateProquint(23, &seed, false) // 23 = 16 + 7
		if err != nil {
			t.Fatalf("GenerateProquint failed: %v", err)
		}

		if result == "" {
			t.Error("Expected non-empty result")
		}

		// Should be deterministic
		result2, err2 := GenerateProquint(23, &seed, false)
		if err2 != nil {
			t.Fatalf("Second call failed: %v", err2)
		}

		if result != result2 {
			t.Errorf("Not deterministic: %q vs %q", result, result2)
		}
	})
}

func TestGenerateProquint_UnseededRandomGeneration(t *testing.T) {
	// Test the unseeded (crypto/rand) path
	t.Run("unseeded generation", func(t *testing.T) {
		result1, err1 := GenerateProquint(4, nil, false)
		if err1 != nil {
			t.Fatalf("First unseeded GenerateProquint failed: %v", err1)
		}

		result2, err2 := GenerateProquint(4, nil, false)
		if err2 != nil {
			t.Fatalf("Second unseeded GenerateProquint failed: %v", err2)
		}

		// Results should be different (since using crypto/rand)
		if result1 == result2 {
			t.Errorf("Expected different results from crypto/rand, got same: %q", result1)
		}

		// Both should be valid proquints
		if result1 == "" {
			t.Error("First result is empty")
		}
		if result2 == "" {
			t.Error("Second result is empty")
		}
	})

	t.Run("unseeded with various lengths", func(t *testing.T) {
		lengths := []int{1, 2, 4, 8, 16}
		for _, length := range lengths {
			result, err := GenerateProquint(length, nil, false)
			if err != nil {
				t.Errorf("GenerateProquint with length %d failed: %v", length, err)
				continue
			}
			if result == "" {
				t.Errorf("Empty result for length %d", length)
			}
		}
	})
}
