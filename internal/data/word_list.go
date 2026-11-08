// Package data holds the embedded word list used for random word IDs.
package data

import (
	"bufio"
	_ "embed"
	"sort"
	"strings"
)

//go:embed five_letter_words.txt
var fiveLetterWords string

// FiveLetterWords is a predefined word list of common five-letter English words.
// It is sorted and contains no blank values.
var FiveLetterWords []string

// WordSet contains every entry of FiveLetterWords, for membership checks.
var WordSet map[string]struct{}

func init() {
	wordSet := make(map[string]struct{})
	sc := bufio.NewScanner(strings.NewReader(fiveLetterWords))
	for sc.Scan() {
		// Trim whitespace and check if the line is a comment or empty
		if word := strings.TrimSpace(sc.Text()); word != "" && !strings.HasPrefix(word, "#") {
			wordSet[word] = struct{}{}
		}
	}

	// Convert to slice and sort
	FiveLetterWords = make([]string, 0, len(wordSet))
	for word := range wordSet {
		FiveLetterWords = append(FiveLetterWords, word)
	}
	sort.Strings(FiveLetterWords)
}
