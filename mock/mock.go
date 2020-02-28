package mock

import (
	"context"
	"time"

	"github.com/shanvl/garbage-events-service/eventing"
	"github.com/shanvl/garbage-events-service/garbage"
)

// EventingRepository is a mock eventing usecase repository
type EventingRepository struct {
	ChangeEventResourcesFn func(ctx context.Context, eventID garbage.EventID,
		pupilID garbage.PupilID, resources map[garbage.Resource]int) (*garbage.Event, *garbage.Pupil, error)
	ChangeEventResourcesInvoked bool

	DeleteEventFn      func(ctx context.Context, id garbage.EventID) (garbage.EventID, error)
	DeleteEventInvoked bool

	EventFn      func(ctx context.Context, id garbage.EventID) (*garbage.Event, error)
	EventInvoked bool

	EventsFn func(ctx context.Context, name string, date time.Time, sortBy eventing.SortBy, amount int,
		skip int) (events []*garbage.Event, total int, err error)
	EventsInvoked bool

	StoreEventFn      func(ctx context.Context, e *garbage.Event) (garbage.EventID, error)
	StoreEventInvoked bool
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

// Events returns an array of sorted events
func (r *EventingRepository) Events(ctx context.Context, name string, date time.Time, sortBy eventing.SortBy,
	amount int, skip int) (events []*garbage.Event, total int, err error) {
	r.EventsInvoked = true
	return r.EventsFn(ctx, name, date, sortBy, amount, skip)
}

// Event returns an event by eventID
func (r *EventingRepository) Event(ctx context.Context, id garbage.EventID) (*garbage.Event, error) {
	r.EventInvoked = true
	return r.EventFn(ctx, id)
}

// ChangeEventResources adds/subtracts resources brought by a pupil to/from the event
func (r *EventingRepository) ChangeEventResources(ctx context.Context, eventID garbage.EventID,
	pupilID garbage.PupilID, resources map[garbage.Resource]int) (*garbage.Event, *garbage.Pupil, error) {
	r.ChangeEventResourcesInvoked = true
	return r.ChangeEventResourcesFn(ctx, eventID, pupilID, resources)
}
