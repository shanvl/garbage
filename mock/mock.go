package mock

import (
	"context"

	"github.com/shanvl/garbage-events-service/eventing"
	"github.com/shanvl/garbage-events-service/garbage"
	"github.com/shanvl/garbage-events-service/schooling"
	"github.com/shanvl/garbage-events-service/sorting"
)

// EventingRepository is a mock repository for eventing use case
type EventingRepository struct {
	ChangeEventResourcesFn func(ctx context.Context, eventID garbage.EventID,
		pupilID garbage.PupilID, resources map[garbage.Resource]int) (*eventing.Event, *eventing.Pupil, error)
	ChangeEventResourcesInvoked bool

	DeleteEventFn      func(ctx context.Context, id garbage.EventID) (garbage.EventID, error)
	DeleteEventInvoked bool

	EventFn      func(ctx context.Context, id garbage.EventID) (*eventing.Event, error)
	EventInvoked bool

	StoreEventFn      func(ctx context.Context, e *garbage.Event) (garbage.EventID, error)
	StoreEventInvoked bool

	EventPupilsFn func(ctx context.Context, eventID garbage.EventID, sortBy sorting.By, amount int,
		skip int) ([]*eventing.Pupil, int, error)
	EventPupilsInvoked bool

	EventClassesFn func(ctx context.Context, eventID garbage.EventID, sortBy sorting.By, amount int,
		skip int) ([]*eventing.Class, int, error)
	EventClassesInvoked bool
}

// DeleteEvent calls DeleteEvent
func (r *EventingRepository) DeleteEvent(ctx context.Context, id garbage.EventID) (garbage.EventID, error) {
	r.StoreEventInvoked = true
	return r.DeleteEventFn(ctx, id)
}

// StoreEvent calls StoreEventFn
func (r *EventingRepository) StoreEvent(ctx context.Context, e *garbage.Event) (garbage.EventID, error) {
	r.StoreEventInvoked = true
	return r.StoreEventFn(ctx, e)
}

// EventByID returns an event by eventID
func (r *EventingRepository) EventByID(ctx context.Context, id garbage.EventID) (*eventing.Event, error) {
	r.EventInvoked = true
	return r.EventFn(ctx, id)
}

// ChangeEventResources adds/subtracts resources brought by a pupil to/from the event
func (r *EventingRepository) ChangeEventResources(ctx context.Context, eventID garbage.EventID,
	pupilID garbage.PupilID, resources map[garbage.Resource]int) (*eventing.Event, *eventing.Pupil, error) {
	r.ChangeEventResourcesInvoked = true
	return r.ChangeEventResourcesFn(ctx, eventID, pupilID, resources)
}

// EventPupils returns an array of sorted pupils for the specified event
func (r *EventingRepository) EventPupils(ctx context.Context, eventID garbage.EventID, sortBy sorting.By, amount int,
	skip int) ([]*eventing.Pupil, int, error) {
	r.EventPupilsInvoked = true
	return r.EventPupilsFn(ctx, eventID, sortBy, amount, skip)
}

// EventClasses returns an array of sorted classes for the specified event
func (r *EventingRepository) EventClasses(ctx context.Context, eventID garbage.EventID, sortBy sorting.By, amount int,
	skip int) ([]*eventing.Class, int, error) {
	r.EventClassesInvoked = true
	return r.EventClassesFn(ctx, eventID, sortBy, amount, skip)
}

// SchoolingRepository is mock repository for schooling use case
type SchoolingRepository struct {
	ClassFn      func(ctx context.Context, letter string, yearFormed int) (*garbage.Class, error)
	ClassInvoked bool

	ChangePupilClassFn func(ctx context.Context, pupilID garbage.PupilID, className string) (garbage.PupilID,
		garbage.ClassID, error)
	ChangePupilClassInvoked bool

	RemovePupilsFn      func(ctx context.Context, pupilIDs []garbage.PupilID) ([]garbage.PupilID, error)
	RemovePupilsInvoked bool

	StorePupilsFn      func(ctx context.Context, pupils []*schooling.Pupil) ([]garbage.PupilID, error)
	StorePupilsInvoked bool
}

func (r *SchoolingRepository) Class(ctx context.Context, letter string, yearFormed int) (*garbage.Class, error) {
	r.ClassInvoked = true
	return r.ClassFn(ctx, letter, yearFormed)
}

// RemovePupils removes pupils and returns their IDs
func (r *SchoolingRepository) RemovePupils(ctx context.Context, pupilIDs []garbage.PupilID) ([]garbage.PupilID, error) {
	r.RemovePupilsInvoked = true
	return r.RemovePupilsFn(ctx, pupilIDs)
}

// StorePupils saves pupils to the repo and returns their IDs
func (r *SchoolingRepository) StorePupils(ctx context.Context, pupils []*schooling.Pupil) ([]garbage.PupilID, error) {
	r.StorePupilsInvoked = true
	return r.StorePupilsFn(ctx, pupils)
}

func (r *SchoolingRepository) ChangePupilClass(ctx context.Context, pupilID garbage.PupilID,
	className string) (garbage.PupilID, garbage.ClassID, error) {

	r.ChangePupilClassInvoked = true
	return r.ChangePupilClassFn(ctx, pupilID, className)
}
