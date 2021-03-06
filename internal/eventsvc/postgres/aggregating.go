package postgres

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/jackc/pgtype"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/jmoiron/sqlx"
	"github.com/shanvl/garbage/internal/eventsvc"
	"github.com/shanvl/garbage/internal/eventsvc/aggregating"
	"github.com/shanvl/garbage/internal/eventsvc/sorting"
	pgtextsearch "github.com/shanvl/garbage/pkg/pg-text-search"
)

// aggregatingRepo is a repository used by Aggregating service
type aggregatingRepo struct {
	db *pgxpool.Pool
}

// NewAggregatingRepo returns an instance of aggregatingRepo
func NewAggregatingRepo(db *pgxpool.Pool) aggregating.Repository {
	return &aggregatingRepo{db}
}

const classesQuery = `
	with query as (
		select class_date_formed,
			   class_letter,
			   e.id,
			   e.date,
			   e.name,
               e.resources_allowed,
			   sum(coalesce(gadgets, 0)) as gadgets,
			   sum(coalesce(paper, 0))   as paper,
			   sum(coalesce(plastic, 0)) as plastic
		from pupil p
				 cross join event e
				 left join resources r on r.event_id = e.id and r.pupil_id = p.id
		where e.date between symmetric greatest(p.class_date_formed, ?) and least(
				p.class_date_formed + (365.25 * 11)::integer, ?) %s
		group by class_date_formed, class_letter, e.id
	),
		 aggr as (
			 select class_date_formed,
					class_letter,
					sum(coalesce(gadgets, 0)) as gadgets_aggr,
					sum(coalesce(paper, 0))   as paper_aggr,
					sum(coalesce(plastic, 0)) as plastic_aggr
			 from query
			 group by class_date_formed, class_letter
		 ),
		 pagination as (
			 select *
			 from aggr
			 order by %s
			 limit ? offset ?
		 )
	select class_date_formed,
		   class_letter,
		   gadgets_aggr,
		   paper_aggr,
		   plastic_aggr,
		   id,
		   date,
		   name,
           resources_allowed::text[],
		   gadgets,
		   paper,
		   plastic,
		   total
	from query
			 inner join pagination using (class_date_formed, class_letter)
			 right join (select count(*) from aggr) as c(total) on true
	order by %s
`

