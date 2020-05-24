package postgres

import (
	"context"
	"fmt"

	"github.com/jackc/pgtype"
	"github.com/jackc/pgx/v4/pgxpool"
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
	sorting.NameAsc: "event_name asc, event_date desc",
	sorting.NameDes: "event_name desc, event_date desc",
	sorting.DateAsc: "event_date, event_name",
	sorting.DateDes: "event_date desc, event_name asc",
	sorting.Gadgets: "gadgets desc, event_date desc",
	sorting.Paper:   "paper desc, event_date desc",
	sorting.Plastic: "plastic desc, event_date desc",
}

func (a *AggregatingRepo) Classes(ctx context.Context, filters aggregating.ClassesFilters, classesSorting,
	eventsSorting sorting.By, amount, skip int) (classes []*aggregating.Class, total int, err error) {
	panic("implement me")
}

func (a *AggregatingRepo) ClassByID(ctx context.Context, id garbage.ClassID, filters aggregating.EventsByDateFilter,
	eventsSorting sorting.By) (*aggregating.Class, error) {
	panic("implement me")
}

func (a *AggregatingRepo) Pupils(ctx context.Context, filters aggregating.PupilsFilters, pupilsSorting,
	eventsSorting sorting.By, amount, skip int) (pupils []*aggregating.Pupil, total int, err error) {
	panic("implement me")
}

const pupilByIDQueryA = `
	select e.id                 as event_id,
       e.name                   as event_name,
       e.date                   as event_date,
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
		filters.To = filters.To.AddDate(999999, 0, 0)
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

func (a *AggregatingRepo) Events(ctx context.Context, filters aggregating.EventsFilters, sortBy sorting.By, amount,
	skip int) (events []*aggregating.Event, total int, err error) {
	panic("implement me")
}
