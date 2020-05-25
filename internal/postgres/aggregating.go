package postgres

import (
	"context"
	"fmt"
	"strings"

	"github.com/jackc/pgtype"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/jmoiron/sqlx"
	"github.com/shanvl/garbage-events-service/internal/garbage"
	"github.com/shanvl/garbage-events-service/internal/sorting"
	"github.com/shanvl/garbage-events-service/internal/usecases/aggregating"
)

// AggregatingRepo is a repository used by Aggregating service
type AggregatingRepo struct {
	db *pgxpool.Pool
}

// NewAggregatingRepo returns an instance of AggregatingRepo
func NewAggregatingRepo(db *pgxpool.Pool) *AggregatingRepo {
	return &AggregatingRepo{db}
}

var eventOrderMap = map[sorting.By]string{
	sorting.NameAsc: "name asc, date desc",
	sorting.NameDes: "name desc, date desc",
	sorting.DateAsc: "date, name",
	sorting.DateDes: "date desc, name asc",
	sorting.Gadgets: "gadgets desc, date desc",
	sorting.Paper:   "paper desc, date desc",
	sorting.Plastic: "plastic desc, date desc",
}

func (a *AggregatingRepo) Classes(ctx context.Context, filters aggregating.ClassesFilters, classesSorting,
	eventsSorting sorting.By, amount, skip int) (classes []*aggregating.Class, total int, err error) {
	panic("implement me")
}

func (a *AggregatingRepo) Pupils(ctx context.Context, filters aggregating.PupilsFilters, pupilsSorting,
	eventsSorting sorting.By, amount, skip int) (pupils []*aggregating.Pupil, total int, err error) {
	panic("implement me")
}

const pupilByIDQueryA = `
	select e.id                 as event_id,
       e.name,
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
			 left join event e on e.date between symmetric greatest(p.class_date_formed, $1)
		and least(p.class_date_formed + (365.25 * 11)::integer, $2)
			 left join resources r on p.id = r.pupil_id and e.id = r.event_id
	where p.id = $3
	order by %s;
`

