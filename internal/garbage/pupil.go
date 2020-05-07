package garbage

import "errors"

// PupilID uniquely identifies a pupil
type PupilID string

// Pupil brings recyclable materials (resources) to events.
// This type is often used by various use cases as a carcass for their own Pupil type
type Pupil struct {
	ID        PupilID
	FirstName string
	LastName  string
}

// ErrNoPupil is used when the pupil wasn't found
var ErrNoPupil = errors.New("pupil wasn't found")
