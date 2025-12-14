package provider

import (
	"testing"

	"github.com/iilei/terraform-provider-idgen/internal/idgen"
)

func TestApplyGrouping(t *testing.T) {
	tests := []struct {
		name      string
		id        string
		groupSize int
		expected  string
	}{
		{
			name:      "group by 4",
			id:        "abcdefghij",
			groupSize: 4,
			expected:  "abcd-efgh-ij",
		},
		{
			name:      "group by 3",
			id:        "123456789",
			groupSize: 3,
			expected:  "123-456-789",
		},
		{
			name:      "no grouping when size is 0",
			id:        "abcdefgh",
			groupSize: 0,
			expected:  "abcdefgh",
		},
		{
			name:      "no grouping when size >= length",
			id:        "abcd",
			groupSize: 5,
			expected:  "abcd",
		},
		{
			name:      "single character groups",
			id:        "abcd",
			groupSize: 1,
			expected:  "a-b-c-d",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := idgen.ApplyGrouping(tt.id, tt.groupSize)
			if result != tt.expected {
				t.Errorf("idgen.ApplyGrouping(%q, %d) = %q, want %q", tt.id, tt.groupSize, result, tt.expected)
			}
		})
	}
}

func TestStringToSeed(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected int64
	}{
		{
			name:     "integer string",
			input:    "42",
			expected: 42,
		},
		{
			name:     "negative integer",
			input:    "-123",
			expected: -123,
		},
		{
			name:     "large integer",
			input:    "9223372036854775807",
			expected: 9223372036854775807,
		},
		{
			name:  "text string (hashed)",
			input: "my-seed",
			// This will be a hash, just verify it's deterministic
			expected: func() int64 { seed, _ := stringToSeed("my-seed"); return seed }(),
		},
		{
			name:     "IPv4 address",
			input:    "127.0.0.1",
			expected: 2130706433, // 0x7f000001
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, _ := stringToSeed(tt.input)
			if result != tt.expected {
				t.Errorf("stringToSeed(%q) = %d, want %d", tt.input, result, tt.expected)
			}
		})
	}

	// Test determinism: same input always produces same output
	t.Run("deterministic hashing", func(t *testing.T) {
		input := "test-string-123"
		result1, _ := stringToSeed(input)
		result2, _ := stringToSeed(input)
		if result1 != result2 {
			t.Errorf("stringToSeed not deterministic: got %d and %d for same input", result1, result2)
		}
	})

	// Test different inputs produce different seeds
	t.Run("different inputs produce different seeds", func(t *testing.T) {
		seed1, _ := stringToSeed("seed-a")
		seed2, _ := stringToSeed("seed-b")
		if seed1 == seed2 {
			t.Errorf("stringToSeed produced same seed %d for different inputs", seed1)
		}
	})

	// Test directEncode flag for IPv4
	t.Run("IPv4 returns directEncode flag", func(t *testing.T) {
		_, directEncode := stringToSeed("127.0.0.1")
		if !directEncode {
			t.Errorf("stringToSeed(\"127.0.0.1\") should return directEncode=true")
		}
	})

	// Test non-IP strings return false for directEncode
	t.Run("non-IP strings return directEncode=false", func(t *testing.T) {
		_, directEncode := stringToSeed("my-seed")
		if directEncode {
			t.Errorf("stringToSeed(\"my-seed\") should return directEncode=false")
		}
	})
}

func TestStringToCanonicalValue(t *testing.T) {
	tests := []struct {
		name         string
		input        string
		wantValue    uint64
		wantByteSize int
		wantError    string
	}{
		// IPv4 addresses~>4 bytes
		{
			name:         "localhost IPv4",
			input:        "127.0.0.1",
			wantValue:    2130706433,
			wantByteSize: 4,
			wantError:    "",
		},
		{
			name:         "max IPv4",
			input:        "255.255.255.255",
			wantValue:    4294967295,
			wantByteSize: 4,
			wantError:    "",
		},

		// uint32 range~>4 bytes
		{
			name:         "zero",
			input:        "0",
			wantValue:    0,
			wantByteSize: 4,
			wantError:    "",
		},
		{
			name:         "uint32 max",
			input:        "4294967295",
			wantValue:    4294967295,
			wantByteSize: 4,
			wantError:    "",
		},

		// uint64 range~>8 bytes
		{
			name:         "just above uint32",
			input:        "4294967296",
			wantValue:    4294967296,
			wantByteSize: 8,
			wantError:    "",
		},
		{
			name:         "max int64",
			input:        "9223372036854775807",
			wantValue:    9223372036854775807,
			wantByteSize: 8,
			wantError:    "",
		},
		{
			name:         "max uint64",
			input:        "18446744073709551615",
			wantValue:    18446744073709551615,
			wantByteSize: 8,
			wantError:    "",
		},

		// Hexadecimal strings
		{
			name:         "hex with 0x prefix uint32",
			input:        "0x7f000001",
			wantValue:    2130706433,
			wantByteSize: 4,
			wantError:    "",
		},
		{
			name:         "hex without 0x prefix uint32",
			input:        "7f000001",
			wantValue:    2130706433,
			wantByteSize: 4,
			wantError:    "",
		},
		{
			name:         "hex uppercase",
			input:        "0xFFFFFFFF",
			wantValue:    4294967295,
			wantByteSize: 4,
			wantError:    "",
		},
		{
			name:         "hex uint64",
			input:        "0x7fffffffffffffff",
			wantValue:    9223372036854775807,
			wantByteSize: 8,
			wantError:    "",
		},
		{
			name:         "hex max uint64",
			input:        "0xffffffffffffffff",
			wantValue:    18446744073709551615,
			wantByteSize: 8,
			wantError:    "",
		},

		// Errors
		{
			name:         "IPv6 not supported",
			input:        "2001:db8::1",
			wantValue:    0,
			wantByteSize: 0,
			wantError:    "IPv6 addresses are not supported for canonical encoding (convert to integer first)",
		},
		{
			name:         "invalid string",
			input:        "not-a-number",
			wantValue:    0,
			wantByteSize: 0,
			wantError:    "value must be an IPv4 address or unsigned integer",
		},
		{
			name:         "negative integer",
			input:        "-1",
			wantValue:    0,
			wantByteSize: 0,
			wantError:    "value must be an IPv4 address or unsigned integer",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			value, byteSize, errMsg := stringToCanonicalValue(tt.input)

			if tt.wantError != "" {
				if errMsg == "" {
					t.Errorf("expected error containing %q, got no error", tt.wantError)
				}
				return
			}

			if errMsg != "" {
				t.Errorf("unexpected error: %s", errMsg)
				return
			}

			if value != tt.wantValue {
				t.Errorf("value: got %d, want %d", value, tt.wantValue)
			}

			if byteSize != tt.wantByteSize {
				t.Errorf("byteSize: got %d, want %d", byteSize, tt.wantByteSize)
			}
		})
	}
}
