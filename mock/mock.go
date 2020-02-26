package mock

import (
	"context"
	"time"

	"github.com/shanvl/garbage-events-service"
	"github.com/shanvl/garbage-events-service/eventing"
)

// EventingRepository is a mock eventing usecase repository
type EventingRepository struct {
	DeleteEventFn      func(ctx context.Context, id garbage.EventID) (garbage.EventID, error)
	DeleteEventInvoked bool

	EventsFn func(ctx context.Context, name string, date time.Time, sortBy eventing.SortBy, amount int,
		skip int) (events []*eventing.Event, total int, err error)
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
	amount int, skip int) (events []*eventing.Event, total int, err error) {
	r.EventsInvoked = true
	return r.EventsFn(ctx, name, date, sortBy, amount, skip)
}
