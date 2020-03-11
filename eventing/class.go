package eventing

import "github.com/shanvl/garbage-events-service/garbage"

// Class is a model of the class, adapted for this usecase.
type Class struct {
	// Name of the class as it was on the date of the event
	Name string
	// Resources brought by the class to the event
	ResourcesBrought map[garbage.Resource]int
}
