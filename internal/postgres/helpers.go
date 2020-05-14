package postgres

import (
	"fmt"
	"strings"
	"time"
	"unicode"

	"github.com/shanvl/garbage-events-service/internal/garbage"
)

// prepareFilterQuery processes a query to make it a valid argument for to_tsquery,
// adding ':*' at the end of each word, concatenating the words with ' & ' and,
// if a word resembles a school class name,
// creating a copy of it with changes needed to hit the indices of the tables. For example,
// if the passed date is 10.10.2020, a query "3A Iv Ig" will be changed to "3A:* | 2018A:* & Iv:* & Ig:*"
func prepareFilterQuery(q string, t time.Time) string {
	ss := strings.Fields(q)
	for i, s := range ss {
		// if the word starts with a number, try to parse it as a class name and concatenate with itself
		if unicode.IsDigit(rune(s[0])) {
			// if the word can be parsed as a class name, process it and concatenate with itself.
			// "3B" becomes "3B:* | 2018A:*"
			letter, formed, err := garbage.ParseClassName(s, t)
			// error means that the word doesn't resemble a class name,
			// so it can be skipped here and processed as a regular word
			if err == nil {
				ss[i] = fmt.Sprintf("%s:* | %d%s:*", s, formed, letter)
				continue
			}
		}
		ss[i] = s + ":*"
	}
	return strings.Join(ss, " & ")
}
