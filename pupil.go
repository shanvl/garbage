package garbage

// PupilID uniquely identifies a pupil
type PupilID string

// Pupil brings recyclables to events
type Pupil struct {
	ID PupilID
	// A pupil can be moved between classes at any given moment
	Class     *Class
	FirstName string
	LastName  string
}
