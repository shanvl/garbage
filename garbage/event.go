package garbage

import (
	"time"
)

// EventID uniquely identifies an event
type EventID string

// Event is a meeting of pupils who bring in recyclables.
// The goal of the event is to gather as many recyclable materials (resources) as possible
type Event struct {
	ID   EventID
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

// NewEvent returns an instance of an event
func NewEvent(id EventID, date time.Time, name string, resourcesAllowed []Resource) *Event {
	return &Event{
		ID:               id,
		Date:             date,
		Name:             name,
		ResourcesAllowed: resourcesAllowed,
	}
}
