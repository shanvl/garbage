// Package schooling is responsible for school management.
// There are no concepts like events or resources, only pupils and classes
package schooling

import (
	"context"
	"fmt"
	"time"

	"github.com/shanvl/garbage-events-service/internal/garbage"
	"github.com/shanvl/garbage-events-service/internal/idgen"
	"github.com/shanvl/garbage-events-service/internal/valid"
)

// Service is an interface providing methods to manage pupils and classes w/o concepts like events or resources
type Service interface {
	// AddPupils adds pupils returning ids of added
	AddPupils(ctx context.Context, pupilInfo []PupilBio) ([]garbage.PupilID, error)
	// ChangePupilClass changes the class of the pupil if such a class exists
	ChangePupilClass(ctx context.Context, pupilID garbage.PupilID, className string) (garbage.PupilID, error)
	// RemovePupils removes pupils using provided IDs and returns their IDs
	RemovePupils(ctx context.Context, pupilIDs []garbage.PupilID) ([]garbage.PupilID, error)
}

type Repository interface {
	PupilByID(ctx context.Context, pupilID garbage.PupilID) (*Pupil, error)
	RemovePupils(ctx context.Context, pupilIDs []garbage.PupilID) ([]garbage.PupilID, error)
	StorePupil(ctx context.Context, pupil *Pupil) (garbage.PupilID, error)
	StorePupils(ctx context.Context, pupils []*Pupil) ([]garbage.PupilID, error)
}

type service struct {
	repo Repository
}

const (
	MaxAddPupils    = 1000
	MaxRemovePupils = 1000
)

// NewService returns an instance of Service with all its dependencies
func NewService(repo Repository) Service {
	return &service{repo}
}

// AddPupils adds pupils, returning ids of added
func (s *service) AddPupils(ctx context.Context, pupilsBio []PupilBio) ([]garbage.PupilID, error) {
	if len(pupilsBio) == 0 {
		return nil, valid.NewError("pupils", "no pupils were provided")
	}
	if len(pupilsBio) > MaxAddPupils {
		return nil, valid.NewError("pupils", fmt.Sprintf("no more than %d pupils are allowed to add", MaxAddPupils))
	}
	// pupils to pass to the repo
	pupils := make([]*Pupil, 0, len(pupilsBio))
	// date needed to derive a pupil's class entity out of its class name
	today := time.Now()

	errVld := valid.EmptyError()
	for i, bio := range pupilsBio {
		// validate a pupil's name and class
		f := fmt.Sprintf("pupils[%d]", i)
		if len(bio.FirstName) == 0 {
			errVld.Add(fmt.Sprintf("%s[firstName]", f), "first name must be provided")
		}
		if len(bio.LastName) == 0 {
			errVld.Add(fmt.Sprintf("%s[lastName]", f), "last name must be provided")
		}
		if len(bio.ClassName) == 0 {
			errVld.Add(fmt.Sprintf("%s[class]", f), "class must be provided")
		}
		// derive the classLetter and classYearFormed from the class name
		classLetter, classYearFormed, err := garbage.ParseClassName(bio.ClassName, today)
		// if invalid className, add it to validation error and go on to the next pupil
		// in order to collect all validation errors
		if err != nil {
			errVld.Add(fmt.Sprintf("%s[class]", f), fmt.Sprintf("invalid class name: %s", bio.ClassName))
			continue
		}
		// create a pupil entity
		pupilID, err := idgen.CreatePupilID()
		if err != nil {
			return nil, err
		}
		p := &Pupil{
			Pupil: garbage.Pupil{
				ID:        pupilID,
				FirstName: bio.FirstName,
				LastName:  bio.LastName,
			},
			Class: garbage.Class{Letter: classLetter, YearFormed: classYearFormed},
		}
		// push the pupil entity to the pupils slice
		pupils = append(pupils, p)
	}
	// if there are validation errors, return them w/o proceeding further
	if !errVld.IsEmpty() {
		return nil, errVld
	}

	// save pupils
	return s.repo.StorePupils(ctx, pupils)
}

// ChangePupilClass changes the class of the pupil
func (s *service) ChangePupilClass(ctx context.Context, pupilID garbage.PupilID, className string) (garbage.PupilID,
	error) {

	// validate args
	errVld := valid.EmptyError()
	if len(pupilID) == 0 {
		errVld.Add("pupilID", "pupilID must be provided")
	}
	if len(className) == 0 {
		errVld.Add("className", "className must be provided")
	}
	if !errVld.IsEmpty() {
		return "", errVld
	}

	// get pupil
	pupil, err := s.repo.PupilByID(ctx, pupilID)
	if err != nil {
		return "", err
	}
	// parse className
	classLetter, classYearFormed, err := garbage.ParseClassName(className, time.Now())
	if err != nil {
		return "", err
	}
	// if the pupil is already in the class, return their id
	if classLetter == pupil.Class.Letter && classYearFormed == pupil.Class.YearFormed {
		return pupil.ID, nil
	}
	// otherwise, change the class' data
	pupil.Class.Letter, pupil.Class.YearFormed = classLetter, classYearFormed
	// save the pupil
	pupilID, err = s.repo.StorePupil(ctx, pupil)
	if err != nil {
		return "", err
	}
	return pupilID, nil
}

// RemovePupils removes pupils using provided IDs and returns their IDs
func (s *service) RemovePupils(ctx context.Context, pupilIDs []garbage.PupilID) ([]garbage.PupilID, error) {
	if len(pupilIDs) == 0 {
		return nil, valid.NewError("pupils", "no pupils ids were provided")
	}
	if len(pupilIDs) > MaxRemovePupils {
		return nil, valid.NewError("pupils", fmt.Sprintf("no more than %d pupils can be deleted in one run",
			MaxRemovePupils))
	}
	// validate pupils ids
	errVld := valid.EmptyError()
	for i, id := range pupilIDs {
		errField := fmt.Sprintf("pupils[%d]", i)
		if len(id) == 0 {
			errVld.Add(errField, "pupil id can't be empty")
		}
	}
	if !errVld.IsEmpty() {
		return nil, errVld
	}

	return s.repo.RemovePupils(ctx, pupilIDs)
}

// PupilBio used in adding new pupils
type PupilBio struct {
	FirstName, LastName, ClassName string
}

// Pupil is a model of the pupil, adapted for this use case.
type Pupil struct {
	garbage.Pupil
	Class garbage.Class
}
