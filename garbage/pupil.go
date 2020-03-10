package garbage

// PupilID uniquely identifies a pupil
type PupilID string

// Pupil brings recyclable materials (resources) to events.
// This type is often used by various use cases as a carcass for their own Pupil type
type Pupil struct {
	ID        PupilID
	FirstName string
	LastName  string
}
