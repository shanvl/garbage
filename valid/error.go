// Package valid provides tools to validate some input and return a convenient error
package valid

import (
	"fmt"
	"strings"
)

// ErrValidation is a validation error which can store several errors as its fields
type ErrValidation struct {
	fields map[string]string
}

// Add simply adds an error as a field and description to the internal map
// which then can used by other methods to obtain a nice human-readable error
func (e *ErrValidation) Add(field, errDesc string) {
	e.fields[field] = errDesc
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

// IsEmpty shows whether the error is empty, i.e. has no errors added to it
func (e *ErrValidation) IsEmpty() bool {
	if len(e.fields) == 0 {
		return true
	}
	return false
}

// EmptyError returns a new empty ErrValidation, so that it can be populated with errors later
func EmptyError() *ErrValidation {
	return &ErrValidation{
		fields: make(map[string]string),
	}
}
