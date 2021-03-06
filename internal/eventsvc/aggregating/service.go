// Package aggregating is responsible for providing aggregated info on how classes and pupils performed at
// events, with various filters applied
package aggregating

import (
	"context"
	"fmt"
	"time"

	"github.com/shanvl/garbage/internal/eventsvc"
	"github.com/shanvl/garbage/internal/eventsvc/sorting"
	"github.com/shanvl/garbage/pkg/valid"
)

// Service is an interface providing methods for obtaining aggregated info on how classes and pupils performed at
// events, with various filters applied
type Service interface {
	// Classes returns a list of sorted classes, each of which has a list of events that passed the given filters
	Classes(ctx context.Context, filters ClassFilters, classesSorting, eventsSorting sorting.By, amount,
		skip int) (classes []*Class, total int, err error)
	// Events returns a list of sorted events that passed the provided filters
	Events(ctx context.Context, filters EventFilters, sortBy sorting.By, amount, skip int) (events []*Event,
		total int, err error)
	// Pupils returns a list of sorted classes, each of which has a list of events that passed the given filters
	Pupils(ctx context.Context, filters PupilFilters, pupilsSorting, eventsSorting sorting.By, amount,
		skip int) (pupils []*Pupil, total int, err error)
	// PupilByID returns a pupil with the given ID and a list of events they has attended
	PupilByID(ctx context.Context, id string, filters EventFilters, eventsSorting sorting.By) (*Pupil,
		error)
}

// Repository provides methods to work with entities persistence
type Repository interface {
	Classes(ctx context.Context, filters ClassFilters, classesSorting, eventsSorting sorting.By, amount,
		skip int) (classes []*Class, total int, err error)
	Pupils(ctx context.Context, filters PupilFilters, pupilsSorting, eventsSorting sorting.By, amount,
		skip int) (pupils []*Pupil, total int, err error)
	PupilByID(ctx context.Context, id string, filters EventFilters, eventsSorting sorting.By) (*Pupil,
		error)
	Events(ctx context.Context, filters EventFilters, sortBy sorting.By, amount, skip int) (events []*Event,
		total int, err error)
}

type service struct {
	repo Repository
}

const (
	DefaultAmount = 50
	DefaultSkip   = 0
	MaxAmount     = 1000
)

// NewService returns an instance of Service with all its dependencies
func NewService(repo Repository) Service {
	return &service{repo}
}

// Classes returns a list of sorted classes, each of which has a list of events that passed the given filters
func (s *service) Classes(ctx context.Context, filters ClassFilters, classesSorting, eventsSorting sorting.By,
	amount, skip int) (classes []*Class, total int, err error) {

	// there should be no class letter at all or only one letter
	if len(filters.Letter) > 1 {
		err := valid.NewError("class letter", fmt.Sprintf("invalid class letter: %s", filters.Letter))
		return nil, 0, err
	}
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

// Events returns a list of sorted events that passed the provided filters
func (s *service) Events(ctx context.Context, filters EventFilters, sortBy sorting.By, amount,
	skip int) (events []*Event, total int, err error) {

	// if provided values are incorrect, use default ones instead
	amount, skip = validateAmountSkip(amount, skip)
	// if eventsSorting is invalid, use default one instead
	sortBy = validateEventsSorting(sortBy)

	return s.repo.Events(ctx, filters, sortBy, amount, skip)
}

// Pupils returns a list of sorted classes, each of which has a list of events that passed the given filters
func (s *service) Pupils(ctx context.Context, filters PupilFilters, pupilsSorting, eventsSorting sorting.By, amount,
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

// PupilByID returns a pupil with the given ID and a list of events they has attended
func (s *service) PupilByID(ctx context.Context, id string, filters EventFilters,
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

// if the events sorting passed is not set to resources, name or date, sets it to DateDes
func validateEventsSorting(s sorting.By) sorting.By {
	if !s.IsResources() && !s.IsName() && !s.IsDate() {
		s = sorting.DateDes
	}
	return s
}

// ensures that amount and skip are valid
func validateAmountSkip(a, s int) (int, int) {
	if a <= 0 || a > MaxAmount {
		a = DefaultAmount
	}
	if s < 0 {
		s = DefaultSkip
	}
	return a, s
}

// Class is a model of the class, adapted for this use case
type Class struct {
	eventsvc.Class
	// all the resources the class brought to the events
	ResourcesBrought eventsvc.ResourceMap
	// list of events with resources brought by the class to each of them
	Events []*Event
}

// Event is a model of the event, adapted for this use case.
type Event struct {
	eventsvc.Event
	// resources collected at this event OR resources brought by the parent entity to this event
	ResourcesBrought eventsvc.ResourceMap
}

// Pupil is a model of the pupil, adapted for this use case
type Pupil struct {
	eventsvc.Pupil
	eventsvc.Class
	// all the resources the pupil brought to the events
	ResourcesBrought eventsvc.ResourceMap
	// list of events with resources brought by the pupil to each of them
	Events []*Event
}

// ClassFilters are used to filter classes and events in which they participated
type ClassFilters struct {
	EventFilters
	// Letter of the class
	Letter string
	// Date the class was formed at
	DateFormed time.Time
}

// EventFilters are used to filter events
type EventFilters struct {
	// include events occurred since this date
	From time.Time
	// include events occurred up to this date
	To time.Time
	// Name of the event
	Name string
	// Recyclables permitted to be brought to this event
	ResourcesAllowed []eventsvc.Resource
}

// PupilFilters are used to filter pupils and events in which they participated
type PupilFilters struct {
	EventFilters
	NameAndClass string
}
