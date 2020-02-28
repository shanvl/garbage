package garbage

import (
	"fmt"
	"time"
)

// EventID uniquely identifies an event
type EventID string

// Event is a gathering of pupils on a specific date, bringing recyclables
type Event struct {
	ID   EventID
	Date time.Time
	Name string
	// Recyclables permitted to be brought to this event
	ResourcesAllowed []Resource
	// Recyclables brought to the event
	ResourcesBrought map[Resource]int
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

// NewEvent returns an instance of a new event. If no name is provided, its date used as the name
func NewEvent(id EventID, date time.Time, name string, resourcesAllowed []Resource) *Event {
	if name == "" {
		year, month, day := date.Date()
		name = fmt.Sprintf("%02d-%02d-%d", day, month, year)
	}
	return &Event{
		ID:               id,
		Date:             date,
		Name:             name,
		ResourcesAllowed: resourcesAllowed,
		ResourcesBrought: make(map[Resource]int),
	}
}
