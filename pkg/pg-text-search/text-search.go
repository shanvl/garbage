package pgtextsearch

import (
	"strings"
	"unicode"
)

// prepareTextSearchClass processes the query so as to make it a valid argument for to_tsquery,
// adding ':*' to the end of each word and concatenating the words with ' & '. If the input contains invalid symbols,
// it simply returns an empty string
func PrepareQuery(q string) string {
	ss := strings.Fields(q)
	for i, s := range ss {
		if !IsValidInput(s) {
			return ""
		}
		ss[i] = s + ":*"
	}
	return strings.Join(ss, " & ")
}

// isValidInput checks if the given string consists only of digits, letters, "'" and "-"
func IsValidInput(s string) bool {
	for _, r := range s {
		if !unicode.IsLetter(r) && !unicode.IsDigit(r) && r != '\'' && r != '-' {
			return false
		}
	}
	return true
}
