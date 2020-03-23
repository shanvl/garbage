package mock

import (
	"context"

	"github.com/shanvl/garbage-events-service/aggregating"
	"github.com/shanvl/garbage-events-service/eventing"
	"github.com/shanvl/garbage-events-service/garbage"
	"github.com/shanvl/garbage-events-service/schooling"
	"github.com/shanvl/garbage-events-service/sorting"
)

// AggregatingRepository is a mock repository for aggregating use case
type AggregatingRepository struct {
	ClassesFn func(ctx context.Context, filters aggregating.ClassesFilters, classesSorting, eventsSorting sorting.By,
		amount, skip int) (classes []*aggregating.Class, total int, err error)
	ClassesInvoked bool
}

// Classes calls ClassesFn
func (r *AggregatingRepository) Classes(ctx context.Context, filters aggregating.ClassesFilters, classesSorting,
	eventsSorting sorting.By, amount, skip int) (classes []*aggregating.Class, total int, err error) {

	r.ClassesInvoked = true
	return r.ClassesFn(ctx, filters, classesSorting, eventsSorting, amount, skip)
}

// EventingRepository is a mock repository for eventing use case
type EventingRepository struct {
	ClassByIDFn      func(ctx context.Context, classID garbage.ClassID, eventID garbage.EventID) (*eventing.Class, error)
	ClassByIDInvoked bool

	ChangeEventResourcesFn func(ctx context.Context, eventID garbage.EventID,
		pupilID garbage.PupilID, resources map[garbage.Resource]int) (*eventing.Event, *eventing.Pupil, error)
	ChangeEventResourcesInvoked bool

	DeleteEventFn      func(ctx context.Context, id garbage.EventID) (garbage.EventID, error)
	DeleteEventInvoked bool

	EventByIDFn      func(ctx context.Context, id garbage.EventID) (*eventing.Event, error)
	EventByIDInvoked bool

	EventPupilsFn func(ctx context.Context, eventID garbage.EventID, filters eventing.EventPupilsFilters,
		sortBy sorting.By, amount int, skip int) ([]*eventing.Pupil, int, error)
	EventPupilsInvoked bool

	EventClassesFn func(ctx context.Context, eventID garbage.EventID, filters eventing.EventClassesFilters,
		sortBy sorting.By, amount int, skip int) ([]*eventing.Class, int, error)
	EventClassesInvoked bool

	PupilByIDFn      func(ctx context.Context, pupilID garbage.PupilID, eventID garbage.EventID) (*eventing.Pupil, error)
	PupilByIDInvoked bool

	StoreEventFn      func(ctx context.Context, e *garbage.Event) (garbage.EventID, error)
	StoreEventInvoked bool
}

// ChangeEventResources calls ChangeEventResourcesFn
func (r *EventingRepository) ChangeEventResources(ctx context.Context, eventID garbage.EventID,
	pupilID garbage.PupilID, resources map[garbage.Resource]int) (*eventing.Event, *eventing.Pupil, error) {
	r.ChangeEventResourcesInvoked = true
	return r.ChangeEventResourcesFn(ctx, eventID, pupilID, resources)
}

// ClassByID calls ClassByIDFn
func (r *EventingRepository) ClassByID(ctx context.Context, classID garbage.ClassID,
	eventID garbage.EventID) (*eventing.Class, error) {

	r.ClassByIDInvoked = true
	return r.ClassByIDFn(ctx, classID, eventID)
}

// DeleteEvent calls DeleteEventFn
func (r *EventingRepository) DeleteEvent(ctx context.Context, id garbage.EventID) (garbage.EventID, error) {
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
	filters eventing.EventClassesFilters, sortBy sorting.By, amount int, skip int) ([]*eventing.Class, int, error) {
	r.EventClassesInvoked = true
	return r.EventClassesFn(ctx, eventID, filters, sortBy, amount, skip)
}

// EventPupils calls EventPupilsFn
func (r *EventingRepository) EventPupils(ctx context.Context, eventID garbage.EventID,
	filters eventing.EventPupilsFilters, sortBy sorting.By, amount int, skip int) ([]*eventing.Pupil, int, error) {
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
	ClassFn      func(ctx context.Context, letter string, yearFormed int) (*garbage.Class, error)
	ClassInvoked bool

	PupilByIDFn      func(ctx context.Context, pupilID garbage.PupilID) (*schooling.Pupil, error)
	PupilByIDInvoked bool

	RemovePupilsFn      func(ctx context.Context, pupilIDs []garbage.PupilID) ([]garbage.PupilID, error)
	RemovePupilsInvoked bool

	StorePupilFn      func(ctx context.Context, pupils *schooling.Pupil) (garbage.PupilID, error)
	StorePupilInvoked bool

	StorePupilsFn      func(ctx context.Context, pupils []*schooling.Pupil) ([]garbage.PupilID, error)
	StorePupilsInvoked bool
}

// Class calls ClassFn
func (r *SchoolingRepository) Class(ctx context.Context, letter string, yearFormed int) (*garbage.Class, error) {
	r.ClassInvoked = true
	return r.ClassFn(ctx, letter, yearFormed)
}

// PupilByID calls PupilByIDFn
func (r *SchoolingRepository) PupilByID(ctx context.Context, pupilID garbage.PupilID) (*schooling.Pupil, error) {
	r.PupilByIDInvoked = true
	return r.PupilByIDFn(ctx, pupilID)
}

// RemovePupils calls RemovePupilsFn
func (r *SchoolingRepository) RemovePupils(ctx context.Context, pupilIDs []garbage.PupilID) ([]garbage.PupilID, error) {
	r.RemovePupilsInvoked = true
	return r.RemovePupilsFn(ctx, pupilIDs)
}

// StorePupil calls StorePupilFn
func (r *SchoolingRepository) StorePupil(ctx context.Context, pupil *schooling.Pupil) (garbage.PupilID, error) {
	r.StorePupilInvoked = true
	return r.StorePupilFn(ctx, pupil)
}

// StorePupils calls StorePupilsFn
func (r *SchoolingRepository) StorePupils(ctx context.Context, pupils []*schooling.Pupil) ([]garbage.PupilID, error) {
	r.StorePupilsInvoked = true
	return r.StorePupilsFn(ctx, pupils)
}
