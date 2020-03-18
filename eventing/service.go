package eventing

import (
	"context"
	"fmt"
	"time"

	"github.com/shanvl/garbage-events-service/garbage"
	"github.com/shanvl/garbage-events-service/idgen"
	"github.com/shanvl/garbage-events-service/sorting"
	"github.com/shanvl/garbage-events-service/valid"
)

// Service is an interface providing methods to manage an event.
// Note that all methods and entities are used in the context of one event.
type Service interface {
	// CreateEvent creates and stores an event
	CreateEvent(ctx context.Context, date time.Time, name string, resources []garbage.Resource) (garbage.EventID, error)
	// ChangeEventResources adds/subtracts resources brought by a pupil to/from the event
	ChangeEventResources(ctx context.Context, eventID garbage.EventID, pupilID garbage.PupilID,
		resources map[garbage.Resource]int) (*Event, *Pupil, error)
	// DeleteEvent deletes an event
	DeleteEvent(ctx context.Context, eventID garbage.EventID) (garbage.EventID, error)
	// EventByID returns an event by its ID
	EventByID(ctx context.Context, eventID garbage.EventID) (*Event, error)
	// EventPupils returns an array of sorted pupils for the specified event
	EventPupils(ctx context.Context, eventID garbage.EventID, sortBy sorting.By, amount int, skip int) ([]*Pupil, int,
		error)
	// EventClasses returns an array of sorted classes for the specified event
	EventClasses(ctx context.Context, eventID garbage.EventID, sortBy sorting.By, amount, skip int) ([]*Class, int, error)
}

// Repository provides methods to work with an event's persistence
type Repository interface {
	ChangeEventResources(ctx context.Context, eventID garbage.EventID, pupilID garbage.PupilID,
		resources map[garbage.Resource]int) (*Event, *Pupil, error)
	DeleteEvent(ctx context.Context, eventID garbage.EventID) (garbage.EventID, error)
	EventByID(ctx context.Context, eventID garbage.EventID) (*Event, error)
	EventClasses(ctx context.Context, eventID garbage.EventID, sortBy sorting.By, amount int, skip int) (classes []*Class,
		total int, err error)
	EventPupils(ctx context.Context, eventID garbage.EventID, sortBy sorting.By, amount int, skip int) (pupils []*Pupil,
		total int, err error)
	StoreEvent(ctx context.Context, event *garbage.Event) (garbage.EventID, error)
}

const (
	DefaultAmount = 5
	DefaultSkip   = 0
)

type service struct {
	repo Repository
}

// NewService returns an instance of Service with all its dependencies
func NewService(repo Repository) Service {
	return &service{repo}
}

// ChangeEventResources adds/subtracts resources brought by a pupil to/from the event
func (s *service) ChangeEventResources(ctx context.Context, eventID garbage.EventID, pupilID garbage.PupilID,
	resources map[garbage.Resource]int) (*Event, *Pupil, error) {

	errVld := valid.EmptyError()
	if len(pupilID) <= 0 {
		errVld.Add("pupilID", "pupilID must be provided")
	}
	if len(eventID) <= 0 {
		errVld.Add("eventID", "eventID must be provided")
	}
	if len(resources) <= 0 {
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
	event, pupil, err := s.repo.ChangeEventResources(ctx, eventID, pupilID, resources)
	if err != nil {
		return nil, nil, err
	}
	return event, pupil, nil
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
	if len(name) <= 0 {
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
	eventID, err := s.repo.StoreEvent(ctx, event)
	if err != nil {
		return "", err
	}
	return eventID, nil
}

// DeleteEvent deletes an event
func (s *service) DeleteEvent(ctx context.Context, eventID garbage.EventID) (garbage.EventID, error) {
	errVld := valid.EmptyError()
	// check if there's eventID
	if len(eventID) <= 0 {
		errVld.Add("eventID", "eventID must be provided")
	}
	if !errVld.IsEmpty() {
		return "", errVld
	}
	// delete event
	deletedID, err := s.repo.DeleteEvent(ctx, eventID)
	if err != nil {
		return "", err
	}
	return deletedID, nil
}

// EventByID returns an event by its ID
func (s *service) EventByID(ctx context.Context, eventID garbage.EventID) (*Event, error) {
	errVld := valid.EmptyError()
	if len(eventID) <= 0 {
		errVld.Add("eventID", "eventID is needed")
	}
	if !errVld.IsEmpty() {
		return nil, errVld
	}
	event, err := s.repo.EventByID(ctx, eventID)
	if err != nil {
		return nil, err
	}
	return event, nil
}

// EventClasses returns an array of sorted classes for the specified event
func (s *service) EventClasses(ctx context.Context, eventID garbage.EventID, sortBy sorting.By, amount,
	skip int) (classes []*Class, total int, err error) {

	// check if eventID was provided
	errVld := valid.EmptyError()
	if len(eventID) <= 0 {
		errVld.Add("eventID", "eventID must be provided")
	}
	if !errVld.IsEmpty() {
		return nil, 0, errVld
	}

	// if provided values are incorrect, use default values instead
	if amount <= 0 {
		amount = DefaultAmount
	}
	if skip < 0 {
		skip = DefaultSkip
	}
	if !sortBy.IsForEventClasses() {
		sortBy = sorting.NameAsc
	}

	classes, total, err = s.repo.EventClasses(ctx, eventID, sortBy, amount, skip)
	if err != nil {
		return nil, 0, err
	}
	return classes, total, nil
}

// EventPupils returns an array of sorted pupils for the specified event
func (s *service) EventPupils(ctx context.Context, eventID garbage.EventID, sortBy sorting.By, amount int,
	skip int) (pupils []*Pupil, total int, err error) {

	// check if eventID was provided
	errVld := valid.EmptyError()
	if len(eventID) <= 0 {
		errVld.Add("eventID", "eventID must be provided")
	}
	if !errVld.IsEmpty() {
		return nil, 0, errVld
	}

	// if provided values are incorrect, use default values instead
	if amount <= 0 {
		amount = DefaultAmount
	}
	if skip < 0 {
		skip = DefaultSkip
	}
	if !sortBy.IsForEventPupils() {
		sortBy = sorting.NameAsc
	}

	pupils, total, err = s.repo.EventPupils(ctx, eventID, sortBy, amount, skip)
	if err != nil {
		return nil, 0, err
	}
	return pupils, total, nil
}
