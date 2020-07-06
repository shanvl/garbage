package postgres

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/jmoiron/sqlx"
	"github.com/shanvl/garbage/internal/eventsvc"
	"github.com/shanvl/garbage/internal/eventsvc/schooling"
)

type schoolingRepo struct {
	db *pgxpool.Pool
}

func NewSchoolingRepo(db *pgxpool.Pool) schooling.Repository {
	return &schoolingRepo{db}
}

const pupilByIDQuery = `
	select id, first_name, last_name, class_letter, class_date_formed
	from pupil
	where pupil.id = $1;
`

// returns pupils with the given id
func (s *schoolingRepo) PupilByID(ctx context.Context, pupilID string) (*schooling.Pupil, error) {
	p := &schooling.Pupil{}
	err := s.db.QueryRow(ctx, pupilByIDQuery, pupilID).Scan(&p.ID, &p.FirstName, &p.LastName, &p.Class.Letter,
		&p.Class.DateFormed)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			err = eventsvc.ErrUnknownPupil
		}
		return nil, err
	}
	return p, nil
}

// "?" instead of "$" here because the query will be modified by sqlx.In
const removePupilsQuery = `delete from pupil where id in(?)`

// removes pupils with the given ids
func (s *schoolingRepo) RemovePupils(ctx context.Context, pupilIDs []string) error {
	// create the query with some magic of sqlx
	q, args, err := sqlx.In(removePupilsQuery, pupilIDs)
	if err != nil {
		return err
	}
	q = sqlx.Rebind(sqlx.BindType("pgx"), q)
	// execute the query
	if _, err := s.db.Exec(ctx, q, args...); err != nil {
		return err
	}
	// if there's no error, all the passed pupils have been removed, so their ids can be returned
	return nil
}

const storePupilsQuery = `insert into pupil (id, first_name, last_name, class_letter, class_date_formed) values`

// saves the given pupils
func (s *schoolingRepo) StorePupils(ctx context.Context, pupils []*schooling.Pupil) error {
	pupilsLen := len(pupils)
	// params placeholders to pass to the query
	queryParams := make([]string, pupilsLen)
	// values to pass to the query
	queryValues := make([]interface{}, 0, pupilsLen*5)
	// ids of the saved pupils
	ids := make([]string, pupilsLen)
	for i, p := range pupils {
		n := i * 5
		// create query params placeholders for the pupil
		qp := fmt.Sprintf("($%d, $%d, $%d, $%d, $%d)", n+1, n+2, n+3, n+4, n+5)
		queryParams[i] = qp
		// push param values
		queryValues = append(queryValues, p.ID, p.FirstName, p.LastName, p.Class.Letter, p.Class.DateFormed)
		// push the pupil's id
		ids[i] = p.ID
	}
	// create and execute the query
	q := fmt.Sprintf("%s %s;", storePupilsQuery, strings.Join(queryParams, ","))
	_, err := s.db.Exec(ctx, q, queryValues...)
	return err
}

const updatePupilQuery = `
	update pupil 
	set (first_name, last_name, class_letter, class_date_formed) = ($2, $3, $4, $5)
	where id = $1
	returning id;
`

// updates the given pupil
func (s *schoolingRepo) UpdatePupil(ctx context.Context, pupil *schooling.Pupil) error {
	var id string
	err := s.db.QueryRow(ctx, updatePupilQuery, pupil.ID, pupil.FirstName, pupil.LastName, pupil.Class.Letter,
		pupil.Class.DateFormed).Scan(&id)
	if err != nil {
		return err
	}
	// no pupil has been found
	if id == "" {
		return eventsvc.ErrUnknownPupil
	}
	return nil
}
