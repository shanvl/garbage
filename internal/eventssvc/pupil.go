package eventssvc

import "errors"

// ErrUnknownPupil is used when the pupil wasn't found
var ErrUnknownPupil = errors.New("unknown pupil")

// Pupil brings recyclable materials (resources) to events.
// This type is often used by various use cases as a carcass for their own Pupil type
type Pupil struct {
	ID        string
	FirstName string
	LastName  string
}
