package mock

import (
	"context"

	"github.com/shanvl/garbage-events-service/internal/garbage"
	"github.com/shanvl/garbage-events-service/internal/sorting"
	"github.com/shanvl/garbage-events-service/internal/usecases/aggregating"
	"github.com/shanvl/garbage-events-service/internal/usecases/eventing"
	"github.com/shanvl/garbage-events-service/internal/usecases/schooling"
)

// AggregatingRepository is a mock repository for aggregating use case
type AggregatingRepository struct {
	ClassesFn func(ctx context.Context, filters aggregating.ClassFilters, classesSorting, eventsSorting sorting.By,
		amount, skip int) (classes []*aggregating.Class, total int, err error)
	ClassesInvoked bool

	EventsFn func(ctx context.Context, filters aggregating.EventFilters, sortBy sorting.By, amount,
		skip int) (events []*aggregating.Event, total int, err error)
	EventsInvoked bool

	PupilsFn func(ctx context.Context, filters aggregating.PupilFilters, pupilsSorting, eventsSorting sorting.By, amount,
		skip int) (pupils []*aggregating.Pupil, total int, err error)
	PupilsInvoked bool

	PupilByIDFn func(ctx context.Context, id garbage.PupilID, filters aggregating.EventDateFilters,
		eventsSorting sorting.By) (*aggregating.Pupil, error)
	PupilByIDInvoked bool
}

// Classes calls ClassesFn
func (r *AggregatingRepository) Classes(ctx context.Context, filters aggregating.ClassFilters, classesSorting,
	eventsSorting sorting.By, amount, skip int) (classes []*aggregating.Class, total int, err error) {

	r.ClassesInvoked = true
	return r.ClassesFn(ctx, filters, classesSorting, eventsSorting, amount, skip)
}

// Events calls EventsFn
func (r *AggregatingRepository) Events(ctx context.Context, filters aggregating.EventFilters, sortBy sorting.By,
	amount, skip int) (events []*aggregating.Event, total int, err error) {

	r.EventsInvoked = true
	return r.EventsFn(ctx, filters, sortBy, amount, skip)
}

// Pupils calls PupilsFn
func (r *AggregatingRepository) Pupils(ctx context.Context, filters aggregating.PupilFilters, pupilsSorting,
	eventsSorting sorting.By, amount, skip int) (pupils []*aggregating.Pupil, total int, err error) {

	r.PupilsInvoked = true
	return r.PupilsFn(ctx, filters, pupilsSorting, eventsSorting, amount, skip)
}

// PupilByID calls PupilByIDFn
func (r *AggregatingRepository) PupilByID(ctx context.Context, id garbage.PupilID,
	filters aggregating.EventDateFilters, eventsSorting sorting.By) (*aggregating.Pupil, error) {

	r.PupilByIDInvoked = true
	return r.PupilByIDFn(ctx, id, filters, eventsSorting)
}

// EventingRepository is a mock repository for eventing use case
type EventingRepository struct {
	ChangePupilResourcesFn func(ctx context.Context, eventID garbage.EventID,
		pupilID garbage.PupilID, resources garbage.ResourceMap) error
	ChangePupilResourcesInvoked bool

	DeleteEventFn      func(ctx context.Context, id garbage.EventID) error
	DeleteEventInvoked bool

	EventByIDFn      func(ctx context.Context, id garbage.EventID) (*eventing.Event, error)
	EventByIDInvoked bool

	EventPupilsFn func(ctx context.Context, eventID garbage.EventID, filters eventing.EventPupilFilters,
		sortBy sorting.By, amount int, skip int) ([]*eventing.Pupil, int, error)
	EventPupilsInvoked bool

	EventClassesFn func(ctx context.Context, eventID garbage.EventID, filters eventing.EventClassFilters,
		sortBy sorting.By, amount int, skip int) ([]*eventing.Class, int, error)
	EventClassesInvoked bool

	PupilByIDFn      func(ctx context.Context, pupilID garbage.PupilID, eventID garbage.EventID) (*eventing.Pupil, error)
	PupilByIDInvoked bool

	StoreEventFn      func(ctx context.Context, e *garbage.Event) (garbage.EventID, error)
	StoreEventInvoked bool
}

