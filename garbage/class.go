// Package garbage contains core domain concepts used in the project
package garbage

import (
	"errors"
	"fmt"
	"math"
	"strconv"
	"strings"
	"time"
	"unicode"
)

// ClassID uniquely identifies a particular class
type ClassID string

// Class is a school class consisting of pupils, which changes its name depending on a given date
// relative to the time when it was formed.
// This type is often used by various use cases as a carcass for their own Class type
type Class struct {
	ID ClassID
	// Year in which the class was formed
	// better than an instance of time.Time because that way it has no pointers, so won't be allocated on the heap
	// TODO: no sense in such droching because class is usually passed as a pointer so it will usually be located on the
	//  heap
	YearFormed int
	Letter     string
}

// ErrNoClass is used when a class couldn't be found
var ErrNoClass = errors.New("class doesn't exists")

// NameOnDate constructs a class name on a specific date. For example, if a class, which has
// a letter Б, was formed on 09.2001, on 09.2002 it was 2Б, on 02.2002 it still was 2Б, on 09.2003
// it was 3Б
func (c *Class) NameOnDate(date time.Time) (string, error) {
	// classes are formed on 1st September
	formedDate := time.Date(c.YearFormed, 9, 1, 0, 0, 0, 0, time.UTC)
	yearsPassed := date.Sub(formedDate).Hours() / 24 / 365
	classNumber := int(math.Ceil(yearsPassed))
	if classNumber <= 0 || classNumber > 11 {
		return "", errors.New("invalid date range")
	}
	return fmt.Sprintf("%d%s", classNumber, c.Letter), nil
}

// ParseClassName derives a class' letter and the year the class was formed from its name and a given date.
// Let the className be "3B" and the date is 10.10.2010. Then the letter is "B" and the class was formed in 2008
func ParseClassName(className string, date time.Time) (letter string, yearFormed int, err error) {
	// string to be Atoi'ed to the class number
	numberBuf := strings.Builder{}
	// parse the className, ignoring non-alphanumeric chars;
	// if there are two letters or a digit after a letter, throw an error
	wasLetter := false
	for _, r := range className {
		if wasLetter && (unicode.IsLetter(r) || unicode.IsNumber(r)) {
			return "", 0, fmt.Errorf("invalid class: %s", className)
		}
		if unicode.IsNumber(r) {
			numberBuf.WriteRune(r)
		}
		if unicode.IsLetter(r) {
			letter = string(unicode.ToUpper(r))
			wasLetter = true
		}
	}
	number, err := strconv.Atoi(numberBuf.String())
	if err != nil {
		return "", 0, err
	}
	if number > 11 || number < 1 {
		return "", 0, fmt.Errorf("invalid class number: %d", number)
	}
	yearFormed = date.Year() - number
	// classes are yearFormed in September
	if date.Month() >= 9 {
		yearFormed += 1
	}
	return letter, yearFormed, nil
}
