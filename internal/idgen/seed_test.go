package idgen

import (
	"encoding/binary"
	"hash/fnv"
	"net"
	"testing"
)

func TestStringToSeed(t *testing.T) {
	t.Run("IPv4 addresses", func(t *testing.T) {
		testCases := []struct {
			name     string
			input    string
			expected int64
		}{
			{
				name:     "localhost IPv4",
				input:    "127.0.0.1",
				expected: int64(binary.BigEndian.Uint32(net.ParseIP("127.0.0.1").To4())),
			},
			{
				name:     "zero IPv4",
				input:    "0.0.0.0",
				expected: int64(binary.BigEndian.Uint32(net.ParseIP("0.0.0.0").To4())),
			},
			{
				name:     "arbitrary IPv4",
				input:    "192.168.1.1",
				expected: int64(binary.BigEndian.Uint32(net.ParseIP("192.168.1.1").To4())),
			},
			{
				name:     "max IPv4",
				input:    "255.255.255.255",
				expected: int64(binary.BigEndian.Uint32(net.ParseIP("255.255.255.255").To4())),
			},
		}

		for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {
				seed, shouldDirectEncode := StringToSeed(tc.input)
				if seed != tc.expected {
					t.Errorf("StringToSeed(%q) seed = %d, want %d", tc.input, seed, tc.expected)
				}
				if !shouldDirectEncode {
					t.Errorf("StringToSeed(%q) shouldDirectEncode = false, want true", tc.input)
				}
			})
		}
	})

	t.Run("integers in uint32 range", func(t *testing.T) {
		testCases := []struct {
			name     string
			input    string
			expected int64
		}{
			{
				name:     "zero",
				input:    "0",
				expected: 0,
			},
			{
				name:     "small positive",
				input:    "42",
				expected: 42,
			},
			{
				name:     "large uint32",
				input:    "4294967295",
				expected: 4294967295,
			},
			{
				name:     "middle range",
				input:    "1000000",
				expected: 1000000,
			},
		}

		for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {
				seed, shouldDirectEncode := StringToSeed(tc.input)
				if seed != tc.expected {
					t.Errorf("StringToSeed(%q) seed = %d, want %d", tc.input, seed, tc.expected)
				}
				if !shouldDirectEncode {
					t.Errorf("StringToSeed(%q) shouldDirectEncode = false, want true", tc.input)
				}
			})
		}
	})

	t.Run("large integers beyond uint32 range", func(t *testing.T) {
		testCases := []struct {
			name     string
			input    string
			expected int64
		}{
			{
				name:     "just above uint32 max",
				input:    "4294967296",
				expected: 4294967296,
			},
			{
				name:     "large positive integer",
				input:    "9223372036854775807",
				expected: 9223372036854775807,
			},
			{
				name:     "negative integer",
				input:    "-1",
				expected: -1,
			},
			{
				name:     "large negative integer",
				input:    "-9223372036854775808",
				expected: -9223372036854775808,
			},
		}

		for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {
				seed, shouldDirectEncode := StringToSeed(tc.input)
				if seed != tc.expected {
					t.Errorf("StringToSeed(%q) seed = %d, want %d", tc.input, seed, tc.expected)
				}
				if shouldDirectEncode {
					t.Errorf("StringToSeed(%q) shouldDirectEncode = true, want false", tc.input)
				}
			})
		}
	})

	t.Run("arbitrary strings", func(t *testing.T) {
		testCases := []struct {
			name  string
			input string
		}{
			{
				name:  "simple string",
				input: "hello",
			},
			{
				name:  "empty string",
				input: "",
			},
			{
				name:  "alphanumeric",
				input: "abc123",
			},
			{
				name:  "special characters",
				input: "hello@world!",
			},
			{
				name:  "unicode string",
				input: "café",
			},
			{
				name:  "long string",
				input: "this is a very long string that should be hashed deterministically",
			},
		}

		for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {
				seed, shouldDirectEncode := StringToSeed(tc.input)

				// Calculate expected hash
				h := fnv.New64a()
				h.Write([]byte(tc.input))
				expectedSeed := int64(h.Sum64())

				if seed != expectedSeed {
					t.Errorf("StringToSeed(%q) seed = %d, want %d", tc.input, seed, expectedSeed)
				}
				if shouldDirectEncode {
					t.Errorf("StringToSeed(%q) shouldDirectEncode = true, want false", tc.input)
				}
			})
		}
	})

	t.Run("consistency - same input produces same output", func(t *testing.T) {
		inputs := []string{
			"127.0.0.1",
			"42",
			"hello world",
			"192.168.1.100",
			"4294967296",
			"",
		}

		for _, input := range inputs {
			seed1, direct1 := StringToSeed(input)
			seed2, direct2 := StringToSeed(input)

			if seed1 != seed2 {
				t.Errorf("StringToSeed(%q) not consistent: got %d and %d", input, seed1, seed2)
			}
			if direct1 != direct2 {
				t.Errorf("StringToSeed(%q) shouldDirectEncode not consistent: got %v and %v", input, direct1, direct2)
			}
		}
	})

	t.Run("edge cases", func(t *testing.T) {
		testCases := []struct {
			name                 string
			input                string
			expectedDirectEncode bool
		}{
			{
				name:                 "IPv6 address should not be direct encoded",
				input:                "2001:db8::1",
				expectedDirectEncode: false,
			},
			{
				name:                 "invalid IP address format",
				input:                "256.256.256.256",
				expectedDirectEncode: false,
			},
			{
				name:                 "integer with leading zeros",
				input:                "00042",
				expectedDirectEncode: true,
			},
			{
				name:                 "integer with plus sign",
				input:                "+42",
				expectedDirectEncode: true,
			},
			{
				name:                 "non-numeric string that looks like number",
				input:                "42abc",
				expectedDirectEncode: false,
			},
			{
				name:                 "floating point number",
				input:                "42.5",
				expectedDirectEncode: false,
			},
		}

		for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {
				_, shouldDirectEncode := StringToSeed(tc.input)
				if shouldDirectEncode != tc.expectedDirectEncode {
					t.Errorf("StringToSeed(%q) shouldDirectEncode = %v, want %v",
						tc.input, shouldDirectEncode, tc.expectedDirectEncode)
				}
			})
		}
	})
}

func TestStringToSeed_IPv4Conversion(t *testing.T) {
	// Test specific IPv4 to seed conversion values
	testCases := []struct {
		ipv4     string
		expected int64
	}{
		{"127.0.0.1", 2130706433},       // 0x7F000001
		{"192.168.1.1", 3232235777},     // 0xC0A80101
		{"10.0.0.1", 167772161},         // 0x0A000001
		{"255.255.255.255", 4294967295}, // 0xFFFFFFFF
		{"0.0.0.0", 0},                  // 0x00000000
	}

	for _, tc := range testCases {
		t.Run(tc.ipv4, func(t *testing.T) {
			seed, shouldDirectEncode := StringToSeed(tc.ipv4)
			if seed != tc.expected {
				t.Errorf("StringToSeed(%q) = %d, want %d", tc.ipv4, seed, tc.expected)
			}
			if !shouldDirectEncode {
				t.Errorf("StringToSeed(%q) shouldDirectEncode = false, want true", tc.ipv4)
			}
		})
	}
}
