package eventing

import (
	"context"
	"fmt"
	"time"

	"github.com/shanvl/garbage-events-service/internal/garbage"
	"github.com/shanvl/garbage-events-service/internal/idgen"
	"github.com/shanvl/garbage-events-service/internal/sorting"
	"github.com/shanvl/garbage-events-service/internal/valid"
)

// Service is an interface providing methods to manage an event.
// Note that all methods and entities are used in the context of one event.
type Service interface {
	// ClassByID returns a class with a given id w/ resources for the specified event
	ClassByID(ctx context.Context, classID garbage.ClassID, eventID garbage.EventID) (*Class, error)
	// ChangeEventResources adds/subtracts resources brought by a pupil to/from the event
	ChangeEventResources(ctx context.Context, eventID garbage.EventID, pupilID garbage.PupilID,
		resources map[garbage.Resource]int) (*Event, *Pupil, error)
	// CreateEvent creates and stores an event
	CreateEvent(ctx context.Context, date time.Time, name string, resources []garbage.Resource) (garbage.EventID, error)
	// DeleteEvent deletes an event
	DeleteEvent(ctx context.Context, eventID garbage.EventID) (garbage.EventID, error)
	// EventByID returns an event by its ID
	EventByID(ctx context.Context, eventID garbage.EventID) (*Event, error)
	// EventClasses returns an array of sorted classes for the specified event
	EventClasses(ctx context.Context, eventID garbage.EventID, filters EventClassesFilters, sortBy sorting.By,
		amount, skip int) (classes []*Class, total int, err error)
	// EventPupils returns an array of sorted pupils for the specified event
	EventPupils(ctx context.Context, eventID garbage.EventID, filters EventPupilsFilters, sortBy sorting.By,
		amount int, skip int) (pupils []*Pupil, total int, err error)
	// PupilByID returns a pupil with a given id w/ resources for the specified event
	PupilByID(ctx context.Context, pupilID garbage.PupilID, eventID garbage.EventID) (*Pupil, error)
}

