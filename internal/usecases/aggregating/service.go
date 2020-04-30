// Package aggregating is responsible for providing aggregated info on how classes and pupils performed at
// events scattered in time, with various filters applied
package aggregating

import (
	"context"
	"time"

	"github.com/shanvl/garbage-events-service/internal/garbage"
	"github.com/shanvl/garbage-events-service/internal/sorting"
	"github.com/shanvl/garbage-events-service/internal/valid"
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
	// Pupils returns a list of sorted pupils with a list of resources they brought to events that passed the given
	// filters
	Pupils(ctx context.Context, filters PupilsFilters, pupilsSorting, eventsSorting sorting.By, amount,
		skip int) (pupils []*Pupil, total int, err error)
	// PupilByID returns a pupil with the given ID with a list of all the resources they has brought to every event that
	// passed the provided filter
	PupilByID(ctx context.Context, id garbage.PupilID, filters EventsByDateFilter, eventsSorting sorting.By) (*Pupil,
		error)
	// Events returns a list of sorted events that passed the provided filters
	Events(ctx context.Context, filters EventsFilters, sortBy sorting.By, amount, skip int) (events []*Event,
		total int, err error)
}

// Repository provides methods to work with entities persistence
type Repository interface {
	Classes(ctx context.Context, filters ClassesFilters, classesSorting, eventsSorting sorting.By, amount,
		skip int) (classes []*Class, total int, err error)
	ClassByID(ctx context.Context, id garbage.ClassID, filters EventsByDateFilter, eventsSorting sorting.By) (*Class,
		error)
	Pupils(ctx context.Context, filters PupilsFilters, pupilsSorting, eventsSorting sorting.By, amount,
		skip int) (pupils []*Pupil, total int, err error)
	PupilByID(ctx context.Context, id garbage.PupilID, filters EventsByDateFilter, eventsSorting sorting.By) (*Pupil,
		error)
	Events(ctx context.Context, filters EventsFilters, sortBy sorting.By, amount, skip int) (events []*Event,
		total int, err error)
}

type service struct {
	repo Repository
}

const (
	DefaultAmount = 25
	DefaultSkip   = 0
)

// NewService returns an instance of Service with all its dependencies
func NewService(repo Repository) Service {
	return &service{repo}
}

// Classes returns a list of sorted classes with resources they brought to events that passed given filters
func (s *service) Classes(ctx context.Context, filters ClassesFilters, classesSorting, eventsSorting sorting.By,
	amount, skip int) (classes []*Class, total int, err error) {

	// if provided values are incorrect, use default ones instead
	amount, skip = validateAmountSkip(amount, skip)
	// classes can be sorted by resources they brought or by name
	if !classesSorting.IsResources() && !classesSorting.IsName() {
		classesSorting = sorting.NameAsc
	}
	// if eventsSorting is invalid, use default one instead
	eventsSorting = validateEventsSorting(eventsSorting)

	return s.repo.Classes(ctx, filters, classesSorting, eventsSorting, amount, skip)
}

// ClassByID returns a class with the given ID with a list of all the resources it has brought to every event that
// passed the provided filter. Events are sorted
func (s *service) ClassByID(ctx context.Context, id garbage.ClassID, filters EventsByDateFilter,
	eventsSorting sorting.By) (*Class, error) {

	// check if eventID is provided
	errVld := valid.EmptyError()
	if len(id) == 0 {
		errVld.Add("eventID", "eventID must be provided")
		return nil, errVld
	}
	// validate events sorting
	eventsSorting = validateEventsSorting(eventsSorting)

	return s.repo.ClassByID(ctx, id, filters, eventsSorting)
}

// Events returns a list of sorted events that passed the provided filters
func (s *service) Events(ctx context.Context, filters EventsFilters, sortBy sorting.By, amount,
	skip int) (events []*Event, total int, err error) {

	// if provided values are incorrect, use default ones instead
	amount, skip = validateAmountSkip(amount, skip)
	// if eventsSorting is invalid, use default one instead
	sortBy = validateEventsSorting(sortBy)

	return s.repo.Events(ctx, filters, sortBy, amount, skip)
}

// Pupils returns a list of sorted pupils with a list of resources they brought to events that passed the given
// filters
func (s *service) Pupils(ctx context.Context, filters PupilsFilters, pupilsSorting, eventsSorting sorting.By, amount,
	skip int) (pupils []*Pupil, total int, err error) {

	// if provided values are incorrect, use default ones instead
	amount, skip = validateAmountSkip(amount, skip)

	// pupils can be sorted by resources they brought or by name
	if !pupilsSorting.IsName() && !pupilsSorting.IsResources() {
		pupilsSorting = sorting.NameAsc
	}
	// if eventsSorting is invalid, use default one instead
	eventsSorting = validateEventsSorting(eventsSorting)

	return s.repo.Pupils(ctx, filters, pupilsSorting, eventsSorting, amount, skip)
}

// PupilByID returns a pupil with the given ID with a list of all the resources they has brought to every event that
// passed the provided filter. Events are sorted
func (s *service) PupilByID(ctx context.Context, id garbage.PupilID, filters EventsByDateFilter,
	eventsSorting sorting.By) (*Pupil, error) {

	// check if pupilID is provided
	errVld := valid.EmptyError()
	if len(id) == 0 {
		errVld.Add("pupilID", "pupilID must be provided")
		return nil, errVld
	}

	// if eventsSorting is invalid, use default one instead
	eventsSorting = validateEventsSorting(eventsSorting)

	return s.repo.PupilByID(ctx, id, filters, eventsSorting)
}

// if events sorting passed is not set to resources, name or date, sets it to DateDes
func validateEventsSorting(s sorting.By) sorting.By {
	if !s.IsResources() && !s.IsName() && !s.IsDate() {
		s = sorting.DateDes
	}
	return s
}

// ensures that amount and skip are valid
func validateAmountSkip(a, s int) (int, int) {
	if a <= 0 {
		a = DefaultAmount
	}
	if s < 0 {
		s = DefaultSkip
	}
	return a, s
}

// Class is a model of the class, adapted for this use case
type Class struct {
	garbage.Class
	// list of events with resources brought by the class to each of them
	Events []Event
}

// Event is a model of the event, adapted for this use case.
type Event struct {
	garbage.Event
	// resources collected at this event OR resources brought by the parent entity to this event
	ResourcesBrought map[garbage.Resource]int
}

// Pupil is a model of the pupil, adapted for this use case
type Pupil struct {
	garbage.Pupil
	// list of events with resources brought by the pupil to each of them
	Events []Event
}

// ClassesFilters are used to filter classes and events in which they participated
type ClassesFilters struct {
	EventsByDateFilter
	// Letter of the class
	Letter string
	// Year the class was formed in
	YearFormed string
}

// EventsFilters are used to filter events
type EventsFilters struct {
	EventsByDateFilter
	// Name of the event
	Name string
	// Recyclables permitted to be brought to this event
	ResourcesAllowed []garbage.Resource
}

// PupilsFilters are used to filter pupils and events in which they participated
type PupilsFilters struct {
	EventsByDateFilter
	// Name of the pupil
	Name string
}

// EventsByDateFilter is used to filter events by date
type EventsByDateFilter struct {
	// include events occurred since this date
	From time.Time
	// include events occurred up to this date
	To time.Time
}
