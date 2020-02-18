// Package garbage contains core concepts used by the project's services
package garbage

import (
	"errors"
	"fmt"
	"math"
	"time"
)

// ClassID uniquely identifies a particular class
type ClassID string

// Class is a school class consisting of pupils, which changes its name depending on a given date
// relative to the time when it was formed
type Class struct {
	ID ClassID
	// Date on which the class was formed
	Formed time.Time
	Letter string
}

// NameFromDate constructs a class name for a specified date. For example, if a class, which has
// a letter Б, was formed on 09.2001, on 09.2002 it was 2Б, on 02.2002 it still was 2Б, on 09.2003
// it was 3Б
func (c *Class) NameFromDate(date time.Time) (string, error) {
	yearsPassed := date.Sub(c.Formed).Hours() / 24 / 365
	classNumber := int(math.Ceil(yearsPassed))
	if classNumber <= 0 || classNumber > 11 {
		return "", errors.New("invalid date range")
	}
	return fmt.Sprintf("%d%s", classNumber, c.Letter), nil
}
