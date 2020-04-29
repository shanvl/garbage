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

// ErrNoClass is used when a class couldn't be found
var ErrNoPupil = errors.New("pupil doesn't exists")
