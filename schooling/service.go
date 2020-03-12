package schooling

import (
	"context"

	"github.com/shanvl/garbage-events-service/garbage"
	"github.com/shanvl/garbage-events-service/valid"
)

type Service interface {
	AddPupils(ctx context.Context, pupilInfo []PupilInfo) ([]garbage.PupilID, error)
	ChangePupilClass(ctx context.Context, pupilID garbage.PupilID, className string) (garbage.PupilID,
		garbage.ClassID, error)
	// RemovePupils removes pupils using provided IDs and returns their IDs
	RemovePupils(ctx context.Context, pupilIDs []garbage.PupilID) ([]garbage.PupilID, error)
}

type Repository interface {
	AddPupils(ctx context.Context, pupils []*Pupil) ([]garbage.PupilID, error)
	ChangePupilClass(ctx context.Context, pupilID garbage.PupilID, className string) (garbage.PupilID,
		garbage.ClassID, error)
	RemovePupils(ctx context.Context, pupilIDs []garbage.PupilID) ([]garbage.PupilID, error)
}

type service struct {
	repo Repository
}

// NewService returns an instance of Service with all its dependencies
func NewService(repo Repository) Service {
	return &service{repo}
}

func (s *service) AddPupils(ctx context.Context, pupilInfo []PupilInfo) ([]garbage.PupilID, error) {
	panic("implement me")
}

func (s *service) ChangePupilClass(ctx context.Context, pupilID garbage.PupilID, className string) (garbage.PupilID,
	garbage.ClassID,
	error) {
	panic("implement me")
}

// RemovePupils removes pupils using provided IDs and returns their IDs
func (s *service) RemovePupils(ctx context.Context, pupilIDs []garbage.PupilID) ([]garbage.PupilID, error) {
	errVld := valid.EmptyError()
	if len(pupilIDs) <= 0 {
		errVld.Add("pupilIDs", "no pupil ids were provided")
	}
	for _, id := range pupilIDs {
		if len(id) <= 0 {
			errVld.Add("pupilIDs", "pupil id can't be empty")
		}
	}
	if !errVld.IsEmpty() {
		return nil, errVld
	}
	pupilIDs, err := s.repo.RemovePupils(ctx, pupilIDs)
	if err != nil {
		return nil, err
	}
	return pupilIDs, nil
}

// PupilInfo used in adding new pupils
type PupilInfo struct {
	FirstName, LastName, ClassName string
}

// Pupil is a model of the pupil, adapted for this use case.
type Pupil struct {
	garbage.Pupil
	Class garbage.Class
}