// PupilByID returns a pupil with the given ID, with a list of all the resources they has brought to every event that
// passed the provided filter
func (a *AggregatingRepo) PupilByID(ctx context.Context, id garbage.PupilID, filters aggregating.EventsByDateFilter,
	eventsSorting sorting.By) (*aggregating.Pupil, error) {

	// get and add "order by" to the query
	orderBy := eventOrderMap[eventsSorting]
	q := fmt.Sprintf(pupilByIDQueryA, orderBy)

	// the query selects the lesser date from "p.class_date_formed + (365.25 * 11)" and "filters.To". So if filters.
	// To is not set, set it to some date in the distant future
	if filters.To.IsZero() {
		filters.To = filters.To.AddDate(2222, 0, 0)
	}

	rows, err := a.db.Query(ctx, q, filters.From, filters.To, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// next types are only needed if no events have been found.
	// we need them to distinguish between the absence of the pupil and the absence of the events
	var (
		eID               pgtype.Varchar
		eDate             pgtype.Date
		eName             pgtype.Varchar
		eResourcesAllowed []string
	)
	p := &aggregating.Pupil{}
	for rows.Next() {
		e := aggregating.Event{}
		if err := rows.Scan(&eID, &eName, &eDate, &eResourcesAllowed, &p.ID, &p.FirstName, &p.LastName, &p.Letter,
			&p.DateFormed, &e.ResourcesBrought.Gadgets, &e.ResourcesBrought.Paper,
			&e.ResourcesBrought.Plastic); err != nil {
			return nil, err
		}
		// if event id is null, then no events have been found for the dates passed.
		// We simply return the pupil with an empty event slice
		if eID.Status != pgtype.Present {
			return p, nil
		}
		// set the event fields
		e.ID = garbage.EventID(eID.String)
		e.Date = eDate.Time
		e.Name = eName.String
		resAllowed, err := garbage.StringSliceToResourceSlice(eResourcesAllowed)
		if err != nil {
			return nil, err
		}
		e.ResourcesAllowed = resAllowed
		// append the event to the pupil's slice of events
		p.Events = append(p.Events, e)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	// if p.ID isn't set, then no pupil has been found
	if p.ID == "" {
		return nil, garbage.ErrUnknownPupil
	}
	return p, nil
}

const eventsQuery = `
	with query as (
		select e.id,
			   e.date,
			   e.name,
			   e.resources_allowed::text[],
			   coalesce(sum(gadgets), 0) as gadgets,
			   coalesce(sum(paper), 0)   as paper,
			   coalesce(sum(plastic), 0) as plastic
		from event e
				 left join resources r on r.event_id = e.id
		where 1=1 %s
		group by e.id
	), pagination as (
		select *
		from query
		order by %s
		limit ? offset ?
	)
	select *
		from pagination
			right join (select count(*) from query) as c(total) on true;
`

// Events returns a list of sorted events that passed the provided filters
func (a *AggregatingRepo) Events(ctx context.Context, filters aggregating.EventsFilters, sortBy sorting.By, amount,
	skip int) (events []*aggregating.Event, total int, err error) {

	// get the "order by" part of the query
	orderBy := eventOrderMap[sortBy]
	// create the "where" part of the query
	where := strings.Builder{}
	var args []interface{}
	if !filters.From.IsZero() || !filters.To.IsZero() {
		where.WriteString("and e.date between symmetric ? and ? ")
		// If From is set and To is not set, set it to some date in the distant future
		if filters.To.IsZero() {
			filters.To = filters.To.AddDate(2222, 0, 0)
		}
		args = append(args, filters.From, filters.To)
	}
	if len(filters.ResourcesAllowed) > 0 {
		where.WriteString("and e.resources_allowed @> ?::text[]::resource[] ")
		args = append(args, garbage.ResourceSliceToStringSlice(filters.ResourcesAllowed))
	}
	if filters.Name != "" {
		textSearch := prepareTextSearch(filters.Name)
		where.WriteString("and e.text_search @@ to_tsquery('simple', ?) ")
		args = append(args, textSearch)
	}
	args = append(args, amount, skip)
	// embed the "where" and the "order by" parts to the query
	q := fmt.Sprintf(eventsQuery, where.String(), orderBy)
	// change "?" to "$" in the query
	q = sqlx.Rebind(sqlx.BindType("pgx"), q)

	rows, err := a.db.Query(ctx, q, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()
	var (
		id            pgtype.Varchar
		date          pgtype.Date
		name          pgtype.Varchar
		resAllowedStr []string
		gadgets       pgtype.Float4
		paper         pgtype.Float4
		plastic       pgtype.Float4
	)
	for rows.Next() {
		err := rows.Scan(&id, &date, &name, &resAllowedStr, &gadgets, &paper, &plastic, &total)
		if err != nil {
			return nil, 0, err
		}
		// next will happen if the offset >= total rows found or no events matching the provided criteria have been
		// found. In that case we simply return total w/o additional work
		if id.Status != pgtype.Present {
			return nil, total, nil
		}
		// convert []string to []garbage.Resource
		resAllowed, err := garbage.StringSliceToResourceSlice(resAllowedStr)
		if err != nil {
			return nil, 0, err
		}
		e := &aggregating.Event{
			Event: garbage.Event{
				ID:               garbage.EventID(id.String),
				Date:             date.Time,
				Name:             name.String,
				ResourcesAllowed: resAllowed,
			},
			ResourcesBrought: garbage.Resources{
				Gadgets: gadgets.Float,
				Paper:   paper.Float,
				Plastic: plastic.Float,
			},
		}
		events = append(events, e)
	}
	if err = rows.Err(); err != nil {
		return nil, 0, err
	}
	return events, total, nil
}
