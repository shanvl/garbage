package schooling

import (
	"context"
	"fmt"
	"time"

	"github.com/shanvl/garbage-events-service/garbage"
	"github.com/shanvl/garbage-events-service/idgen"
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
	AddPupils(ctx context.Context, pupils []Pupil) ([]garbage.PupilID, error)
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

// AddPupils adds pupils, returning ids of added
func (s *service) AddPupils(ctx context.Context, pupilsInfo []PupilInfo) ([]garbage.PupilID, error) {
	// pupils to pass to the repo
	pupils := make([]Pupil, 0, len(pupilsInfo))
	// time needed to derive a pupil's class entity out of its class name
	today := time.Now()
	// map of seen classes
	classesCache := map[string]garbage.Class{}

	errVld := valid.EmptyError()
	for i, info := range pupilsInfo {
		// validate a pupil's name and class
		f := fmt.Sprintf("pupils[%d]", i)
		if len(info.FirstName) <= 0 {
			errVld.Add(fmt.Sprintf("%s[firstName]", f), "first name must be provided")
		}
		if len(info.LastName) <= 0 {
			errVld.Add(fmt.Sprintf("%s[lastName]", f), "last name must be provided")
		}
		if len(info.ClassName) <= 0 {
			errVld.Add(fmt.Sprintf("%s[class]", f), "class must be provided")
		}
		// get a pupil's class entity
		class, ok := classesCache[info.ClassName]
		// if there's no class in the cache, create a new one
		if !ok {
			letter, yearFormed, err := garbage.ParseClassName(info.ClassName, today)
			if err != nil {
				errVld.Add(fmt.Sprintf("%s[class]", f), fmt.Sprintf("invalid class name: %s", info.ClassName))
			}
			classID, err := idgen.CreateClassID()
			if err != nil {
				return nil, err
			}
			class = garbage.Class{
				ID:         classID,
				YearFormed: yearFormed,
				Letter:     letter,
			}
			// add it to the cache
			classesCache[info.ClassName] = class
		}
		// create a pupil entity
		pupilID, err := idgen.CreatePupilID()
		if err != nil {
			return nil, err
		}
		p := Pupil{
			Pupil: garbage.Pupil{
				ID:        pupilID,
				FirstName: info.FirstName,
				LastName:  info.LastName,
			},
			Class: class,
		}
		// push the pupil entity to the pupils slice
		pupils = append(pupils, p)
	}
	if !errVld.IsEmpty() {
		return nil, errVld
	}
	// repo.AddPupils doesn't overwrite duplicate pupils or classes if it finds them. It moves on silently,
	// raising no errors
	pupilIDS, err := s.repo.AddPupils(ctx, pupils)
	if err != nil {
		return nil, err
	}
	return pupilIDS, nil
}

func (s *service) ChangePupilClass(ctx context.Context, pupilID garbage.PupilID, className string) (garbage.PupilID,
	garbage.ClassID,
	error) {
	panic("implement me")
}

// RemovePupils removes pupils using provided IDs and returns their IDs
func (s *service) RemovePupils(ctx context.Context, pupilIDs []garbage.PupilID) ([]garbage.PupilID, error) {

	// validate pupils ids
	errVld := valid.EmptyError()
	if len(pupilIDs) <= 0 {
		errVld.Add("pupils", "no pupils ids were provided")
	}
	for i, id := range pupilIDs {
		errField := fmt.Sprintf("pupils[%d]", i)
		if len(id) <= 0 {
			errVld.Add(errField, "pupil id can't be empty")
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
