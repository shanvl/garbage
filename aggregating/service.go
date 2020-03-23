// Package aggregating is responsible for providing aggregated info on how classes and pupils performed at
// events scattered in time, with various filters applied
package aggregating

import (
	"context"
	"time"

	"github.com/shanvl/garbage-events-service/garbage"
	"github.com/shanvl/garbage-events-service/sorting"
)

// Service is an interface providing methods for obtaining aggregated info on how classes and pupils performed at
// events scattered in time, with various filters applied
type Service interface {
	// Classes returns a list of sorted classes with a list of resources they brought to events that passed the given
	// filters
	Classes(ctx context.Context, filters ClassesFilters, classesSorting, eventsSorting sorting.By, amount,
		skip int) (classes []*Class, total int, err error)
	// ClassByID returns a class with the given ID with a list of all the resources it has brought to every event that
	// passed the provided filter
	ClassByID(ctx context.Context, id garbage.ClassID, filters EventsByDateFilter, eventsSorting sorting.By) (*Class,
		error)
}

// Repository provides methods to work with entities persistence
type Repository interface {
	Classes(ctx context.Context, filters ClassesFilters, classesSorting, eventsSorting sorting.By, amount,
		skip int) (classes []*Class, total int, err error)
	ClassByID(ctx context.Context, id garbage.ClassID, filters EventsByDateFilter, eventsSorting sorting.By) (*Class,
		error)
}

type service struct {
	repo Repository
}

// NewService returns an instance of Service with all its dependencies
func NewService(repo Repository) Service {
	return &service{repo}
}

const (
	DefaultAmount = 25
	DefaultSkip   = 0
)

// Classes returns a list of sorted classes with resources they brought to events that passed given filters
func (s *service) Classes(ctx context.Context, filters ClassesFilters, classesSorting, eventsSorting sorting.By,
	amount, skip int) (classes []*Class, total int, err error) {

	// if provided values are incorrect, use default ones instead
	if amount <= 0 {
		amount = DefaultAmount
	}
	if skip < 0 {
		skip = DefaultSkip
	}
	// classes can be sorted by resources they brought or by name
	if !classesSorting.IsResources() && !classesSorting.IsName() {
		classesSorting = sorting.NameAsc
	}
	// validate events sorting
	eventsSorting = validateEventsSorting(eventsSorting)

	return s.repo.Classes(ctx, filters, classesSorting, eventsSorting, amount, skip)
}

// ClassByID returns a class with the given ID with a list of all the resources it has brought to every event that
// passed the provided filter. Events are sorted
func (s *service) ClassByID(ctx context.Context, id garbage.ClassID, filters EventsByDateFilter,
	eventsSorting sorting.By) (*Class, error) {
	// validate events sorting
	eventsSorting = validateEventsSorting(eventsSorting)

	return s.repo.ClassByID(ctx, id, filters, eventsSorting)
}

// if events sorting passed is not set to resources, name or date, sets it to DateDes
func validateEventsSorting(s sorting.By) sorting.By {
	if !s.IsResources() && !s.IsName() && !s.IsDate() {
		s = sorting.DateDes
	}
	return s
}

// Class is a model of the class, adapted for this use case
type Class struct {
	garbage.Class
	// list of events with resources brought by the class to them
	Events []Event
}

// Event is a model of the event, adapted for this use case.
type Event struct {
	garbage.Event
	// resources collected at this event OR resources brought by the parent entity to this event
	ResourcesBrought map[garbage.Resource]int
}

// ClassesFilters are used to filter classes and events in which they participated
type ClassesFilters struct {
	EventsByDateFilter
	// Letter of the class
	Letter string
	// Year the class was formed in
	YearFormed string
}

// EventsByDateFilter is used to filter events by date
type EventsByDateFilter struct {
	// include events occurred since this date
	From time.Time
	// include events occurred up to this date
	To time.Time
}
