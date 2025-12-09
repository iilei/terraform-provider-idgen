package provider

import "testing"

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
			name:      "remove existing dashes",
			id:        "abc-def-ghi",
			groupSize: 4,
			expected:  "abcd-efgh-i",
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
			result := applyGrouping(tt.id, tt.groupSize)
			if result != tt.expected {
				t.Errorf("applyGrouping(%q, %d) = %q, want %q", tt.id, tt.groupSize, result, tt.expected)
			}
		})
	}
}
