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
	// Recyclables brought by a pupil to selected events
	ResourcesBrought map[Resource]int
}

// NewPupil returns an instance of a new pupil
func NewPupil(id PupilID, class *Class, firstName string, lastName string) *Pupil {
	return &Pupil{
		ID:               id,
		Class:            class,
		FirstName:        firstName,
		LastName:         lastName,
		ResourcesBrought: make(map[Resource]int),
	}
}
