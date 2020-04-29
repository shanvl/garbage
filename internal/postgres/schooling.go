package postgres

import (
	"context"
	"database/sql"
	"fmt"

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

func (s *SchoolingRepo) RemovePupils(ctx context.Context, pupilIDs []garbage.PupilID) ([]garbage.PupilID, error) {
	panic("implement me")
}

const storePupilQuery = `
	insert into pupil (id, first_name, last_name, class_letter, class_year_formed)
	values ($1, $2, $3, $4, $5)
	returning id;
`

func (s *SchoolingRepo) StorePupil(ctx context.Context, pupil *schooling.Pupil) (garbage.PupilID, error) {
	var id garbage.PupilID
	err := s.db.QueryRowContext(ctx, storePupilQuery, pupil.ID, pupil.FirstName, pupil.LastName, pupil.Class.Letter,
		pupil.Class.YearFormed).Scan(&id)
	if err != nil {
		return "", err
	}
	return id, nil
}

func (s *SchoolingRepo) StorePupils(ctx context.Context, pupils []*schooling.Pupil) ([]garbage.PupilID, error) {
	panic("implement me")
}
