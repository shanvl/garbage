package eventing

import "time"

// Filters are used to narrow down a list when searching for events
type Filters struct {
	// event name
	Name string
	// start date of events search
	From time.Time
	// end date of events search
	To time.Time
}
