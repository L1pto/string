package stringtools

import (
	"strings"
	"unicode"
)

func ToLower(s string) string {
	var b strings.Builder

	for _, symbol := range s {
		b.WriteRune(unicode.ToLower(symbol))
	}
	return b.String()
}

func SplitToRunes(s string) []rune {
	return []rune(s)
}

func CleanNonLetters(s string) string {
	var b strings.Builder
	for _, symbol := range s {
		if unicode.IsLetter(symbol) {
			b.WriteRune(symbol)
		}
	}

	return b.String()
}

func SplitWords(s string) []string {
	if len(s) == 0 {
		return []string{}
	}

	result := strings.Fields(s)

	return result
}
