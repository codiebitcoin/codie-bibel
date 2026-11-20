package bible

import (
	"strings"
	"unicode/utf8"
)

// fold is what a query folds to before matching: its lower case.
func fold(s string) string {
	return strings.ToLower(s)
}

// indexOfFold finds the folded query in text without case, as the rune
// offset it starts at, or -1. The offset counts text's runes, so the views
// can lay the match over the text as it reads.
func indexOfFold(text, query string) int {
	if query == "" {
		return -1
	}
	lowered := strings.ToLower(text)
	i := strings.Index(lowered, query)
	if i < 0 {
		return -1
	}
	return utf8.RuneCountInString(lowered[:i])
}

// runes is s's runes, as views size a slice of text by. Kept next to
// indexOfFold, whose offset it reads.
func runes(s string) int {
	return utf8.RuneCountInString(s)
}
