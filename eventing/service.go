// Package eventing is responsible for event management
package eventing

import (
	"context"
	"fmt"
	"time"

	"github.com/shanvl/garbage-events-service/garbage"
	"github.com/shanvl/garbage-events-service/idgen"
	"github.com/shanvl/garbage-events-service/valid"
)

// Service is an interface providing methods to manage events
type Service interface {
	// CreateEvent creates and stores an event
	CreateEvent(ctx context.Context, date time.Time, name string, resources []garbage.Resource) (garbage.EventID, error)
	// ChangeEventResources adds/subtracts resources brought by a pupil to/from the event
	ChangeEventResources(ctx context.Context, eventID garbage.EventID, pupilID garbage.PupilID,
		resources map[garbage.Resource]int) (*garbage.Event, *garbage.Pupil, error)
	// DeleteEvent deletes an event
	DeleteEvent(ctx context.Context, eventID garbage.EventID) (garbage.EventID, error)
	// Event returns an event by its ID
	Event(ctx context.Context, eventID garbage.EventID) (*garbage.Event, error)
	// Events returns an array of sorted events
	Events(ctx context.Context, name string, date time.Time, sortBy SortBy, amount int,
		skip int) (events []*garbage.Event, total int, err error)
}

// Repository provides methods to work with event's persistence
type Repository interface {
	ChangeEventResources(ctx context.Context, eventID garbage.EventID, pupilID garbage.PupilID,
		resources map[garbage.Resource]int) (*garbage.Event, *garbage.Pupil, error)
	DeleteEvent(ctx context.Context, eventID garbage.EventID) (garbage.EventID, error)
	Event(ctx context.Context, eventID garbage.EventID) (*garbage.Event, error)
	Events(ctx context.Context, name string, date time.Time, sortBy SortBy, amount int,
		skip int) (events []*garbage.Event,
		total int, err error)
	StoreEvent(ctx context.Context, event *garbage.Event) (garbage.EventID, error)
}

const (
	DefaultAmount = 5
	DefaultSkip   = 0
	DefaultSort   = DateDesc
)

type service struct {
	repo Repository
}

// ChangeEventResources adds/subtracts resources brought by a pupil to/from the event
func (s *service) ChangeEventResources(ctx context.Context, eventID garbage.EventID, pupilID garbage.PupilID,
	resources map[garbage.Resource]int) (*garbage.Event, *garbage.Pupil, error) {

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
	event, err := s.repo.Event(ctx, eventID)
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
	if !errVld.IsEmpty() {
		return "", errVld
	}
	// use previous functions to validate the arguments
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

// Events returns an array of sorted events
func (s *service) Events(ctx context.Context, name string, date time.Time, sortBy SortBy, amount int,
	skip int) (events []*garbage.Event, total int, err error) {

	// if provided values are incorrect, use default values instead
	if amount <= 0 {
		amount = DefaultAmount
	}
	if skip < 0 {
		skip = DefaultSkip
	}
	if !sortBy.IsValid() {
		sortBy = DefaultSort
	}
	// get the events
	events, total, err = s.repo.Events(ctx, name, date, sortBy, amount, skip)
	if err != nil {
		return nil, 0, err
	}
	return events, total, nil
}

// Event returns an event by its ID
func (s *service) Event(ctx context.Context, eventID garbage.EventID) (*garbage.Event, error) {
	errVld := valid.EmptyError()
	if len(eventID) <= 0 {
		errVld.Add("eventID", "eventID is needed")
	}
	if !errVld.IsEmpty() {
		return nil, errVld
	}
	event, err := s.repo.Event(ctx, eventID)
	if err != nil {
		return nil, err
	}
	return event, nil
}

// NewService returns an instance of Service with all its dependencies
func NewService(repo Repository) Service {
	return &service{repo}
}
