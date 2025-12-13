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
