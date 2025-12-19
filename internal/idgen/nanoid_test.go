package idgen

import (
	"strings"
	"testing"
)

func TestGenerateNanoID(t *testing.T) {
	t.Run("unseeded generation", func(t *testing.T) {
		id, err := GenerateNanoID(Alphanumeric, 21, nil, 0)
		if err != nil {
			t.Fatalf("GenerateNanoID() error = %v", err)
		}

		if len(id) != 21 {
			t.Errorf("GenerateNanoID() length = %d, want 21", len(id))
		}

		// Verify all characters are from the alphabet
		for _, char := range id {
			if !strings.ContainsRune(Alphanumeric, char) {
				t.Errorf("GenerateNanoID() contains invalid character %c", char)
			}
		}
	})

	t.Run("seeded generation is deterministic", func(t *testing.T) {
		seed := int64(12345)
		id1, err1 := GenerateNanoID(Alphanumeric, 21, &seed, 0)
		id2, err2 := GenerateNanoID(Alphanumeric, 21, &seed, 0)

		if err1 != nil {
			t.Fatalf("GenerateNanoID() error1 = %v", err1)
		}
		if err2 != nil {
			t.Fatalf("GenerateNanoID() error2 = %v", err2)
		}

		if id1 != id2 {
			t.Errorf("GenerateNanoID() seeded not deterministic: %q != %q", id1, id2)
		}

		if len(id1) != 21 {
			t.Errorf("GenerateNanoID() length = %d, want 21", len(id1))
		}
	})

	t.Run("with grouping", func(t *testing.T) {
		seed := int64(12345)
		// Request 15 chars total with grouping of 4 -> "xxxx-xxxx-xxx" (11 chars + 2 separators = 13, adjust to fit)
		id, err := GenerateNanoID(Alphanumeric, 13, &seed, 4)
		if err != nil {
			t.Fatalf("GenerateNanoID() error = %v", err)
		}

		// Should contain dashes
		if !strings.Contains(id, "-") {
			t.Errorf("GenerateNanoID() with grouping should contain dashes, got %q", id)
		}

		parts := strings.Split(id, "-")
		if len(parts) < 2 {
			t.Errorf("GenerateNanoID() with grouping should have multiple parts, got %q", id)
		}

		// First parts should be exactly 4 characters
		for i, part := range parts[:len(parts)-1] {
			if len(part) != 4 {
				t.Errorf("GenerateNanoID() part %d length = %d, want 4", i, len(part))
			}
		}
	})

	t.Run("different alphabets", func(t *testing.T) {
		seed := int64(12345)

		idNumeric, err := GenerateNanoID(Numeric, 10, &seed, 0)
		if err != nil {
			t.Fatalf("GenerateNanoID() numeric error = %v", err)
		}

		// Should only contain digits
		for _, char := range idNumeric {
			if !strings.ContainsRune(Numeric, char) {
				t.Errorf("GenerateNanoID() numeric contains invalid character %c", char)
			}
		}

		idReadable, err := GenerateNanoID(Readable, 10, &seed, 0)
		if err != nil {
			t.Fatalf("GenerateNanoID() readable error = %v", err)
		}

		// Should only contain readable characters
		for _, char := range idReadable {
			if !strings.ContainsRune(Readable, char) {
				t.Errorf("GenerateNanoID() readable contains invalid character %c", char)
			}
		}

		// Different alphabets should produce different results
		if idNumeric == idReadable {
			t.Errorf("GenerateNanoID() with different alphabets should produce different results")
		}
	})
}

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
			id:        "abc",
			groupSize: 1,
			expected:  "a-b-c",
		},
		{
			name:      "exact multiple of group size",
			id:        "abcdef",
			groupSize: 2,
			expected:  "ab-cd-ef",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ApplyGrouping(tt.id, tt.groupSize)
			if result != tt.expected {
				t.Errorf("ApplyGrouping(%q, %d) = %q, want %q", tt.id, tt.groupSize, result, tt.expected)
			}
		})
	}
}

func TestGenerateSeededNanoID(t *testing.T) {
	t.Run("deterministic generation", func(t *testing.T) {
		seed := int64(12345)
		id1 := generateSeededNanoID(seed, Alphanumeric, 21)
		id2 := generateSeededNanoID(seed, Alphanumeric, 21)

		if id1 != id2 {
			t.Errorf("generateSeededNanoID() not deterministic: %q != %q", id1, id2)
		}

		if len(id1) != 21 {
			t.Errorf("generateSeededNanoID() length = %d, want 21", len(id1))
		}
	})

	t.Run("different seeds produce different results", func(t *testing.T) {
		id1 := generateSeededNanoID(12345, Alphanumeric, 21)
		id2 := generateSeededNanoID(67890, Alphanumeric, 21)

		if id1 == id2 {
			t.Errorf("generateSeededNanoID() with different seeds should produce different results")
		}
	})

	t.Run("respects alphabet", func(t *testing.T) {
		seed := int64(12345)
		id := generateSeededNanoID(seed, Numeric, 10)

		// Should only contain digits
		for _, char := range id {
			if !strings.ContainsRune(Numeric, char) {
				t.Errorf("generateSeededNanoID() contains invalid character %c", char)
			}
		}

		if len(id) != 10 {
			t.Errorf("generateSeededNanoID() length = %d, want 10", len(id))
		}
	})

	t.Run("empty alphabet", func(t *testing.T) {
		seed := int64(12345)
		// This should panic or produce empty string
		defer func() {
			if r := recover(); r == nil {
				t.Error("generateSeededNanoID() with empty alphabet should panic")
			}
		}()
		generateSeededNanoID(seed, "", 10)
	})

	t.Run("zero length", func(t *testing.T) {
		seed := int64(12345)
		id := generateSeededNanoID(seed, Alphanumeric, 0)

		if id != "" {
			t.Errorf("generateSeededNanoID() with zero length = %q, want empty string", id)
		}
	})
}

func TestGenerateNanoID_ErrorCases(t *testing.T) {
	t.Run("empty alphabet unseeded should error", func(t *testing.T) {
		// This should trigger the error path in GenerateNanoID when calling gonanoid.Generate
		_, err := GenerateNanoID("", 21, nil, 0)
		if err == nil {
			t.Error("Expected error for empty alphabet in unseeded generation, but got none")
		}

		// Error should mention the issue
		if err != nil {
			errMsg := err.Error()
			if errMsg == "" {
				t.Error("Error should have a message")
			}
		}
	})

	t.Run("empty alphabet seeded should panic", func(t *testing.T) {
		// This should panic in generateSeededNanoID due to division by zero in rng.Intn(alphabetLen)
		seed := int64(42)

		defer func() {
			if r := recover(); r == nil {
				t.Error("Expected panic for empty alphabet in seeded generation")
			}
		}()

		GenerateNanoID("", 21, &seed, 0)
	})
}
