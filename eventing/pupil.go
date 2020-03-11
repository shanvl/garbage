package eventing

import "github.com/shanvl/garbage-events-service/garbage"

// Pupil is a model of the pupil, adapted for this use case.
type Pupil struct {
	garbage.Pupil
	// Note that Class here is a string, not a garbage.Class instance.
	// It is the name of the class as it was on the date of the event
	Class string
	// Resources brought by the pupil to the event
	ResourcesBrought map[garbage.Resource]int
}
