package data

import (
	"strings"
	"testing"
)

func TestFiveLetterWords(t *testing.T) {
	t.Run("init function loads words and not comments", func(t *testing.T) {
		if len(FiveLetterWords) == 0 {
			t.Error("FiveLetterWords should not be empty after init")
		}

		// Verify all words are 5 letters and not comments
		for i, word := range FiveLetterWords {
			if len(word) != 5 {
				t.Errorf("FiveLetterWords[%d] = %q, length = %d, want 5", i, word, len(word))
			}

			// Verify no empty or whitespace-only words
			if strings.TrimSpace(word) != word {
				t.Errorf("FiveLetterWords[%d] = %q contains whitespace", i, word)
			}

			// Verify no comment
			if strings.HasPrefix(word, "#") {
				t.Errorf("FiveLetterWords[%d] = %q starts with #", i, word)
			}
		}

		// Verify list is sorted
		for i := 1; i < len(FiveLetterWords); i++ {
			if FiveLetterWords[i-1] >= FiveLetterWords[i] {
				t.Errorf("FiveLetterWords not sorted: %q >= %q at index %d", FiveLetterWords[i-1], FiveLetterWords[i], i-1)
			}
		}
	})

	t.Run("contains words", func(t *testing.T) {
		// Test that some common five-letter words are present
		expectedWords := []string{"brave", "shiny", "windy"}

		wordMap := make(map[string]bool)
		for _, word := range FiveLetterWords {
			wordMap[word] = true
		}

		for _, expected := range expectedWords {
			if !wordMap[expected] {
				t.Errorf("Expected word %q not found in FiveLetterWords", expected)
			}
		}
	})

	t.Run("no duplicates", func(t *testing.T) {
		wordMap := make(map[string]bool)
		for i, word := range FiveLetterWords {
			if wordMap[word] {
				t.Errorf("Duplicate word %q found at index %d", word, i)
			}
			wordMap[word] = true
		}
	})
}
