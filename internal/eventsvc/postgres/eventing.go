package postgres

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/jackc/pgconn"
	"github.com/jackc/pgtype"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/jmoiron/sqlx"
	"github.com/shanvl/garbage/internal/eventsvc"
	"github.com/shanvl/garbage/internal/eventsvc/eventing"
	"github.com/shanvl/garbage/internal/eventsvc/sorting"
)

// eventingRepo is a repository used by Eventing service
type eventingRepo struct {
	db *pgxpool.Pool
}

// NewEventingRepo returns an instance of eventingRepo
func NewEventingRepo(db *pgxpool.Pool) eventing.Repository {
	return &eventingRepo{db}
}

const changePupilResourcesQuery = `
	insert into resources (pupil_id, event_id, gadgets, paper, plastic)
	values ($1, $2, $3, $4, $5)
	on conflict (pupil_id, event_id) do update
		set (gadgets, paper, plastic) = ($3, $4, $5);
`

// ChangePupilResources adds/subtracts resources brought by a pupil to/from the event, updating `resources` table
func (e *eventingRepo) ChangePupilResources(ctx context.Context, eventID string, pupilID string,
	resources eventsvc.ResourceMap) error {

	_, err := e.db.Exec(ctx, changePupilResourcesQuery, pupilID, eventID, resources[eventsvc.Gadgets],
		resources[eventsvc.Paper], resources[eventsvc.Plastic])

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
func (e *eventingRepo) DeleteEvent(ctx context.Context, eventID string) error {
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
func (e *eventingRepo) EventByID(ctx context.Context, eventID string) (*eventing.Event, error) {
	ev := &eventing.Event{}
	// need this one in order to scan resources_allowed into it
	var (
		resAllowedStr []string
		gadgets       float32
		paper         float32
		plastic       float32
	)
	// query db
	err := e.db.QueryRow(ctx, eventByIDQuery, eventID).Scan(&ev.ID, &ev.Name, &ev.Date,
		&resAllowedStr, &gadgets, &paper, &plastic)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, eventsvc.ErrUnknownEvent
		}
		return nil, err
	}
	// convert a string slice resources_allowed to a slice of resources
	resAllowed, err := eventsvc.StringSliceToResourceSlice(resAllowedStr)
	if err != nil {
		return nil, err
	}
	ev.ResourcesAllowed = resAllowed
	ev.ResourcesBrought = newResourceMap(resAllowed, gadgets, paper, plastic)
	return ev, nil
}

const eventClassesQuery = `
	with query as (
		select p.class_letter,
			   p.class_date_formed,
			   e.date,
               e.resources_allowed::text[],
			   coalesce(sum(gadgets), 0) as gadgets,
			   coalesce(sum(paper), 0)   as paper,
			   coalesce(sum(plastic), 0) as plastic
		from pupil p
				 cross join event e
				 left join resources r on r.pupil_id = p.id and r.event_id = e.id
		where e.id = ? %s
		  and e.date between symmetric p.class_date_formed and p.class_date_formed + (365.25 * 11)::integer
		group by p.class_date_formed, p.class_letter, e.date, e.resources_allowed
	),  pagination as (
			 select *
			 from query
			 order by %s
			 limit ? offset ?
)
	select *
	from pagination
			 right join (select count(*) FROM query) as c(total) on true
             left join (select id from event where id = ?) as d(event_id) on true;
`

