package garbage

// PupilID uniquely identifies a pupil
type PupilID string

// Pupil brings recyclable materials (resources) to events
type Pupil struct {
	ID        PupilID
	FirstName string
	LastName  string
}
