package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"github.com/jmoiron/sqlx"
	"github.com/shanvl/garbage-events-service/internal/garbage"
	"github.com/shanvl/garbage-events-service/internal/usecases/schooling"
)

type SchoolingRepo struct {
	db *sqlx.DB
}

func NewSchoolingRepo(db *sqlx.DB) *SchoolingRepo {
	return &SchoolingRepo{db}
}

const pupilByIDQuery = `
	select id, first_name, last_name, class_letter, class_year_formed
	from pupil
	where pupil.id = $1;
`

// returns pupils with the given id
func (s *SchoolingRepo) PupilByID(ctx context.Context, pupilID garbage.PupilID) (*schooling.Pupil, error) {
	p := &schooling.Pupil{}
	err := s.db.QueryRowContext(ctx, pupilByIDQuery, pupilID).Scan(&p.ID, &p.FirstName, &p.LastName, &p.Class.Letter,
		&p.Class.YearFormed)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, garbage.ErrNoPupil
		}
		return nil, fmt.Errorf("error finding a pupil by id: %w", err)
	}
	return p, nil
}

// "?" instead of "$" here because the query will be modified by sqlx.In
const removePupilsQuery = `delete from pupil where id in(?)`

// removes pupils with the given ids
func (s *SchoolingRepo) RemovePupils(ctx context.Context, pupilIDs []garbage.PupilID) error {
	// create the query with some magic of sqlx
	q, args, err := sqlx.In(removePupilsQuery, pupilIDs)
	if err != nil {
		return err
	}
	q = s.db.Rebind(q)
	// execute the query
	if _, err := s.db.ExecContext(ctx, q, args...); err != nil {
		return err
	}
	// if there's no error, all the passed pupils have been removed, so their ids can be returned
	return nil
}

const storePupilQuery = `
	insert into pupil (id, first_name, last_name, class_letter, class_year_formed)
	values ($1, $2, $3, $4, $5);
`

// saves the given pupil
func (s *SchoolingRepo) StorePupil(ctx context.Context, pupil *schooling.Pupil) error {
	_, err := s.db.ExecContext(ctx, storePupilQuery, pupil.ID, pupil.FirstName, pupil.LastName, pupil.Class.Letter,
		pupil.Class.YearFormed)
	return err
}

const storePupilsQuery = `insert into pupil (id, first_name, last_name, class_letter, class_year_formed) values`

// saves the given pupils
func (s *SchoolingRepo) StorePupils(ctx context.Context, pupils []*schooling.Pupil) error {
	pupilsLen := len(pupils)
	// params placeholders to pass to the query
	queryParams := make([]string, pupilsLen)
	// values to pass to the query
	queryValues := make([]interface{}, 0, pupilsLen*5)
	// ids of the saved pupils
	ids := make([]garbage.PupilID, pupilsLen)
	for i, p := range pupils {
		n := i * 5
		// create query params placeholders for the pupil
		qp := fmt.Sprintf("($%d, $%d, $%d, $%d, $%d)", n+1, n+2, n+3, n+4, n+5)
		queryParams[i] = qp
		// push param values
		queryValues = append(queryValues, p.ID, p.FirstName, p.LastName, p.Class.Letter, p.Class.YearFormed)
		// push the pupil's id
		ids[i] = p.ID
	}
	// create and execute the query
	q := fmt.Sprintf("%s %s;", storePupilsQuery, strings.Join(queryParams, ","))
	_, err := s.db.ExecContext(ctx, q, queryValues...)
	return err
}