// EventClasses returns a sorted and paginated list of classes that match the passed filters
func (e *eventingRepo) EventClasses(ctx context.Context, eventID string,
	filters eventing.EventClassFilters, sortBy sorting.By, amount int, skip int) (classes []*eventing.Class,
	total int, err error) {

	var q string
	// derive the "order by" query part from the sortBy passed
	orderBy := classOrderMap[sortBy]
	// create a slice of the query arguments
	args := []interface{}{eventID}
	// if there're no filters passed, create a simple query. Otherwise, create a query w/ a conditional "where" part
	if filters.Name == "" {
		q = fmt.Sprintf(eventClassesQuery, "", orderBy)
		args = append(args, amount, skip, eventID)
	} else {
		// the event's date is needed when we create a eventsvc.Class from it's name.
		// A class' name changes depending on the event date
		var eDate time.Time
		if err := e.db.QueryRow(ctx, `select date from event where id = $1`, eventID).Scan(&eDate); err != nil {
			if errors.Is(err, pgx.ErrNoRows) {
				return nil, 0, eventsvc.ErrUnknownEvent
			}
			return nil, 0, err
		}
		// get the class' letter and date formed, using the event's date
		letter, dateFormed, err := eventsvc.ParseClassName(filters.Name, eDate)
		// if the letter or the date can't be derived from the class name, then simply return the total w/o any errors
		if err != nil {
			return nil, 0, nil
		}
		// create the "where" clause based on the presence of the letter and the dateFormed and append the
		// arguments to the slice
		where := ""
		if letter != "" {
			where += " and p.class_letter = ? "
			args = append(args, letter)
		}
		if !dateFormed.IsZero() {
			where += " and p.class_date_formed = ? "
			args = append(args, dateFormed)
		}
		// add other arguments
		args = append(args, amount, skip, eventID)
		// add the where clause to the query
		q = fmt.Sprintf(eventClassesQuery, where, orderBy)
	}
	// swap "?" for "$" in the query
	q = sqlx.Rebind(sqlx.BindType("pgx"), q)
	// exec the query
	rows, err := e.db.Query(ctx, q, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	// "total" column will always be returned, so other columns might be null
	var (
		classLetter pgtype.Varchar
		classDate   pgtype.Date
		eventDate   pgtype.Date
		gadgets     pgtype.Float4
		paper       pgtype.Float4
		plastic     pgtype.Float4
		eID         pgtype.Varchar
	)
	for rows.Next() {
		var resAllowedStr []string
		if err := rows.Scan(&classLetter, &classDate, &eventDate, &resAllowedStr, &gadgets, &paper, &plastic, &total,
			&eID); err != nil {
			return nil, 0, err
		}
		// if the event hasn't been found, return an error
		if eID.Status != pgtype.Present {
			return nil, total, eventsvc.ErrUnknownEvent
		}
		// next will happen if the offset >= total rows found or no classes with such class names have been found
		// In that case we simply return the total w/o additional work
		if classLetter.Status != pgtype.Present {
			return nil, total, nil
		}
		// derive a className from its letter, dateFormed and the event's date
		c := eventsvc.Class{Letter: classLetter.String, DateFormed: classDate.Time}
		className, err := c.NameOnDate(eventDate.Time)
		if err != nil {
			return nil, 0, err
		}
		// convert []string to []eventsvc.Resource
		resAllowed, err := eventsvc.StringSliceToResourceSlice(resAllowedStr)
		// create a map of the resources brought to the event
		resBrought := newResourceMap(resAllowed, gadgets.Float, paper.Float, plastic.Float)
		if err != nil {
			return nil, 0, err
		}
		// append the class to other classes
		classes = append(classes, &eventing.Class{
			Name:             className,
			ResourcesBrought: resBrought,
		})
	}
	if err := rows.Err(); err != nil {
		return nil, 0, err
	}
	return classes, total, nil
}

const eventPupilsQuery = `
	with query as (
    select p.id,
           p.first_name,
           p.last_name,
           p.class_letter,
           p.class_date_formed,
		   e.date,
		   e.resources_allowed::text[],
           coalesce(r.gadgets, 0) as gadgets,
           coalesce(r.paper, 0)   as paper,
           coalesce(r.plastic, 0) as plastic
    from pupil p
             cross join event e
             left join resources r on r.pupil_id = p.id and r.event_id = e.id
    where e.id = ? %s
      and e.date between symmetric p.class_date_formed and p.class_date_formed + (365.25 * 11)::integer
	),
	pagination as (
		select *
		from query
		order by %s
		limit ? offset ?
	)
	select *
	from pagination
	right join (select count(*) from query) as c(total) on true
	left join (select id from event where id = ?) as d(event_id) on true;
`

// EventPupils returns a paginated and sorted list of the pupils who have participated in the specified event
func (e *eventingRepo) EventPupils(ctx context.Context, eventID string, filters eventing.EventPupilFilters,
	sortBy sorting.By, amount int, skip int) (pupils []*eventing.Pupil, total int, err error) {

	var q string
	// derive the "order by" query part from the sortBy passed
	orderBy := pupilOrderMap[sortBy]
	// create a slice of the query arguments
	args := []interface{}{eventID}
	// if there're no filters passed, create a simple query. Otherwise, create a query w/ the text search
	if filters.NameAndClass == "" {
		q = fmt.Sprintf(eventPupilsQuery, "", orderBy)
		args = append(args, amount, skip, eventID)
	} else {
		// we need to know the event's date in order to create a text search query. Every word,
		// which resembles a class name, will be copied,
		// processed and concatenated with itself so as to hit the table indices. The event's date is needed there.
		// "3B" will become "3B:* | 2018B:*" if the event's date is 10.10.2020. Also,
		// the event's date is needed when we create a class name from it's letter and the year it was formed in
		var eDate time.Time
		if err := e.db.QueryRow(ctx, `select date from event where id = $1`, eventID).Scan(&eDate); err != nil {
			if errors.Is(err, pgx.ErrNoRows) {
				return nil, 0, eventsvc.ErrUnknownEvent
			}
			return nil, 0, err
		}
		// create the text search part of the query from the filters passed
		textSearchQuery := prepareTextSearchClass(filters.NameAndClass, eDate)
		q = fmt.Sprintf(eventPupilsQuery, " and p.text_search @@ to_tsquery('simple', ?)", orderBy)
		args = append(args, textSearchQuery, amount, skip, eventID)
	}
	// change "?" to "$" in the query
	q = sqlx.Rebind(sqlx.BindType("pgx"), q)
	rows, err := e.db.Query(ctx, q, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()
	// "total" column will always be returned, so other columns might be null
	var (
		id          pgtype.Varchar
		firstName   pgtype.Varchar
		lastName    pgtype.Varchar
		classLetter pgtype.Varchar
		classDate   pgtype.Date
		eventDate   pgtype.Date
		gadgets     pgtype.Float4
		paper       pgtype.Float4
		plastic     pgtype.Float4
		eID         pgtype.Varchar
	)

	for rows.Next() {
		var resAllowedStr []string
		err := rows.Scan(&id, &firstName, &lastName, &classLetter, &classDate, &eventDate, &resAllowedStr, &gadgets,
			&paper, &plastic, &total, &eID)
		if err != nil {
			return nil, 0, err
		}
		if eID.Status != pgtype.Present {
			return nil, total, eventsvc.ErrUnknownEvent
		}
		// next will happen if the offset >= total rows found or no pupils with such names/classNames were found
		// In that case we simply return the total w\o additional work
		if id.Status != pgtype.Present {
			return nil, total, nil
		}
		// convert []string to []eventsvc.Resource
		resAllowed, err := eventsvc.StringSliceToResourceSlice(resAllowedStr)
		// create a map of the resources brought to the event
		resBrought := newResourceMap(resAllowed, gadgets.Float, paper.Float, plastic.Float)
		if err != nil {
			return nil, 0, err
		}
		p := &eventing.Pupil{
			Pupil: eventsvc.Pupil{
				ID:        id.String,
				FirstName: firstName.String,
				LastName:  lastName.String,
			},
			ResourcesBrought: resBrought,
		}
		c := eventsvc.Class{Letter: classLetter.String, DateFormed: classDate.Time}
		// derive a class name from its letter and a year it was formed in
		className, err := c.NameOnDate(eventDate.Time)
		if err != nil {
			return nil, 0, err
		}
		// append it to the pupil
		p.Class = className
		// append the pupil to the array of pupils
		pupils = append(pupils, p)
	}
	err = rows.Err()
	if err != nil {
		return nil, 0, err
	}

	return pupils, total, nil
}

// doing it via the left join so as to differ between the absence of the event or the pupil.
// If no rows have been returned, then there's no pupil, otherwise, if e.id is null, there's no event
const evPupilByIDQuery = `
	select e.id,
           e.date,
		   e.resources_allowed::text[],
		   p.id,
		   p.first_name,
		   p.last_name,
		   p.class_letter,
		   p.class_date_formed,
		   coalesce(gadgets, 0) as gadgets,
		   coalesce(paper, 0)   as paper,
		   coalesce(plastic, 0) as plastic
	from pupil p
			 left join event e on e.id = $1
			 left join resources r on e.id = r.event_id and p.id = r.pupil_id
	where p.id = $2;
`

func (e *eventingRepo) PupilByID(ctx context.Context, pupilID string,
	eventID string) (*eventing.Pupil, error) {

	p := &eventing.Pupil{}

	var (
		// if there's no e.id has been returned, then there's no event with such an event_id
		eID pgtype.Varchar
		// class name is always relative to the date of the event, not to the current date
		eDate          pgtype.Date
		eResAllowedStr []string
		gadgets        float32
		paper          float32
		plastic        float32
	)

	// class instance to derive a class name from
	c := eventsvc.Class{}

	err := e.db.QueryRow(ctx, evPupilByIDQuery, eventID, pupilID).Scan(&eID, &eDate, &eResAllowedStr, &p.ID,
		&p.FirstName, &p.LastName, &c.Letter, &c.DateFormed, &gadgets, &paper, &plastic)

	if err != nil {
		// if no rows have been returned, there's no such a pupil
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, eventsvc.ErrUnknownPupil
		}
		return nil, err
	}
	// if eID is null, there's no such an event
	if eID.Status != pgtype.Present {
		return nil, eventsvc.ErrUnknownEvent
	}

	// create a class name and assign it to the pupil
	class, err := c.NameOnDate(eDate.Time)
	if err != nil {
		// if there was no such class on the event's date, then the pupil didn't participate in the event
		if errors.Is(err, eventsvc.ErrNoClassOnDate) {
			return nil, eventing.ErrNoEventPupil
		}
		return nil, err
	}
	p.Class = class

	// convert []string to []eventsvc.Resource
	resAllowed, err := eventsvc.StringSliceToResourceSlice(eResAllowedStr)
	// create a map of the resources brought to the event
	resBrought := newResourceMap(resAllowed, gadgets, paper, plastic)
	if err != nil {
		return nil, err
	}
	p.ResourcesBrought = resBrought

	return p, nil
}

// ::text[]::resource[] is a workaround for pgx to save an enum array. It won't break or slow anything
const storeEventQuery = `
	insert into event (id, name, date, resources_allowed)
	values ($1, $2, $3, $4::text[]::resource[]);
`

// StoreEvent stores event into the db
func (e *eventingRepo) StoreEvent(ctx context.Context, event *eventsvc.Event) error {
	_, err := e.db.Exec(ctx, storeEventQuery, event.ID, event.Name, event.Date,
		eventsvc.ResourceSliceToStringSlice(event.ResourcesAllowed))
	return err
}
