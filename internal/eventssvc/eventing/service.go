package eventing

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/shanvl/garbage/internal/eventssvc"
	"github.com/shanvl/garbage/internal/eventssvc/idgen"
	"github.com/shanvl/garbage/internal/eventssvc/sorting"
	"github.com/shanvl/garbage/pkg/valid"
)

// Service is an interface providing methods to manage an event.
// Note that all methods and entities are used in the context of one event.
type Service interface {
	// ChangePupilResources adds/subtracts resources brought by a pupil to/from the event
	ChangePupilResources(ctx context.Context, eventID eventssvc.EventID, pupilID eventssvc.PupilID,
		resources eventssvc.ResourceMap) error
	// CreateEvent creates and stores an event
	CreateEvent(ctx context.Context, date time.Time, name string, resources []eventssvc.Resource) (eventssvc.EventID, error)
	// DeleteEvent deletes an event with the id passed
	DeleteEvent(ctx context.Context, eventID eventssvc.EventID) error
	// EventByID returns an event by its ID
	EventByID(ctx context.Context, eventID eventssvc.EventID) (*Event, error)
	// EventClasses returns an array of sorted classes for the specified event
	EventClasses(ctx context.Context, eventID eventssvc.EventID, filters EventClassFilters, sortBy sorting.By,
		amount, skip int) (classes []*Class, total int, err error)
	// EventPupils returns an array of sorted pupils for the specified event
	EventPupils(ctx context.Context, eventID eventssvc.EventID, filters EventPupilFilters, sortBy sorting.By,
		amount int, skip int) (pupils []*Pupil, total int, err error)
	// PupilByID returns a pupil with a given id w/ resources for the specified event
	PupilByID(ctx context.Context, pupilID eventssvc.PupilID, eventID eventssvc.EventID) (*Pupil, error)
}

// Repository provides methods to work with an event's persistence
type Repository interface {
	ChangePupilResources(ctx context.Context, eventID eventssvc.EventID, pupilID eventssvc.PupilID,
		resources eventssvc.ResourceMap) error
	DeleteEvent(ctx context.Context, eventID eventssvc.EventID) error
	EventByID(ctx context.Context, eventID eventssvc.EventID) (*Event, error)
	EventClasses(ctx context.Context, eventID eventssvc.EventID, filters EventClassFilters, sortBy sorting.By,
		amount int, skip int) (classes []*Class, total int, err error)
	EventPupils(ctx context.Context, eventID eventssvc.EventID, filters EventPupilFilters, sortBy sorting.By,
		amount int, skip int) (pupils []*Pupil, total int, err error)
	PupilByID(ctx context.Context, pupilID eventssvc.PupilID, eventID eventssvc.EventID) (*Pupil, error)
	StoreEvent(ctx context.Context, event *eventssvc.Event) (eventssvc.EventID, error)
}

type service struct {
	repo Repository
}

const (
	DefaultAmount = 50
	DefaultSkip   = 0
	MaxAmount     = 1000
)

// ErrNoEventPupil indicates that the pupil didn't participate in the event
var ErrNoEventPupil = errors.New("pupil didn't participate in the event")

// NewService returns an instance of Service with all its dependencies
func NewService(repo Repository) Service {
	return &service{repo}
}

// ChangePupilResources adds/subtracts resources brought by a pupil to/from the event
func (s *service) ChangePupilResources(ctx context.Context, eventID eventssvc.EventID, pupilID eventssvc.PupilID,
	resources eventssvc.ResourceMap) error {

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
		return errVld
	}
	// find an event by its id
	event, err := s.repo.EventByID(ctx, eventID)
	if err != nil {
		return err
	}
	// check that provided resources are allowed at this event
	for res := range resources {
		if !event.IsResourceAllowed(res) {
			return valid.NewError("resources", fmt.Sprintf("%s is not allowed", res))
		}
	}
	err = s.repo.ChangePupilResources(ctx, eventID, pupilID, resources)
	if errors.Is(err, ErrNoEventPupil) {
		// we already know for sure that the event exists —— we've checked it earlier. Hence,
		// we can be certain that only the pupil hasn't been found
		err = eventssvc.ErrUnknownPupil
	}
	return err
}

// CreateEvent creates and stores an event
func (s *service) CreateEvent(ctx context.Context, date time.Time, name string,
	resourcesAllowed []eventssvc.Resource) (eventssvc.EventID, error) {

	errVld := valid.EmptyError()
	// new event mustn't occur in the past
	if time.Now().After(date) {
		errVld.Add("date", "event's date must be in the future")
	}
	// check that provided resourcesAllowed exist and are known
	if len(resourcesAllowed) == 0 {
		errVld.Add("resourcesAllowed", "at least one resource must be specified")
	}
	for i, resource := range resourcesAllowed {
		if !resource.IsKnown() {
			errVld.Add(fmt.Sprintf("resourcesAllowed[%d]", i), fmt.Sprintf("unknown resource: %s", resource))
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
	// create an event
	event := &eventssvc.Event{
		ID:               id,
		Date:             date,
		Name:             name,
		ResourcesAllowed: resourcesAllowed,
	}
	// store the event
	return s.repo.StoreEvent(ctx, event)
}

// DeleteEvent deletes an event
func (s *service) DeleteEvent(ctx context.Context, eventID eventssvc.EventID) error {
	// check if there's eventID
	if len(eventID) == 0 {
		return valid.NewError("eventID", "eventID must be provided")
	}
	// delete the event
	return s.repo.DeleteEvent(ctx, eventID)
}

// EventByID returns an event by its ID
func (s *service) EventByID(ctx context.Context, eventID eventssvc.EventID) (*Event, error) {
	// check if there's eventID
	if len(eventID) == 0 {
		return nil, valid.NewError("eventID", "eventID is needed")
	}
	return s.repo.EventByID(ctx, eventID)
}

// EventClasses returns an array of sorted classes for the specified event
func (s *service) EventClasses(ctx context.Context, eventID eventssvc.EventID, filters EventClassFilters,
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
func (s *service) EventPupils(ctx context.Context, eventID eventssvc.EventID, filters EventPupilFilters,
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
func (s *service) PupilByID(ctx context.Context, pupilID eventssvc.PupilID, eventID eventssvc.EventID) (*Pupil, error) {
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
	if a <= 0 || a > MaxAmount {
		a = DefaultAmount
	}
	if s < 0 {
		s = DefaultSkip
	}
	return a, s
}

// Class is a model of the class, adapted for this use case.
type Class struct {
	// Name of the class as it was on the date of the event
	Name string
	// Resources brought by the class to the event
	ResourcesBrought eventssvc.ResourceMap
}

// Event is a model of the event, adapted for this use case.
// It indicates how many resources have been collected at this event
type Event struct {
	eventssvc.Event
	// resources brought by pupils to this event
	ResourcesBrought eventssvc.ResourceMap
}

// Pupil is a model of the pupil, adapted for this use case.
type Pupil struct {
	eventssvc.Pupil
	// Note that Class here is a string, not a eventssvc.Class instance.
	// It is the name of the class as it was on the date of the event
	Class string
	// Resources brought by the pupil to the event
	ResourcesBrought eventssvc.ResourceMap
}

// EventClassFilters are used to filter classes participating in an event
type EventClassFilters struct {
	Name string
}

// EventPupilFilters are used to filter pupils participating in an event
type EventPupilFilters struct {
	NameAndClass string
}
