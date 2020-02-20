// Package eventing is responsible for event management
package eventing

import (
	"time"

	"github.com/shanvl/garbage-events-service"
)

// Service is an interface providing methods to manage events
type Service interface {
	// CreateEvent creates and stores an event
	CreateEvent(date time.Time, name string, resources []garbage.Resource) (garbage.EventID, error)
}

// Repository provides methods to work with event's persistence
type Repository interface {
	StoreEvent(event *garbage.Event) (garbage.EventID, error)
}

// IDGenerator is used to to generate unique IDs
type IDGenerator interface {
	GenerateEventID() garbage.EventID
}

// Validator provides functions helping with params validation and returning a convenient error
type Validator interface {
	Validate(validateFunctions ...func() (isValid bool, errKey string, errDesc string)) error
}

type service struct {
	repo  Repository
	idGen IDGenerator
	valid Validator
}

// CreateEvent creates and stores an event
func (s *service) CreateEvent(date time.Time, name string, resourcesAllowed []garbage.Resource) (garbage.EventID, error) {
	// new event mustn't occur in the past
	validateDate := func() (isValid bool, errKey string, errDesc string) {
		if time.Now().After(date) {
			return false, "date", "event's date must be in the future"
		}
		return true, "", ""
	}
	// check that given resources are known
	validateResources := func() (isValid bool, errKey string, errDesc string) {
		for _, resource := range resourcesAllowed {
			if !resource.IsValid() {
				return false, "resourcesAllowed", "unknown resource"
			}
		}
		return true, "", ""
	}
	// use previous functions to validate the arguments
	if err := s.valid.Validate(validateDate, validateResources); err != nil {
		return "", err
	}
	id := s.idGen.GenerateEventID()
	event := garbage.NewEvent(id, date, name, resourcesAllowed)
	eventID, err := s.repo.StoreEvent(event)
	if err != nil {
		return "", err
	}
	return eventID, nil
}

// NewService returns an instance of Service w/ all its dependencies
func NewService(repo Repository, idGen IDGenerator, valid Validator) Service {
	return &service{repo, idGen, valid}
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
