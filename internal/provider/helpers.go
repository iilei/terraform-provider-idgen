package provider

import "strings"

// applyGrouping inserts dashes between groups of characters in the ID.
// For example, with groupSize=4, "abcdefghij" becomes "abcd-efgh-ij".
func applyGrouping(id string, groupSize int) string {
	if groupSize <= 0 || groupSize >= len(id) {
		return id
	}

	// Remove any existing dashes first
	id = strings.ReplaceAll(id, "-", "")

	var grouped strings.Builder
	for i, char := range id {
		if i > 0 && i%groupSize == 0 {
			grouped.WriteRune('-')
		}
		grouped.WriteRune(char)
	}

	return grouped.String()
}