// Classes returns a sorted list of the classes that passed the classes filters with a sorted list of the events that
// passed the event filters
func (a *aggregatingRepo) Classes(ctx context.Context, filters aggregating.ClassFilters, classesSorting,
	eventsSorting sorting.By, amount, skip int) (classes []*aggregating.Class, total int, err error) {

	// create the "order by" parts of the query
	classOrderBy := classAggrOrderMap[classesSorting]
	eventOrderBy := eventOrderMap[eventsSorting]
	orderBy := fmt.Sprintf("%s, class_date_formed desc, class_letter, %s", classOrderBy, eventOrderBy)

	// create the "where" part of the query
	where := strings.Builder{}
	var args []interface{}
	// if filters.To is not set, set it to some date in the distant future
	if filters.To.IsZero() {
		filters.To = filters.To.AddDate(2222, 0, 0)
	}
	args = append(args, filters.From, filters.To)
	if len(filters.ResourcesAllowed) > 0 {
		where.WriteString("and e.resources_allowed @> ?::text[]::resource[] ")
		args = append(args, eventsvc.ResourceSliceToStringSlice(filters.ResourcesAllowed))
	}
	// event's name text search
	if filters.Name != "" {
		eventTextSearch := pgtextsearch.PrepareQuery(filters.Name)
		where.WriteString("and e.text_search @@ to_tsquery('simple', ?) ")
		args = append(args, eventTextSearch)
	}
	if filters.Letter != "" {
		where.WriteString(" and p.class_letter = ? ")
		args = append(args, filters.Letter)
	}
	if !filters.DateFormed.IsZero() {
		where.WriteString(" and p.class_date_formed = ? ")
		args = append(args, filters.DateFormed)
	}

	// add limit and offset to the query
	args = append(args, amount, skip)
	q := fmt.Sprintf(classesQuery, where.String(), classOrderBy, orderBy)
	// change "?" to "$"
	q = sqlx.Rebind(sqlx.BindType("pgx"), q)

	rows, err := a.db.Query(ctx, q, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	// "total" column will always be returned, so other columns might be null
	var (
		cDate       pgtype.Date
		cLetter     pgtype.Varchar
		gadgetsAggr pgtype.Float4
		paperAggr   pgtype.Float4
		plasticAggr pgtype.Float4
		eID         pgtype.Varchar
		eDate       pgtype.Date
		eName       pgtype.Varchar
		eResAllowed []string
		gadgets     pgtype.Float4
		paper       pgtype.Float4
		plastic     pgtype.Float4
	)
	// map to fast search the class in the classes slice. class_year+class_letter -> index in the slice
	classSliceIndex := map[string]int{}
	for rows.Next() {
		if err = rows.Scan(&cDate, &cLetter, &gadgetsAggr, &paperAggr, &plasticAggr, &eID, &eDate, &eName,
			&eResAllowed, &gadgets, &paper, &plastic, &total); err != nil {
			return nil, 0, err
		}
		// next will happen if the offset >= total rows found or no classes matching the provided criteria have been
		// found. In that case we simply return total w/o additional work
		if cDate.Status != pgtype.Present {
			return nil, total, nil
		}
		classID := createClassID(cDate.Time, cLetter.String)
		var c *aggregating.Class
		if i, ok := classSliceIndex[classID]; ok {
			// if the class is already in the classes, get it
			c = classes[i]
		} else {
			// otherwise, create it
			c = &aggregating.Class{
				Class: eventsvc.Class{
					Letter:     cLetter.String,
					DateFormed: cDate.Time,
				},
				ResourcesBrought: eventsvc.ResourceMap{
					eventsvc.Gadgets: gadgetsAggr.Float,
					eventsvc.Paper:   paperAggr.Float,
					eventsvc.Plastic: plasticAggr.Float,
				},
			}
			// append the class to the slice and put its index to the map
			classes = append(classes, c)
			classSliceIndex[classID] = len(classes) - 1
		}
		resAllowed, err := eventsvc.StringSliceToResourceSlice(eResAllowed)
		if err != nil {
			return nil, 0, err
		}
		// create a map of the resources brought by the class to the event
		resBrought := newResourceMap(resAllowed, gadgets.Float, paper.Float, plastic.Float)
		// create an event and append it to the pupil's slice of events
		e := &aggregating.Event{
			Event: eventsvc.Event{
				ID:               eID.String,
				Date:             eDate.Time,
				Name:             eName.String,
				ResourcesAllowed: resAllowed,
			},
			ResourcesBrought: resBrought,
		}
		c.Events = append(c.Events, e)
	}
	if err = rows.Err(); err != nil {
		return nil, 0, err
	}
	return classes, total, nil
}

const pupilsQuery = `
	with query as (
    select p.id,
           first_name,
           last_name,
           class_date_formed,
           class_letter,
           e.id                 as event_id,
           e.date,
           e.name,
           e.resources_allowed,
           coalesce(gadgets, 0) as gadgets,
           coalesce(paper, 0)   as paper,
           coalesce(plastic, 0) as plastic

    from pupil p
             cross join event e
             left join resources r on r.event_id = e.id and r.pupil_id = p.id
    where e.date between symmetric greatest(p.class_date_formed, ?) and least(
            p.class_date_formed + (365.25 * 11)::integer, ?) %s
),
     aggr as (
         select id,
                first_name,
                last_name,
                class_date_formed,
                class_letter,
                sum(coalesce(gadgets, 0)) as gadgets_aggr,
                sum(coalesce(paper, 0))   as paper_aggr,
                sum(coalesce(plastic, 0)) as plastic_aggr
         from query
         group by id, class_date_formed, class_letter, first_name, last_name
     ),
     pagination as (
         select *
         from aggr
         order by %s
         limit ? offset ?
     )
select id,
       query.first_name,
       query.last_name,
       query.class_date_formed,
       query.class_letter,
       gadgets_aggr,
       paper_aggr,
       plastic_aggr,
       event_id,
       date,
       name,
       resources_allowed::text[],
       gadgets,
       paper,
       plastic,
       total
from query
         inner join pagination using (id)
         right join (select count(*) from aggr) as c(total) on true
order by %s;
`

// Pupils returns a sorted list of the pupils that passed the pupil filters with a sorted list of the events that
// passed the event filters
func (a *aggregatingRepo) Pupils(ctx context.Context, filters aggregating.PupilFilters, pupilsSorting,
	eventsSorting sorting.By, amount, skip int) (pupils []*aggregating.Pupil, total int, err error) {

	// create the "order by" parts of the query
	pupilOrderBy := pupilAggrOrderMap[pupilsSorting]
	eventOrderBy := eventOrderMap[eventsSorting]
	orderBy := fmt.Sprintf("%s, id, %s", pupilOrderBy, eventOrderBy)

	// create the "where" part of the query
	where := strings.Builder{}
	var args []interface{}
	// if filters.To is not set, set it to some date in the distant future
	if filters.To.IsZero() {
		filters.To = filters.To.AddDate(2222, 0, 0)
	}
	args = append(args, filters.From, filters.To)
	if len(filters.ResourcesAllowed) > 0 {
		where.WriteString("and e.resources_allowed @> ?::text[]::resource[] ")
		args = append(args, eventsvc.ResourceSliceToStringSlice(filters.ResourcesAllowed))
	}
	// event's name text search
	if filters.Name != "" {
		eventTextSearch := pgtextsearch.PrepareQuery(filters.Name)
		where.WriteString("and e.text_search @@ to_tsquery('simple', ?) ")
		args = append(args, eventTextSearch)
	}
	// pupil's name and class text search
	if filters.NameAndClass != "" {
		pupilTextSearch := pgtextsearch.PrepareQuery(filters.NameAndClass)
		where.WriteString("and p.text_search @@ to_tsquery('simple', ?)")
		args = append(args, pupilTextSearch)
	}

	// add limit and offset to the query
	args = append(args, amount, skip)
	q := fmt.Sprintf(pupilsQuery, where.String(), pupilOrderBy, orderBy)
	// change "?" to "$"
	q = sqlx.Rebind(sqlx.BindType("pgx"), q)

	rows, err := a.db.Query(ctx, q, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()
	// "total" column will always be returned, so other columns might be null
	var (
		pID         pgtype.Varchar
		pFName      pgtype.Varchar
		pLName      pgtype.Varchar
		cDate       pgtype.Date
		cLetter     pgtype.Varchar
		gadgetsAggr pgtype.Float4
		paperAggr   pgtype.Float4
		plasticAggr pgtype.Float4
		eID         pgtype.Varchar
		eDate       pgtype.Date
		eName       pgtype.Varchar
		eResAllowed []string
		gadgets     pgtype.Float4
		paper       pgtype.Float4
		plastic     pgtype.Float4
	)
	// map to fast search the pupil in the pupils slice. ID -> index in the slice
	pupilSliceIndex := map[string]int{}
	for rows.Next() {
		err := rows.Scan(&pID, &pFName, &pLName, &cDate, &cLetter, &gadgetsAggr, &paperAggr, &plasticAggr, &eID, &eDate,
			&eName, &eResAllowed, &gadgets, &paper, &plastic, &total)
		if err != nil {
			return nil, 0, err
		}
		// next will happen if the offset >= total rows found or no pupils matching the provided criteria have been
		// found. In that case we simply return total w/o additional work
		if pID.Status != pgtype.Present {
			return nil, total, nil
		}
		var p *aggregating.Pupil
		if i, ok := pupilSliceIndex[string(pID.String)]; ok {
			// if the pupil is already in the pupils, get it
			p = pupils[i]
		} else {
			// otherwise, create it
			p = &aggregating.Pupil{
				Pupil: eventsvc.Pupil{
					ID:        pID.String,
					FirstName: pFName.String,
					LastName:  pLName.String,
				},
				Class: eventsvc.Class{
					Letter:     cLetter.String,
					DateFormed: cDate.Time,
				},
				ResourcesBrought: eventsvc.ResourceMap{
					eventsvc.Gadgets: gadgetsAggr.Float,
					eventsvc.Paper:   paperAggr.Float,
					eventsvc.Plastic: plasticAggr.Float,
				},
			}
			// append the pupil to the slice and put its index to the map
			pupils = append(pupils, p)
			pupilSliceIndex[p.ID] = len(pupils) - 1
		}
		resAllowed, err := eventsvc.StringSliceToResourceSlice(eResAllowed)
		if err != nil {
			return nil, 0, err
		}
		// create a map of the resources brought by the pupil to the event
		resBrought := newResourceMap(resAllowed, gadgets.Float, paper.Float, plastic.Float)
		// create an event and append it to the pupil's slice of events
		e := &aggregating.Event{
			Event: eventsvc.Event{
				ID:               eID.String,
				Date:             eDate.Time,
				Name:             eName.String,
				ResourcesAllowed: resAllowed,
			},
			ResourcesBrought: resBrought,
		}
		p.Events = append(p.Events, e)
	}
	if err = rows.Err(); err != nil {
		return nil, 0, err
	}
	return pupils, total, nil
}

const pupilByIDQueryA = `
	select e.id                              as event_id,
		   e.date,
		   e.name,
		   e.resources_allowed::text[],
		   coalesce(gadgets, 0)              as gadgets,
		   coalesce(paper, 0)                as paper,
		   coalesce(plastic, 0)              as plastic,
		   p.id,
		   p.first_name,
		   p.last_name,
		   p.class_date_formed,
		   p.class_letter,
		   sum(coalesce(gadgets, 0)) over () as gadgets_aggr,
		   sum(coalesce(paper, 0)) over ()   as paper_aggr,
		   sum(coalesce(plastic, 0)) over () as plastic_aggr
	from pupil p
			 left join event e on e.date between symmetric greatest(p.class_date_formed, ?) and least(
				p.class_date_formed + (365.25 * 11)::integer, ?) %s
			 left join resources r on p.id = r.pupil_id and e.id = r.event_id
	where p.id = ?
	order by %s;
`

// PupilByID returns a pupil with the given ID, with a list of all the resources they has brought to every event that
// passed the provided filter
func (a *aggregatingRepo) PupilByID(ctx context.Context, id string, filters aggregating.EventFilters,
	eventsSorting sorting.By) (*aggregating.Pupil, error) {

	// create the "left join event e on" part of the query
	joinOn := strings.Builder{}
	var args []interface{}
	// if filters.To is not set, set it to some date in the distant future
	if filters.To.IsZero() {
		filters.To = filters.To.AddDate(2222, 0, 0)
	}
	args = append(args, filters.From, filters.To)
	if len(filters.ResourcesAllowed) > 0 {
		joinOn.WriteString("and e.resources_allowed @> ?::text[]::resource[] ")
		args = append(args, eventsvc.ResourceSliceToStringSlice(filters.ResourcesAllowed))
	}
	if filters.Name != "" {
		eventTextSearch := pgtextsearch.PrepareQuery(filters.Name)
		joinOn.WriteString("and e.text_search @@ to_tsquery('simple', ?) ")
		args = append(args, eventTextSearch)
	}
	args = append(args, id)

	// get and add "order by" to the query
	orderBy := eventOrderMap[eventsSorting]
	// create the query
	q := fmt.Sprintf(pupilByIDQueryA, joinOn.String(), orderBy)
	// change all the "?" in the query to "$"
	q = sqlx.Rebind(sqlx.BindType("pgx"), q)

	rows, err := a.db.Query(ctx, q, args...)
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
		eGadgets          pgtype.Float4
		ePaper            pgtype.Float4
		ePlastic          pgtype.Float4
		pGadgets          pgtype.Float4
		pPaper            pgtype.Float4
		pPlastic          pgtype.Float4
	)
	p := &aggregating.Pupil{}
	for rows.Next() {
		e := &aggregating.Event{}
		if err := rows.Scan(&eID, &eDate, &eName, &eResourcesAllowed, &eGadgets, &ePaper, &ePlastic, &p.ID,
			&p.FirstName, &p.LastName, &p.Class.DateFormed, &p.Class.Letter, &pGadgets, &pPaper,
			&pPlastic); err != nil {

			return nil, err
		}
		// if event id is null, then no events have been found for the dates passed.
		// We simply return the pupil with an empty event slice
		if eID.Status != pgtype.Present {
			return p, nil
		}
		// set the event fields
		e.ID = eID.String
		e.Date = eDate.Time
		e.Name = eName.String
		resAllowed, err := eventsvc.StringSliceToResourceSlice(eResourcesAllowed)
		resBrought := newResourceMap(resAllowed, eGadgets.Float, ePaper.Float, ePlastic.Float)
		if err != nil {
			return nil, err
		}
		e.ResourcesAllowed = resAllowed
		e.ResourcesBrought = resBrought
		// append the event to the pupil's slice of events
		p.Events = append(p.Events, e)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	// if p.ID isn't set, then no pupil has been found
	if p.ID == "" {
		return nil, eventsvc.ErrUnknownPupil
	}
	// add the resources the pupil brought to the event
	p.ResourcesBrought = eventsvc.ResourceMap{eventsvc.Gadgets: pGadgets.Float, eventsvc.Paper: pPaper.Float,
		eventsvc.Plastic: pPlastic.Float}
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
func (a *aggregatingRepo) Events(ctx context.Context, filters aggregating.EventFilters, sortBy sorting.By, amount,
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
		args = append(args, eventsvc.ResourceSliceToStringSlice(filters.ResourcesAllowed))
	}
	if filters.Name != "" {
		textSearch := pgtextsearch.PrepareQuery(filters.Name)
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
	// "total" column will always be returned, so other columns might be null
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
		// convert []string to []eventsvc.Resource
		resAllowed, err := eventsvc.StringSliceToResourceSlice(resAllowedStr)
		resBrought := newResourceMap(resAllowed, gadgets.Float, paper.Float, plastic.Float)
		if err != nil {
			return nil, 0, err
		}
		e := &aggregating.Event{
			Event: eventsvc.Event{
				ID:               id.String,
				Date:             date.Time,
				Name:             name.String,
				ResourcesAllowed: resAllowed,
			},
			ResourcesBrought: resBrought,
		}
		events = append(events, e)
	}
	if err = rows.Err(); err != nil {
		return nil, 0, err
	}
	return events, total, nil
}

func createClassID(date time.Time, letter string) string {
	return fmt.Sprintf("%d%s", date.Year(), letter)
}

func newResourceMap(resAllowed []eventsvc.Resource, gadgets float32, paper float32,
	plastic float32) eventsvc.ResourceMap {

	resBrought := eventsvc.ResourceMap{}
	for _, res := range resAllowed {
		switch res {
		case eventsvc.Gadgets:
			resBrought[res] = gadgets
		case eventsvc.Paper:
			resBrought[res] = paper
		case eventsvc.Plastic:
			resBrought[res] = plastic
		}
	}
	return resBrought
}
