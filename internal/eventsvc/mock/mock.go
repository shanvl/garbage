package mock

import (
	"context"

	"github.com/shanvl/garbage/internal/eventsvc"
	"github.com/shanvl/garbage/internal/eventsvc/aggregating"
	"github.com/shanvl/garbage/internal/eventsvc/eventing"
	"github.com/shanvl/garbage/internal/eventsvc/schooling"
	"github.com/shanvl/garbage/internal/eventsvc/sorting"
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

	PupilByIDFn func(ctx context.Context, id string, filters aggregating.EventFilters,
		eventsSorting sorting.By) (*aggregating.Pupil, error)
	PupilByIDInvoked bool
}

func (r *AggregatingRepository) Classes(ctx context.Context, filters aggregating.ClassFilters, classesSorting,
	eventsSorting sorting.By, amount, skip int) (classes []*aggregating.Class, total int, err error) {

	r.ClassesInvoked = true
	return r.ClassesFn(ctx, filters, classesSorting, eventsSorting, amount, skip)
}

func (r *AggregatingRepository) Events(ctx context.Context, filters aggregating.EventFilters, sortBy sorting.By,
	amount, skip int) (events []*aggregating.Event, total int, err error) {

	r.EventsInvoked = true
	return r.EventsFn(ctx, filters, sortBy, amount, skip)
}

func (r *AggregatingRepository) Pupils(ctx context.Context, filters aggregating.PupilFilters, pupilsSorting,
	eventsSorting sorting.By, amount, skip int) (pupils []*aggregating.Pupil, total int, err error) {

	r.PupilsInvoked = true
	return r.PupilsFn(ctx, filters, pupilsSorting, eventsSorting, amount, skip)
}

func (r *AggregatingRepository) PupilByID(ctx context.Context, id string,
	filters aggregating.EventFilters, eventsSorting sorting.By) (*aggregating.Pupil, error) {

	r.PupilByIDInvoked = true
	return r.PupilByIDFn(ctx, id, filters, eventsSorting)
}

// EventingRepository is a mock repository for eventing use case
type EventingRepository struct {
	ChangePupilResourcesFn func(ctx context.Context, eventID string,
		pupilID string, resources eventsvc.ResourceMap) error
	ChangePupilResourcesInvoked bool

	DeleteEventFn      func(ctx context.Context, id string) error
	DeleteEventInvoked bool

	EventByIDFn      func(ctx context.Context, id string) (*eventing.Event, error)
	EventByIDInvoked bool

	EventPupilsFn func(ctx context.Context, eventID string, filters eventing.EventPupilFilters,
		sortBy sorting.By, amount int, skip int) ([]*eventing.Pupil, int, error)
	EventPupilsInvoked bool

	EventClassesFn func(ctx context.Context, eventID string, filters eventing.EventClassFilters,
		sortBy sorting.By, amount int, skip int) ([]*eventing.Class, int, error)
	EventClassesInvoked bool

	PupilByIDFn      func(ctx context.Context, pupilID string, eventID string) (*eventing.Pupil, error)
	PupilByIDInvoked bool

	StoreEventFn      func(ctx context.Context, e *eventsvc.Event) error
	StoreEventInvoked bool
}

// ChangePupilResources calls ChangeEventResourcesFn
func (r *EventingRepository) ChangePupilResources(ctx context.Context, eventID string,
	pupilID string, resources eventsvc.ResourceMap) error {
	r.ChangePupilResourcesInvoked = true
	return r.ChangePupilResourcesFn(ctx, eventID, pupilID, resources)
}

// DeleteEvent calls DeleteEventFn
func (r *EventingRepository) DeleteEvent(ctx context.Context, id string) error {
	r.StoreEventInvoked = true
	return r.DeleteEventFn(ctx, id)
}

// EventByID calls EventByIDFn
func (r *EventingRepository) EventByID(ctx context.Context, id string) (*eventing.Event, error) {
	r.EventByIDInvoked = true
	return r.EventByIDFn(ctx, id)
}

// EventClasses calls EventClassesFn
func (r *EventingRepository) EventClasses(ctx context.Context, eventID string,
	filters eventing.EventClassFilters, sortBy sorting.By, amount int, skip int) ([]*eventing.Class, int, error) {
	r.EventClassesInvoked = true
	return r.EventClassesFn(ctx, eventID, filters, sortBy, amount, skip)
}

// EventPupils calls EventPupilsFn
func (r *EventingRepository) EventPupils(ctx context.Context, eventID string,
	filters eventing.EventPupilFilters, sortBy sorting.By, amount int, skip int) ([]*eventing.Pupil, int, error) {
	r.EventPupilsInvoked = true
	return r.EventPupilsFn(ctx, eventID, filters, sortBy, amount, skip)
}

// PupilByID calls PupilByIDFn
func (r *EventingRepository) PupilByID(ctx context.Context, pupilID string,
	eventID string) (*eventing.Pupil, error) {

	r.PupilByIDInvoked = true
	return r.PupilByIDFn(ctx, pupilID, eventID)
}

// StoreEvent calls StoreEventFn
func (r *EventingRepository) StoreEvent(ctx context.Context, e *eventsvc.Event) error {
	r.StoreEventInvoked = true
	return r.StoreEventFn(ctx, e)
}

// SchoolingRepository is mock repository for schooling use case
type SchoolingRepository struct {
	PupilByIDFn      func(ctx context.Context, pupilID string) (*schooling.Pupil, error)
	PupilByIDInvoked bool

	RemovePupilsFn      func(ctx context.Context, pupilIDs []string) error
	RemovePupilsInvoked bool

	UpdatePupilFn     func(ctx context.Context, pupils *schooling.Pupil) error
	StorePupilInvoked bool

	StorePupilsFn      func(ctx context.Context, pupils []*schooling.Pupil) error
	StorePupilsInvoked bool
}

// PupilByID calls PupilByIDFn
func (r *SchoolingRepository) PupilByID(ctx context.Context, pupilID string) (*schooling.Pupil, error) {
	r.PupilByIDInvoked = true
	return r.PupilByIDFn(ctx, pupilID)
}

// RemovePupils calls RemovePupilsFn
func (r *SchoolingRepository) RemovePupils(ctx context.Context, pupilIDs []string) error {
	r.RemovePupilsInvoked = true
	return r.RemovePupilsFn(ctx, pupilIDs)
}

// StorePupil calls StorePupilFn
func (r *SchoolingRepository) UpdatePupil(ctx context.Context, pupil *schooling.Pupil) error {
	r.StorePupilInvoked = true
	return r.UpdatePupilFn(ctx, pupil)
}

// StorePupils calls StorePupilsFn
func (r *SchoolingRepository) StorePupils(ctx context.Context, pupils []*schooling.Pupil) error {
	r.StorePupilsInvoked = true
	return r.StorePupilsFn(ctx, pupils)
}
