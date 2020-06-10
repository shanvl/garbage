// Package eventssvc contains the core domain concepts used in the project
package eventssvc

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

// ErrNoClassOnDate indicates that there was no class with the specified credentials on the given date.
// It can occur, for example, if the class was formed in 2002 and the date is 01.01.2020
var ErrNoClassOnDate = errors.New("there was no such class on that date")

// ErrInvalidClassName is used when something wrong with extracting a class info from a string
var ErrInvalidClassName = errors.New("invalid class name")

// Class is a school class consisting of pupils, which changes its name depending on a given date
// relative to the time when it was formed.
// This type is often used by various use cases as a carcass for their own Class types
type Class struct {
	Letter     string
	DateFormed time.Time
}

// NameOnDate constructs a class name on the specified date. For example, if a class with a letter B
// was formed on 09.2001, on 09.2002 it was 2B, on 02.2002 it still was 2B, on 09.2003 it became 3B
func (c Class) NameOnDate(date time.Time) (string, error) {
	// classes are formed on 1st September
	yearsPassed := date.Sub(c.DateFormed).Hours() / 24 / 365
	classNumber := int(math.Ceil(yearsPassed))
	if classNumber <= 0 || classNumber > 11 {
		return "", ErrNoClassOnDate
	}
	return fmt.Sprintf("%d%s", classNumber, c.Letter), nil
}

// ClassFromClassName derives a class' instance from its name and the given date.
// Let the className be "3B" and the date is 10.10.2010. Then the letter is "B" and the class was formed in 2008
func ClassFromClassName(className string, date time.Time) (Class, error) {
	letter, dateFormed, err := ParseClassName(className, date)
	if err != nil {
		return Class{}, err
	}
	if len(letter) == 0 {
		return Class{}, fmt.Errorf("%w: class must have a letter", ErrInvalidClassName)
	}
	if dateFormed.IsZero() {
		return Class{}, fmt.Errorf("%w: class must have a number", ErrInvalidClassName)
	}
	return Class{letter, dateFormed}, nil
}

// ParseClassName extracts the letter and the date the class was formed from the given class name and date.
// It's ok if one of them is not present
func ParseClassName(className string, date time.Time) (letter string, dateFormed time.Time, err error) {
	if len(className) == 0 {
		return "", time.Time{}, fmt.Errorf("%w: className can't be empty", ErrInvalidClassName)
	}
	// string to be Atoi'ed to the class number
	numberBuf := strings.Builder{}
	// parse the className, ignoring non-alphanumeric chars;
	// if there are two letters or a digit after a letter, throw an error
	wasLetter := false
	for _, r := range className {
		if wasLetter && (unicode.IsLetter(r) || unicode.IsNumber(r)) {
			return "", time.Time{}, fmt.Errorf("%w: %s", ErrInvalidClassName, className)
		}
		if unicode.IsNumber(r) {
			numberBuf.WriteRune(r)
		}
		if unicode.IsLetter(r) {
			letter = string(unicode.ToLower(r))
			wasLetter = true
		}
	}
	if numberBuf.Len() > 0 {
		number, err := strconv.Atoi(numberBuf.String())
		if err != nil {
			return "", time.Time{}, err
		}
		if number > 11 || number < 1 {
			return "", time.Time{}, fmt.Errorf("%w: invalid number: %d", ErrInvalidClassName, number)
		}
		yearFormed := date.Year() - number
		// classes are formed in September
		if date.Month() >= 9 {
			yearFormed += 1
		}
		dateFormed = time.Date(yearFormed, 9, 1, 0, 0, 0, 0, time.UTC)
	}
	return letter, dateFormed, nil
}
