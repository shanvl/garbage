package valid

import (
	"fmt"
	"strings"
)

// ErrValidation is a valid error with some handy methods
type ErrValidation struct {
	fields map[string]string
}

// Error pretty prints all collected errors
func (e *ErrValidation) Error() string {
	var builder strings.Builder
	for field, err := range e.fields {
		fmt.Fprintf(&builder, "%s: %s\n", field, err)
	}
	return builder.String()
}

// Fields returns a map of errors collected
func (e *ErrValidation) Fields() map[string]string {
	return e.fields
}

// add simply adds an error field and its description to the internal map
// which then is used by some methods to create a nice human-readable error
func (e *ErrValidation) add(field, err string) {
	if e.fields == nil {
		e.fields = make(map[string]string, 10)
	}
	e.fields[field] = err
}
