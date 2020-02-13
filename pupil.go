package garbage

type PupilID string

type Pupil struct {
	ID        PupilID
	Class     *Class
	FirstName string
	LastName  string
}