// Repository provides methods to work with an event's persistence
type Repository interface {
	ClassByID(ctx context.Context, classID garbage.ClassID, eventID garbage.EventID) (*Class, error)
	ChangeEventResources(ctx context.Context, eventID garbage.EventID, pupilID garbage.PupilID,
		resources map[garbage.Resource]int) (*Event, *Pupil, error)
	DeleteEvent(ctx context.Context, eventID garbage.EventID) (garbage.EventID, error)
	EventByID(ctx context.Context, eventID garbage.EventID) (*Event, error)
	EventClasses(ctx context.Context, eventID garbage.EventID, filters EventClassesFilters, sortBy sorting.By,
		amount int, skip int) (classes []*Class, total int, err error)
	EventPupils(ctx context.Context, eventID garbage.EventID, filters EventPupilsFilters, sortBy sorting.By,
		amount int, skip int) (pupils []*Pupil, total int, err error)
	PupilByID(ctx context.Context, pupilID garbage.PupilID, eventID garbage.EventID) (*Pupil, error)
	StoreEvent(ctx context.Context, event *garbage.Event) (garbage.EventID, error)
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

// ClassByID returns a class with a given id w/ resources for a specified event
func (s *service) ClassByID(ctx context.Context, classID garbage.ClassID, eventID garbage.EventID) (*Class, error) {
	// check if classID and eventID are provided
	errVld := valid.EmptyError()
	if len(classID) == 0 {
		errVld.Add("classID", "classID must be provided")
	}
	if len(eventID) == 0 {
		errVld.Add("eventID", "eventID must be provided")
	}
	if !errVld.IsEmpty() {
		return nil, errVld
	}

	// get the class
	return s.repo.ClassByID(ctx, classID, eventID)
}

// ChangeEventResources adds/subtracts resources brought by a pupil to/from the event
func (s *service) ChangeEventResources(ctx context.Context, eventID garbage.EventID, pupilID garbage.PupilID,
	resources map[garbage.Resource]int) (*Event, *Pupil, error) {

	errVld := valid.EmptyError()
	if len(pupilID) == 0 {
		errVld.Add("pupilID", "pupilID must be provided")
	}
	if len(eventID) == 0 {
		errVld.Add("eventID", "eventID must be provided")
	}
	if len(resources) == 0 {
		errVld.Add("resources", "no resources were provided")
	}
	if !errVld.IsEmpty() {
		return nil, nil, errVld
	}
	// find an event by its id
	event, err := s.repo.EventByID(ctx, eventID)
	if err != nil {
		return nil, nil, err
	}
	// check that provided resources are allowed at this event
	for res := range resources {
		if !event.IsResourceAllowed(res) {
			errVld.Add("resources", fmt.Sprintf("%s not allowed", res))
			return nil, nil, errVld
		}
	}
	// make changes
	return s.repo.ChangeEventResources(ctx, eventID, pupilID, resources)
}

// CreateEvent creates and stores an event
func (s *service) CreateEvent(ctx context.Context, date time.Time, name string,
	resourcesAllowed []garbage.Resource) (garbage.EventID, error) {

	errVld := valid.EmptyError()
	// new event mustn't occur in the past
	if time.Now().After(date) {
		errVld.Add("date", "event's date must be in the future")
	}
	// check that provided resourcesAllowed exist and are known
	if len(resourcesAllowed) == 0 {
		errVld.Add("resourcesAllowed", "at least one resource must be specified")
	}
	if !errVld.IsEmpty() {
		return "", errVld
	}
	for _, resource := range resourcesAllowed {
		if !resource.IsKnown() {
			errVld.Add("resourcesAllowed", "unknown resource")
			break
		}
	}
	// if there are validation errors, return them
	if !errVld.IsEmpty() {
		return "", errVld
	}
	// If no name was provided, create it from the event's date
	if len(name) == 0 {
		year, month, day := date.Date()
		name = fmt.Sprintf("%02d-%02d-%d", day, month, year)
	}
	// generate eventID
	id, err := idgen.CreateEventID()
	if err != nil {
		return "", err
	}
	// create event
	event := garbage.NewEvent(id, date, name, resourcesAllowed)

	return s.repo.StoreEvent(ctx, event)
}

// DeleteEvent deletes an event
func (s *service) DeleteEvent(ctx context.Context, eventID garbage.EventID) (garbage.EventID, error) {
	errVld := valid.EmptyError()
	// check if there's eventID
	if len(eventID) == 0 {
		errVld.Add("eventID", "eventID must be provided")
	}
	if !errVld.IsEmpty() {
		return "", errVld
	}

	// delete the event
	return s.repo.DeleteEvent(ctx, eventID)
}

// EventByID returns an event by its ID
func (s *service) EventByID(ctx context.Context, eventID garbage.EventID) (*Event, error) {
	errVld := valid.EmptyError()
	if len(eventID) == 0 {
		errVld.Add("eventID", "eventID is needed")
	}
	if !errVld.IsEmpty() {
		return nil, errVld
	}

	return s.repo.EventByID(ctx, eventID)
}

// EventClasses returns an array of sorted classes for the specified event
// TODO: classes names must be relative to the event's date
func (s *service) EventClasses(ctx context.Context, eventID garbage.EventID, filters EventClassesFilters,
	sortBy sorting.By, amount, skip int) (classes []*Class, total int, err error) {

	// check if eventID was provided
	errVld := valid.EmptyError()
	if len(eventID) == 0 {
		errVld.Add("eventID", "eventID must be provided")
	}
	if !errVld.IsEmpty() {
		return nil, 0, errVld
	}

	// if provided values are incorrect, use default values instead
	amount, skip = validateAmountSkip(amount, skip)

	if !sortBy.IsName() && !sortBy.IsResources() {
		sortBy = sorting.NameAsc
	}

	return s.repo.EventClasses(ctx, eventID, filters, sortBy, amount, skip)
}

// EventPupils returns an array of sorted pupils for the specified event
// TODO: classes names must be relative to the event's date
func (s *service) EventPupils(ctx context.Context, eventID garbage.EventID, filters EventPupilsFilters,
	sortBy sorting.By, amount int, skip int) (pupils []*Pupil, total int, err error) {

	// check if eventID was provided
	errVld := valid.EmptyError()
	if len(eventID) == 0 {
		errVld.Add("eventID", "eventID must be provided")
	}
	if !errVld.IsEmpty() {
		return nil, 0, errVld
	}

	// if provided values are incorrect, use default values instead
	amount, skip = validateAmountSkip(amount, skip)

	if !sortBy.IsName() && !sortBy.IsResources() {
		sortBy = sorting.NameAsc
	}

	return s.repo.EventPupils(ctx, eventID, filters, sortBy, amount, skip)
}

// PupilByID returns a pupil with a given id w/ resources for a specified event
// TODO: classes names must be relative to the event's date
func (s *service) PupilByID(ctx context.Context, pupilID garbage.PupilID, eventID garbage.EventID) (*Pupil, error) {
	// check if eventID and pupilID are provided
	errVld := valid.EmptyError()
	if len(eventID) == 0 {
		errVld.Add("eventID", "eventID must be provided")
	}
	if len(pupilID) == 0 {
		errVld.Add("pupilID", "pupilID must be provided")
	}
	if !errVld.IsEmpty() {
		return nil, errVld
	}

	// get the pupil
	return s.repo.PupilByID(ctx, pupilID, eventID)
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

// Class is a model of the class, adapted for this use case.
type Class struct {
	ID garbage.ClassID
	// Name of the class as it was on the date of the event
	Name string
	// Resources brought by the class to the event
	ResourcesBrought map[garbage.Resource]int
}

// Event is a model of the event, adapted for this use case.
// It indicates how many resources have been collected at this event
type Event struct {
	garbage.Event
	// resources brought by pupils to this event
	ResourcesBrought map[garbage.Resource]int
}

// Pupil is a model of the pupil, adapted for this use case.
type Pupil struct {
	garbage.Pupil
	// Note that Class here is a string, not a garbage.Class instance.
	// It is the name of the class as it was on the date of the event
	Class string
	// Resources brought by the pupil to the event
	ResourcesBrought map[garbage.Resource]int
}

// EventClassesFilters are used to filter classes participating in an event
type EventClassesFilters struct {
	Name string
}

// EventPupilsFilters are used to filter pupils participating in an event
type EventPupilsFilters struct {
	Name string
}
