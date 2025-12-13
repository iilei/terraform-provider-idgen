package idgen

import (
	"sort"

	"github.com/iilei/terraform-provider-idgen/internal/data"
)

// GetWordBySeed returns a random word from the wordlist using a deterministic seed
func GetWordBySeed(seed string, wordlist []string) string {
	// Use default wordlist if none provided
	if len(wordlist) == 0 {
		wordlist = data.FiveLetterWords
	}

	if len(wordlist) == 0 {
		return ""
	}

	// Sort to ensure deterministic ordering
	sort.Strings(wordlist)

	seedVal, _ := StringToSeed(seed)
	wordCount := int64(len(wordlist))

	// Use the seed directly to pick an index
	index := int(seedVal % wordCount)
	if index < 0 {
		index += int(wordCount) // Handle negative seeds
	}

	return wordlist[index]
}
