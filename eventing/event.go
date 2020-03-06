package eventing

import "github.com/shanvl/garbage-events-service/garbage"

// Event is a model of the event, adapted for this usecase.
// It indicates how many resources have been collected at this event
type Event struct {
	garbage.Event
	// resources brought by pupils to this event
	ResourcesBrought map[garbage.Resource]int
}
