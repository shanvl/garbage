package garbage

import (
	"errors"
	"fmt"
	"math"
	"time"
)

type ClassID string

type Class struct {
	ID     ClassID
	Formed time.Time
	Letter string
}

// construct a class name for a specified date. For example, if a class, which has a letter Б,
// was formed on 09.2001, on 09.2002 it was 2Б, on 02.2002 it still was 2Б, on 09.2003 it was 3Б
func (c *Class) NameFromDate(date time.Time) (string, error) {
	yearsPassed := date.Sub(c.Formed).Hours() / 24 / 365
	classNumber := int(math.Ceil(yearsPassed))
	if classNumber <= 0 || classNumber > 11 {
		return "", errors.New("invalid date range")
	}
	return fmt.Sprintf("%d%s", classNumber, c.Letter), nil
}
