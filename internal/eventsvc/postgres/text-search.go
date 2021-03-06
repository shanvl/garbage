package postgres

import (
	"fmt"
	"strings"
	"time"
	"unicode"

	"github.com/shanvl/garbage/internal/eventsvc"
	pgtextsearch "github.com/shanvl/garbage/pkg/pg-text-search"
)

// prepareTextSearchClass processes the query to make it a valid argument for to_tsquery,
// adding ':*' at the end of each word, concatenating the words with ' & ' and,
// if a word resembles a school class name,
// creating a copy of it with changes needed to hit the indices of the tables. For example,
// if the passed date is 10.10.2020, a query "3A Iv Ig" will be changed to "(3A:* | 2018A:*) & Iv:* & Ig:*".
// If the input contains invalid symbols, it simply returns an empty string
func prepareTextSearchClass(q string, t time.Time) string {
	ss := strings.Fields(q)
	for i, s := range ss {
		if !pgtextsearch.IsValidInput(s) {
			return ""
		}
		// if the word starts with a number, try to parse it as a class name and concatenate with itself
		if unicode.IsDigit(rune(s[0])) {
			// if the word can be parsed as a class name, process it and concatenate with itself.
			// "3B" becomes "3B:* | 2018A:*"
			letter, dateFormed, err := eventsvc.ParseClassName(s, t)
			// error means that neither letter, nor dateFormed can be extracted from the word,
			// so it shouldn't be processed as a class
			if err == nil {
				if len(letter) > 0 && !dateFormed.IsZero() {
					ss[i] = fmt.Sprintf("(%s:* | %d%s:*)", s, dateFormed.Year(), letter)
				} else if !dateFormed.IsZero() {
					ss[i] = fmt.Sprintf("(%s:* | %d:*)", s, dateFormed.Year())
				}
				continue
			}
		}
		ss[i] = s + ":*"
	}
	return strings.Join(ss, " & ")
}
