// Package eventing is responsible for event management
package eventing

import (
	"context"
	"time"

	"github.com/shanvl/garbage-events-service"
	"github.com/shanvl/garbage-events-service/idgen"
	"github.com/shanvl/garbage-events-service/valid"
)

// Service is an interface providing methods to manage events
type Service interface {
	// CreateEvent creates and stores an event
	CreateEvent(ctx context.Context, date time.Time, name string, resources []garbage.Resource) (garbage.EventID, error)
	// DeleteEvent deletes an event
	DeleteEvent(ctx context.Context, eventID garbage.EventID) (garbage.EventID, error)
	// Events returns an array of sorted events
	Events(ctx context.Context, name string, date time.Time, sortBy SortBy, amount int,
		skip int) (events []*Event, total int, err error)
}

// Repository provides methods to work with event's persistence
type Repository interface {
	StoreEvent(ctx context.Context, event *garbage.Event) (garbage.EventID, error)
	DeleteEvent(ctx context.Context, eventID garbage.EventID) (garbage.EventID, error)
	Events(ctx context.Context, name string, date time.Time, sortBy SortBy, amount int, skip int) (events []*Event,
		total int, err error)
}

type service struct {
	repo Repository
}

// CreateEvent creates and stores an event
func (s *service) CreateEvent(ctx context.Context, date time.Time, name string,
	resourcesAllowed []garbage.Resource) (garbage.EventID, error) {
	// new event mustn't occur in the past
	validateDate := func() (isValid bool, errKey string, errDesc string) {
		if time.Now().After(date) {
			return false, "date", "event's date must be in the future"
		}
		return true, "", ""
	}
	// check that provided resourcesAllowed exist and are known
	validateResources := func() (isValid bool, errKey string, errDesc string) {
		if len(resourcesAllowed) == 0 {
			return false, "resourcesAllowed", "at least one resource must be specified"
		}
		for _, resource := range resourcesAllowed {
			if !resource.IsKnown() {
				return false, "resourcesAllowed", "unknown resource"
			}
		}
		return true, "", ""
	}
	// use previous functions to validate the arguments
	if err := valid.Check(validateDate, validateResources); err != nil {
		return "", err
	}
	id, err := idgen.CreateEventID()
	if err != nil {
		return "", err
	}
	event := garbage.NewEvent(id, date, name, resourcesAllowed)
	eventID, err := s.repo.StoreEvent(ctx, event)
	if err != nil {
		return "", err
	}
	return eventID, nil
}

// DeleteEvent deletes an event
func (s *service) DeleteEvent(ctx context.Context, eventID garbage.EventID) (garbage.EventID, error) {
	// check if there's eventID
	if err := valid.Check(func() (isValid bool, errKey string, errDesc string) {
		if len(eventID) <= 0 {
			return false, "eventID", "eventID must be provided"
		}
		return true, "", ""
	}); err != nil {
		return "", err
	}
	deletedID, err := s.repo.DeleteEvent(ctx, eventID)
	if err != nil {
		return "", err
	}
	return deletedID, nil
}

// Events returns an array of sorted events
func (s *service) Events(ctx context.Context, name string, date time.Time, sortBy SortBy, amount int,
	skip int) (events []*Event, total int, err error) {
	if amount < 0 {
		amount = 0
	}
	if skip < 0 {
		skip = 0
	}
	if !sortBy.IsValid() {
		sortBy = DateDesc
	}
	e, t, err := s.repo.Events(ctx, name, date, sortBy, amount, skip)
	if err != nil {
		return nil, 0, err
	}
	return e, t, nil
}

// NewService returns an instance of Service with all its dependencies
func NewService(repo Repository) Service {
	return &service{repo}
}

type Class struct {
	garbage.Class
	ResourcesBrought map[garbage.Resource]int
}

type Event struct {
	garbage.Event
	ResourcesCollected map[garbage.Resource]int
}

type Pupil struct {
	garbage.Pupil
	ResourcesBrought map[garbage.Resource]int
}
