package postgres

import (
	"context"
	"errors"
	"time"

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
			return nil, garbage.ErrUnknownEvent
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

//
// var eventClassesOrderMap = map[sorting.By]string{
// 	sorting.NameAsc: "class_date_formed, class_letter asc",
// 	sorting.NameDes: "class_date_formed, class_letter desc",
// 	sorting.Gadgets: "gadgets desc",
// 	sorting.Paper:   "paper desc",
// 	sorting.Plastic: "plastic desc",
// }
//
// const eventClassesQuery = `
// 	with query as (
// 		select p.class_letter,
// 			   p.class_date_formed,
// 			   e.date,
// 			   coalesce(sum(paper), 0)   as paper,
// 			   coalesce(sum(plastic), 0) as plastic,
// 			   coalesce(sum(gadgets), 0) as gadgets
// 		from pupil p
// 				 cross join event e
// 				 left join resources r on r.pupil_id = p.id and r.event_id = e.id
// 		where e.id = ? %s
// 		  and e.date between symmetric p.class_date_formed and p.class_date_formed + (365.25 * 11)::integer
// 		group by p.class_date_formed, p.class_letter, e.date
// 	),
// 		 pagination as (
// 			 select *
// 			 from query
// 			 order by ?
// 			 limit ? offset ?
// 		 )
// 	select *
// 	from pagination
// 			 right join (SELECT count(*) FROM query) as c(total) on true;
// `

func (e *EventingRepo) EventClasses(ctx context.Context, eventID garbage.EventID,
	filters eventing.EventClassesFilters, sortBy sorting.By, amount int, skip int) (classes []*eventing.Class,
	total int, err error) {
	/*
		orderBy := eventClassesOrderMap[sortBy]
		// if there's a text filter, derive a garbage.Class from it
		if filters.Name != "" {
			// the event's date is needed when we create a garbage.Class from it's name.
			// A class' name changes depending on the event date
			var eDate time.Time
			if err := e.db.QueryRow(ctx, `select date from event where id = $1`, eventID).Scan(&eDate); err != nil {
				if errors.Is(err, pgx.ErrNoRows) {
					return nil, 0, garbage.ErrUnknownEvent
				}
				return nil, 0, err
			}
			// get the class' letter and date formed, using the event's date
			letter, dateFormed, err := garbage.ParseClassName(filters.Name, eDate)
			if err != nil {
				return nil, 0, err
			}
			// create a "where" clause based on the presence of the letter and the dateFormed
			var args []interface{}
			where := ""
			if len(letter) > 0 {
				where += " and p.class_letter = ? "
				args = append(args, letter)
			}
			if !dateFormed.IsZero() {
				where += " and p.class_date_formed = ? "
				args = append(args, dateFormed)
			}
			// add the where clause to the query
			q := fmt.Sprintf(eventClassesQuery, where)
			q, args, err = sqlx.In(q, eventID, args, orderBy, amount, skip)
			q = sqlx.Rebind(sqlx.BindType("pgx"), q)
			// execute the query
			e.db.Query(ctx, q, args...)
		}
		// use the constant query w/o the class' letter or date

		panic("implement me")*/
	panic("implement me")
}

const eventPupilsQuery = `
	with query as (
    select p.id,
           p.first_name,
           p.last_name,
           p.class_letter,
           p.class_date_formed,
		   e.date,
           coalesce(r.gadgets, 0) as gadgets,
           coalesce(r.paper, 0)   as paper,
           coalesce(r.plastic, 0) as plastic
    from pupil p
             cross join event e
             left join resources r on r.pupil_id = p.id and r.event_id = e.id
    where e.id = $1
      and e.date between symmetric p.class_date_formed and p.class_date_formed + (365.25 * 11)::integer
	),
	pagination as (
		select *
		from query
		order by $2
		limit $3 offset $4
	)
	select *
	from pagination
	right join (SELECT count(*) FROM query) as c(total) on true;
`

const eventPupilsQueryTextSearch = `
	with query as (
    select p.id,
           p.first_name,
           p.last_name,
           p.class_letter,
           p.class_date_formed,
           e.date,
           coalesce(r.gadgets, 0) as gadgets,
           coalesce(r.paper, 0)   as paper,
           coalesce(r.plastic, 0) as plastic
    from pupil p
             cross join event e
             left join resources r on r.pupil_id = p.id and r.event_id = e.id
    where e.id = $1
      and p.text_search @@ to_tsquery('simple', $2)
      and e.date between symmetric p.class_date_formed and p.class_date_formed + (365.25 * 11)::integer
	),
	pagination as (
		select *
		from query
		order by $3
		limit $4 offset $5
	)
	select *
	from pagination
	right join (SELECT count(*) FROM query) as c(total) on true;
`

var eventPupilsOrderMap = map[sorting.By]string{
	sorting.NameAsc: "class_date_formed, class_letter, last_name, first_name asc",
	sorting.NameDes: "class_date_formed, class_letter, last_name, first_name desc",
	sorting.Gadgets: "gadgets desc",
	sorting.Paper:   "paper desc",
	sorting.Plastic: "plastic desc",
}

// EventPupils returns a paginated and sorted list of the pupils who have participated in the specified event
func (e *EventingRepo) EventPupils(ctx context.Context, eventID garbage.EventID, filters eventing.EventPupilsFilters,
	sortBy sorting.By, amount int, skip int) (pupils []*eventing.Pupil, total int, err error) {

	var rows pgx.Rows
	// derive orderBy query from the sortBy passed
	orderBy := eventPupilsOrderMap[sortBy]
	// if there's no text filter, just query the db. Otherwise,
	// make a text search query from it and query the db with it
	if filters.NameAndClass == "" {
		rows, err = e.db.Query(ctx, eventPupilsQuery, eventID, orderBy, amount, skip)
	} else {
		// we need to know the event's date in order to create a text search query. Every word,
		// which resembles a class name, will be copied,
		// processed and concatenated with itself so as to hit the table indices. The event's date is needed there.
		// "3B" will become "3B:* | 2018B:*" if the event's date is 10.10.2020. Also,
		// the event's date is needed when we create a class name from it's letter and the year it was formed in
		var eDate time.Time
		if err := e.db.QueryRow(ctx, `select date from event where id = $1`, eventID).Scan(&eDate); err != nil {
			if errors.Is(err, pgx.ErrNoRows) {
				return nil, 0, garbage.ErrUnknownEvent
			}
			return nil, 0, err
		}
		// create a query from the filters passed
		textSearchQuery := prepareTextSearchQuery(filters.NameAndClass, eDate)
		// if the text query is empty, query the db without it. Otherwise, with it
		if textSearchQuery == "" {
			rows, err = e.db.Query(ctx, eventPupilsQuery, eventID, orderBy, amount, skip)
		} else {
			rows, err = e.db.Query(ctx, eventPupilsQueryTextSearch, eventID, textSearchQuery, orderBy, amount, skip)
		}
	}
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()
	// next types are only needed in case the offset is >= total rows found. Then all the columns, except 'total',
	// will be null
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
	)

	for rows.Next() {
		err := rows.Scan(&id, &firstName, &lastName, &classLetter, &classDate, &eventDate, &gadgets, &paper, &plastic,
			&total)
		if err != nil {
			return nil, 0, err
		}
		// this will only happen if the offset is >= total rows found.
		// In that case we simply return total.
		if id.Status != pgtype.Present {
			return nil, total, nil
		}
		p := &eventing.Pupil{
			Pupil: garbage.Pupil{
				ID:        garbage.PupilID(id.String),
				FirstName: firstName.String,
				LastName:  lastName.String,
			},
			ResourcesBrought: garbage.Resources{
				Gadgets: gadgets.Float,
				Paper:   paper.Float,
				Plastic: plastic.Float,
			},
		}
		c := garbage.Class{Letter: classLetter.String, DateFormed: classDate.Time}
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
	c := garbage.Class{}

	err := e.db.QueryRow(ctx, evPupilByIDQuery, eventID, pupilID).Scan(&eID, &eDate, &p.ID, &p.FirstName, &p.LastName,
		&c.Letter, &c.DateFormed, &p.ResourcesBrought.Gadgets, &p.ResourcesBrought.Paper, &p.ResourcesBrought.Plastic)

	if err != nil {
		// if no rows have been returned, there's no such a pupil
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, garbage.ErrUnknownPupil
		}
		return nil, err
	}
	// if eID is null, there's no such an event
	if eID.Status != pgtype.Present {
		return nil, garbage.ErrUnknownEvent
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
