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
	ChangePupilClass(ctx context.Context, pupilID garbage.PupilID, className string) (garbage.PupilID,
		garbage.ClassID, error)
	// RemovePupils removes pupils using provided IDs and returns their IDs
	RemovePupils(ctx context.Context, pupilIDs []garbage.PupilID) ([]garbage.PupilID, error)
}

type Repository interface {
	Class(ctx context.Context, letter string, yearFormed int) (*garbage.Class, error)
	PupilByID(ctx context.Context, pupilID garbage.PupilID) (*Pupil, error)
	RemovePupils(ctx context.Context, pupilIDs []garbage.PupilID) ([]garbage.PupilID, error)
	StorePupil(ctx context.Context, pupil *Pupil) (garbage.PupilID, error)
	// transaction
	StorePupils(ctx context.Context, pupils []*Pupil) ([]garbage.PupilID, error)
}

type service struct {
	repo Repository
}

// NewService returns an instance of Service with all its dependencies
func NewService(repo Repository) Service {
	return &service{repo}
}

// AddPupils adds pupils returning ids of added
func (s *service) AddPupils(ctx context.Context, pupilsBio []PupilBio) ([]garbage.PupilID, error) {
	// pupils to pass to the repo
	pupils := make([]*Pupil, 0, len(pupilsBio))
	// map of seen classes
	classesCache := map[string]*garbage.Class{}
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
		// get a pupil's class entity
		class, ok := classesCache[bio.ClassName]
		// if there's no class in the cache, get or create a new one
		if !ok {
			letter, yearFormed, err := garbage.ParseClassName(bio.ClassName, today)
			// if invalid className, add it to validation error and go on to the next pupil
			// in order to collect all validation errors
			if err != nil {
				errVld.Add(fmt.Sprintf("%s[class]", f), fmt.Sprintf("invalid class name: %s", bio.ClassName))
				continue
			}
			// check if such a class already exists
			class, err = s.repo.Class(ctx, letter, yearFormed)
			if err != nil && err != garbage.ErrNoClass {
				return nil, err
			}
			// if no class exists, create one
			if err == garbage.ErrNoClass {
				classID, err := idgen.CreateClassID()
				if err != nil {
					return nil, err
				}
				class = &garbage.Class{
					ID:         classID,
					YearFormed: yearFormed,
					Letter:     letter,
				}
			}
			// add it to the cache
			classesCache[bio.ClassName] = class
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
			Class: *class,
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

// ChangePupilClass changes the class of the pupil if such a class exists
func (s *service) ChangePupilClass(ctx context.Context, pupilID garbage.PupilID, className string) (garbage.PupilID,
	garbage.ClassID, error) {

	// validate args
	errVld := valid.EmptyError()
	if len(pupilID) == 0 {
		errVld.Add("pupilID", "pupilID must be provided")
	}
	if len(className) == 0 {
		errVld.Add("className", "className must be provided")
	}
	if !errVld.IsEmpty() {
		return "", "", errVld
	}

	// get pupil
	pupil, err := s.repo.PupilByID(ctx, pupilID)
	if err != nil {
		return "", "", err
	}
	// parse className
	letter, yearFormed, err := garbage.ParseClassName(className, time.Now())
	if err != nil {
		return "", "", err
	}
	// if the pupil is already in the class, return their ids
	if letter == pupil.Class.Letter && yearFormed == pupil.Class.YearFormed {
		return pupil.ID, pupil.Class.ID, nil
	}
	// otherwise, find such a class
	class, err := s.repo.Class(ctx, letter, yearFormed)
	if err != nil {
		return "", "", err
	}
	// change class in the pupil entity
	pupil.Class = *class
	// save the pupil
	pupilID, err = s.repo.StorePupil(ctx, pupil)
	if err != nil {
		return "", "", err
	}
	return pupilID, class.ID, nil
}

// RemovePupils removes pupils using provided IDs and returns their IDs
func (s *service) RemovePupils(ctx context.Context, pupilIDs []garbage.PupilID) ([]garbage.PupilID, error) {
	// validate pupils ids
	errVld := valid.EmptyError()
	if len(pupilIDs) == 0 {
		errVld.Add("pupils", "no pupils ids were provided")
	}
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
