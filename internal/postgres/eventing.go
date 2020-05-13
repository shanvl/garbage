package postgres

import (
	"context"
	"errors"

	"github.com/jackc/pgconn"
	"github.com/jackc/pgtype"
	"github.com/jackc/pgx/v4"
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
		if errors.Is(err, pgx.ErrNoRows) {
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

// doing it via the left join so as to differ between the absence of the event or the pupil.
// If no rows have been returned, then there's no pupil, otherwise, if e.id is null, there's no event
const evPupilByIDQuery = `
	select e.id,
           e.date,
		   p.id,
		   p.first_name,
		   p.last_name,
		   p.class_letter,
		   p.class_year_formed,
		   coalesce(gadgets, 0) as gadgets,
		   coalesce(paper, 0)   as paper,
		   coalesce(plastic, 0) as plastic
	from pupil p
			 left join event e on e.id = $1
			 left join resources r on e.id = r.event_id and p.id = r.pupil_id
	where p.id = $2;
`

func (e *EventingRepo) PupilByID(ctx context.Context, pupilID garbage.PupilID,
	eventID garbage.EventID) (*eventing.Pupil, error) {

	p := &eventing.Pupil{
		Pupil: garbage.Pupil{
			ID:        "",
			FirstName: "",
			LastName:  "",
		},
		Class:            "",
		ResourcesBrought: garbage.Resources{},
	}

	// if there's no e.id has been returned, then there's no event with such an event_id
	var eID pgtype.Varchar
	// class name is always relative to the date of the event, not to the current date
	var eDate pgtype.Date
	// class instance to derive a class name from
	c := &garbage.Class{}

	err := e.db.QueryRow(ctx, evPupilByIDQuery, eventID, pupilID).Scan(&eID, &eDate, &p.ID, &p.FirstName, &p.LastName,
		&c.Letter, &c.YearFormed, &p.ResourcesBrought.Gadgets, &p.ResourcesBrought.Paper, &p.ResourcesBrought.Plastic)

	if err != nil {
		// if no rows have been returned, there's no such a pupil
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, garbage.ErrNoPupil
		}
		return nil, err
	}
	// if eID is null, there's no such an event
	if eID.Status != pgtype.Present {
		return nil, garbage.ErrNoEvent
	}

	// create a class name and assign it to the pupil
	class, err := c.NameOnDate(eDate.Time)
	if err != nil {
		// if there was no such class on the event's date, then the pupil didn't participate in the event
		if errors.Is(err, garbage.ErrNoClassOnDate) {
			return nil, eventing.ErrNoEventPupil
		}
		return nil, err
	}
	p.Class = class

	return p, nil
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
