package eventsvc

import (
	"errors"
	"time"
)

// ErrUnknownPupil is used when the pupil wasn't found
var ErrUnknownEvent = errors.New("unknown event")

// Event is a gathering of pupils who bring in recyclables.
// The goal of the event is to gather as many recyclable materials (resources) as possible
// This type is often used by various use cases as a carcass for their own Event type
type Event struct {
	ID   string
	Date time.Time
	Name string
	// Recyclables permitted to be brought to this event
	ResourcesAllowed []Resource
}

// IsResourcesAllowed checks if a given resource is allowed on this event
func (e *Event) IsResourceAllowed(r Resource) bool {
	for _, rAllowed := range e.ResourcesAllowed {
		if r == rAllowed {
			return true
		}
	}
	return false
}
