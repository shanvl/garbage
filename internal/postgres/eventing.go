package postgres

import (
	"context"
	"database/sql"
	"errors"

	"github.com/jackc/pgconn"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/shanvl/garbage-events-service/internal/garbage"
	"github.com/shanvl/garbage-events-service/internal/sorting"
	"github.com/shanvl/garbage-events-service/internal/usecases/eventing"
)

type EventingRepo struct {
	db *pgxpool.Pool
}

func NewEventingRepo(db *pgxpool.Pool) *EventingRepo {
	return &EventingRepo{db}
}

const changePupilResourcesQuery = `
	insert into resources (pupil_id, event_id, gadgets, paper, plastic)
	values ($1, $2, $3, $4, $5)
	on conflict (pupil_id, event_id) do update
		set (gadgets, paper, plastic) = ($3, $4, $5);
`

// ChangePupilResources adds/subtracts resources brought by a pupil to/from the event, updating `resources` table
func (e *EventingRepo) ChangePupilResources(ctx context.Context, eventID garbage.EventID, pupilID garbage.PupilID,
	resources garbage.ResourceMap) error {

	_, err := e.db.Exec(ctx, changePupilResourcesQuery, pupilID, eventID, resources[garbage.Gadgets],
		resources[garbage.Paper], resources[garbage.Plastic])

	// violation of the foreign key (pupil_id, event_id) means that there has been no pupil or event found
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		if pgErr.Code == foreignKeyViolationCode {
			err = eventing.ErrNoEventPupil
		}
	}
	return err
}

const deleteEventQuery = `
	delete
    from event e
	where e.id = $1;
`

// DeleteEvent deletes an event with the id passed
func (e *EventingRepo) DeleteEvent(ctx context.Context, eventID garbage.EventID) error {
	_, err := e.db.Exec(ctx, deleteEventQuery, eventID)
	return err
}

const eventByIDQuery = `
	select e.id,
       e.name,
       e.date,
       e.resources_allowed::text[],
       coalesce(sum(gadgets), 0) as gadgets,
       coalesce(sum(r.paper), 0)   as paper,
       coalesce(sum(plastic), 0) as plastic
	from event e
			 left join resources r on e.id = r.event_id
	where e.id = $1
	group by e.id;
`

// EventByID returns an event by its ID
func (e *EventingRepo) EventByID(ctx context.Context, eventID garbage.EventID) (*eventing.Event, error) {
	ev := &eventing.Event{}
	// need this one in order to scan resources_allowed into it
	var resAllowedStr []string
	// query db
	err := e.db.QueryRow(ctx, eventByIDQuery, eventID).Scan(&ev.ID, &ev.Name, &ev.Date,
		&resAllowedStr, &ev.ResourcesBrought.Gadgets, &ev.ResourcesBrought.Paper, &ev.ResourcesBrought.Plastic)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, garbage.ErrNoEvent
		}
		return nil, err
	}
	// convert a string slice resources_allowed to a slice of resources
	resAllowed, err := garbage.StringSliceToResourceSlice(resAllowedStr)
	if err != nil {
		return nil, err
	}
	ev.ResourcesAllowed = resAllowed
	return ev, nil
}

func (e *EventingRepo) EventClasses(ctx context.Context, eventID garbage.EventID,
	filters eventing.EventClassesFilters, sortBy sorting.By, amount int, skip int) (classes []*eventing.Class,
	total int, err error) {

	panic("implement me")
}

func (e *EventingRepo) EventPupils(ctx context.Context, eventID garbage.EventID, filters eventing.EventPupilsFilters,
	sortBy sorting.By, amount int, skip int) (pupils []*eventing.Pupil, total int, err error) {
	panic("implement me")
}

func (e *EventingRepo) PupilByID(ctx context.Context, pupilID garbage.PupilID,
	eventID garbage.EventID) (*eventing.Pupil, error) {

	panic("implement me")
}

// ::text[]::resource[] is a workaround for pgx to save an enum array. It won't break or slow anything
const storeEventQuery = `
	insert into event (id, name, date, resources_allowed)
	values ($1, $2, $3, $4::text[]::resource[]);
`

// StoreEvent stores event into the db
func (e *EventingRepo) StoreEvent(ctx context.Context, event *garbage.Event) error {
	_, err := e.db.Exec(ctx, storeEventQuery, event.ID, event.Name, event.Date,
		garbage.ResourceSliceToStringSlice(event.ResourcesAllowed))
	return err
}
