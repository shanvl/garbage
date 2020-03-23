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
	// Classes returns a list of sorted classes with resources they brought to events that passed given filters
	Classes(ctx context.Context, filters ClassesFilters, classesSorting, eventsSorting sorting.By, amount, skip int) (classes []*Class,
		total int, err error)
	// ClassByID(ctx context.Context, filters ClassByIDFilters) (*Class, error)
}

// Repository provides methods to work with entities persistence
type Repository interface {
	Classes(ctx context.Context, filters ClassesFilters, classesSorting, eventsSorting sorting.By, amount,
		skip int) (classes []*Class, total int, err error)
	// ClassByID(ctx context.Context, filters ClassByIDFilters) (*Class, error)
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
	// events in every class found can be sorted by resources, name or date
	if !eventsSorting.IsResources() && !eventsSorting.IsName() && !eventsSorting.IsDate() {
		eventsSorting = sorting.DateDes
	}
	return s.repo.Classes(ctx, filters, classesSorting, eventsSorting, amount, skip)
}

// func (s *service) ClassByID(ctx context.Context, filters ClassByIDFilters) (*Class, error) {
//
// }

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

type ClassByIDFilters struct {
	EventsByDateFilter
}

// EventsByDateFilters are used to filter events by date
type EventsByDateFilter struct {
	// include events occurred since this date
	From time.Time
	// include events occurred up to this date
	To time.Time
}
