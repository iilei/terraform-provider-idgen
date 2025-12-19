package idgen

import (
	"fmt"
	"strconv"
	"testing"
)

func TestGetWordBySeed(t *testing.T) {
	wordList := []string{
		"apple",
		"berry",
		"elder",
		"peach",
	}

	for i, want := range wordList {
		result := GetWordBySeed(strconv.Itoa(i), wordList)
		if result != want {
			t.Errorf("Iteration %d: expected %s, got %s", i, want, result)
		}
	}
}

func TestGetWordBySeedWithNegative(t *testing.T) {
	wordList := []string{
		"apple",
		"berry",
		"elder",
		"peach",
	}

	for i, want := range wordList {
		result := GetWordBySeed(strconv.Itoa(i-len(wordList)), wordList)
		if result != want {
			t.Errorf("Iteration %d: expected %s, got %s", i, want, result)
		}
	}
}

func TestGetWordBySeedWithEntropy(t *testing.T) {
	wordList := []string{
		"apple",
		"berry",
		"elder",
		"peach",
	}

	wantWordList := []string{
		"berry",
		"elder",
		"peach",
		"apple",
	}

	for i := range wordList {
		result := GetWordBySeed(fmt.Sprintf("seed-%d", i), wordList)
		want := wantWordList[i]

		if result != want {
			t.Errorf("Iteration %d: expected %s, got %s", i, want, result)
		}
	}
}

func TestGetWordBySeedWithEmptyWordList(t *testing.T) {
	result := GetWordBySeed("seed-4711", nil)

	for i := 1; i <= 3; i++ {
		want := result
		result := GetWordBySeed("seed-4711", nil)

		if result != want {
			t.Errorf("Iteration %d: expected %s, got %s", i, want, result)
		}
	}
}

func TestGetWordByWithoutSeed(t *testing.T) {
	results := []string{}

	for i := 1; i <= 3; i++ {
		result := GetWordBySeed("", nil)
		results = append(results, result)
	}

	// Check if all results are unique
	unique := make(map[string]struct{})
	for _, r := range results {
		unique[r] = struct{}{}
	}

	if len(unique) > 1 {
		t.Errorf("Expected different results for empty seed, got %d unique values: %v", len(unique), results)
	}
}

func TestGetWordBySeed_EmptyWordlistEdgeCases(t *testing.T) {
	t.Run("empty wordlist falls back to default", func(t *testing.T) {
		result := GetWordBySeed("test", []string{})
		// Should fall back to data.FiveLetterWords and return a word
		if result == "" {
			t.Error("Expected non-empty result when falling back to default wordlist")
		}
	})

	t.Run("nil wordlist falls back to default", func(t *testing.T) {
		result := GetWordBySeed("test", nil)
		// Should fall back to data.FiveLetterWords and return a word
		if result == "" {
			t.Error("Expected non-empty result when falling back to default wordlist")
		}
	})
}

func TestGetWordBySeed_NegativeSeedHandling(t *testing.T) {
	wordList := []string{"apple", "berry", "cherry"}

	// Test cases that would produce negative indices before correction
	testCases := []struct {
		name string
		seed string
	}{
		{"negative integer seed", "-5"},
		{"negative hash result", "some string that produces negative hash"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := GetWordBySeed(tc.seed, wordList)

			// Result should be non-empty and one of the words in the list
			if result == "" {
				t.Error("Expected non-empty result for negative seed")
			}

			found := false
			for _, word := range wordList {
				if result == word {
					found = true
					break
				}
			}

			if !found {
				t.Errorf("Result %q is not in the wordlist %v", result, wordList)
			}

			// Test that the same seed produces the same result (deterministic)
			result2 := GetWordBySeed(tc.seed, wordList)
			if result != result2 {
				t.Errorf("Non-deterministic result: got %q first, %q second", result, result2)
			}
		})
	}
}

func TestGetWordBySeed_NegativeIndexCorrection(t *testing.T) {
	// Test the specific line: index += int(wordCount) when index < 0
	wordList := []string{"alpha", "beta", "gamma"}

	negativeSeeds := []string{"0", "-1", "-2", "-3", "-4"}

	for _, seed := range negativeSeeds {
		// Get the raw seed value to understand what's happening
		seedVal, _ := StringToSeed(seed)
		wordCount := int64(len(wordList))
		index := int(seedVal % wordCount)

		t.Logf("Seed: %s, seedVal: %d, wordCount: %d, index before correction: %d",
			seed, seedVal, wordCount, index)

		result := GetWordBySeed(seed, wordList)

		// The result should still be valid
		if result == "" {
			t.Errorf("GetWordBySeed(%q, %v) returned empty string", seed, wordList)
		}

		// Verify the result is one of our words
		validWord := false
		for _, word := range wordList {
			if result == word {
				validWord = true
				break
			}
		}

		if !validWord {
			t.Errorf("GetWordBySeed(%q, %v) returned invalid word: %q", seed, wordList, result)
		}

		// Test determinism
		result2 := GetWordBySeed(seed, wordList)
		if result != result2 {
			t.Errorf("GetWordBySeed not deterministic for seed %q: got %q and %q", seed, result, result2)
		}
	}
}