// ChangePupilResources calls ChangeEventResourcesFn
func (r *EventingRepository) ChangePupilResources(ctx context.Context, eventID garbage.EventID,
	pupilID garbage.PupilID, resources garbage.ResourceMap) error {
	r.ChangePupilResourcesInvoked = true
	return r.ChangePupilResourcesFn(ctx, eventID, pupilID, resources)
}

// DeleteEvent calls DeleteEventFn
func (r *EventingRepository) DeleteEvent(ctx context.Context, id garbage.EventID) error {
	r.StoreEventInvoked = true
	return r.DeleteEventFn(ctx, id)
}

// EventByID calls EventByIDFn
func (r *EventingRepository) EventByID(ctx context.Context, id garbage.EventID) (*eventing.Event, error) {
	r.EventByIDInvoked = true
	return r.EventByIDFn(ctx, id)
}

// EventClasses calls EventClassesFn
func (r *EventingRepository) EventClasses(ctx context.Context, eventID garbage.EventID,
	filters eventing.EventClassFilters, sortBy sorting.By, amount int, skip int) ([]*eventing.Class, int, error) {
	r.EventClassesInvoked = true
	return r.EventClassesFn(ctx, eventID, filters, sortBy, amount, skip)
}

// EventPupils calls EventPupilsFn
func (r *EventingRepository) EventPupils(ctx context.Context, eventID garbage.EventID,
	filters eventing.EventPupilFilters, sortBy sorting.By, amount int, skip int) ([]*eventing.Pupil, int, error) {
	r.EventPupilsInvoked = true
	return r.EventPupilsFn(ctx, eventID, filters, sortBy, amount, skip)
}

// PupilByID calls PupilByIDFn
func (r *EventingRepository) PupilByID(ctx context.Context, pupilID garbage.PupilID,
	eventID garbage.EventID) (*eventing.Pupil, error) {

	r.PupilByIDInvoked = true
	return r.PupilByIDFn(ctx, pupilID, eventID)
}

// StoreEvent calls StoreEventFn
func (r *EventingRepository) StoreEvent(ctx context.Context, e *garbage.Event) (garbage.EventID, error) {
	r.StoreEventInvoked = true
	return r.StoreEventFn(ctx, e)
}

// SchoolingRepository is mock repository for schooling use case
type SchoolingRepository struct {
	PupilByIDFn      func(ctx context.Context, pupilID garbage.PupilID) (*schooling.Pupil, error)
	PupilByIDInvoked bool

	RemovePupilsFn      func(ctx context.Context, pupilIDs []garbage.PupilID) error
	RemovePupilsInvoked bool

	StorePupilFn      func(ctx context.Context, pupils *schooling.Pupil) error
	StorePupilInvoked bool

	StorePupilsFn      func(ctx context.Context, pupils []*schooling.Pupil) error
	StorePupilsInvoked bool
}

// PupilByID calls PupilByIDFn
func (r *SchoolingRepository) PupilByID(ctx context.Context, pupilID garbage.PupilID) (*schooling.Pupil, error) {
	r.PupilByIDInvoked = true
	return r.PupilByIDFn(ctx, pupilID)
}

// RemovePupils calls RemovePupilsFn
func (r *SchoolingRepository) RemovePupils(ctx context.Context, pupilIDs []garbage.PupilID) error {
	r.RemovePupilsInvoked = true
	return r.RemovePupilsFn(ctx, pupilIDs)
}

// StorePupil calls StorePupilFn
func (r *SchoolingRepository) StorePupil(ctx context.Context, pupil *schooling.Pupil) error {
	r.StorePupilInvoked = true
	return r.StorePupilFn(ctx, pupil)
}

// StorePupils calls StorePupilsFn
func (r *SchoolingRepository) StorePupils(ctx context.Context, pupils []*schooling.Pupil) error {
	r.StorePupilsInvoked = true
	return r.StorePupilsFn(ctx, pupils)
}
