package garbage

import "time"

type EventID string

type Event struct {
	ID               EventID
	Classes          []*Class
	Date             time.Time
	Name             string
	ResourcesAllowed []Resource
}
