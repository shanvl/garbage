// Package schooling is responsible for the school management such as adding or deleting pupils or changing their
// classes
package schooling

import (
	"context"
	"fmt"
	"time"

	gonanoid "github.com/matoous/go-nanoid"
	"github.com/shanvl/garbage/internal/eventsvc"
	"github.com/shanvl/garbage/pkg/valid"
)

// Service is an interface providing methods to manage pupils and classes w/o concepts like events or resources
type Service interface {
	// AddPupils adds the given pupils and returns ids of the added
	AddPupils(ctx context.Context, pupilInfo []PupilBio) ([]string, error)
	// ChangePupilClass changes the class of the pupil if such a class exists
	ChangePupilClass(ctx context.Context, pupilID string, className string) error
	// RemovePupils removes pupils with provided IDs
	RemovePupils(ctx context.Context, pupilIDs []string) error
}

type Repository interface {
	PupilByID(ctx context.Context, pupilID string) (*Pupil, error)
	RemovePupils(ctx context.Context, pupilIDs []string) error
	StorePupils(ctx context.Context, pupils []*Pupil) error
	UpdatePupil(ctx context.Context, pupil *Pupil) error
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
func (s *service) AddPupils(ctx context.Context, pupilsBio []PupilBio) ([]string, error) {
	if len(pupilsBio) == 0 {
		return nil, valid.NewError("pupils", "no pupils were provided")
	}
	if len(pupilsBio) > MaxAddPupils {
		return nil, valid.NewError("pupils", fmt.Sprintf("no more than %d pupils are allowed to add", MaxAddPupils))
	}
	// pupils to pass to the repo
	pupils := make([]*Pupil, 0, len(pupilsBio))
	pupilIDs := make([]string, 0, len(pupilsBio))
	// date needed to derive a pupil's class entity out of its class name
	today := time.Now()

	errVld := valid.EmptyError()
	for i, bio := range pupilsBio {
		// validate a pupil's name and class
		f := fmt.Sprintf("pupils[%d]", i)
		if len(bio.FirstName) == 0 {
			errVld.Add(fmt.Sprintf("%s[firstName]", f), "first name must be provided")
		}
		if len(bio.FirstName) > 35 {
			errVld.Add(fmt.Sprintf("%s[firstName]", f), "length of the first name can't be more than 35")
		}
		if len(bio.LastName) == 0 {
			errVld.Add(fmt.Sprintf("%s[lastName]", f), "last name must be provided")
		}
		if len(bio.LastName) > 35 {
			errVld.Add(fmt.Sprintf("%s[lastName]", f), "length of the last name can't be more than 35")
		}
		if len(bio.ClassName) == 0 {
			errVld.Add(fmt.Sprintf("%s[class]", f), "class must be provided")
		}
		// derive the classLetter and classYearFormed from the class name
		class, err := eventsvc.ClassFromClassName(bio.ClassName, today)
		// if invalid className, add it to validation error and go on to the next pupil
		// in order to collect all validation errors
		if err != nil {
			errVld.Add(fmt.Sprintf("%s[class]", f), err.Error())
			continue
		}
		// create a pupil entity
		pupilID, err := gonanoid.Nanoid(14)
		if err != nil {
			return nil, fmt.Errorf("couldn't generate pupil id: %w", err)
		}
		p := &Pupil{
			Pupil: eventsvc.Pupil{
				ID:        pupilID,
				FirstName: bio.FirstName,
				LastName:  bio.LastName,
			},
			Class: class,
		}
		// push the pupil entity to the pupils slice
		pupils = append(pupils, p)
		// push the pupil's id to to the slice of pupil's ids
		pupilIDs = append(pupilIDs, pupilID)
	}
	// if there are validation errors, return them w/o proceeding further
	if !errVld.IsEmpty() {
		return nil, errVld
	}

	// save the pupils
	err := s.repo.StorePupils(ctx, pupils)
	if err != nil {
		return nil, err
	}
	return pupilIDs, nil
}

// ChangePupilClass changes the class of the pupil
func (s *service) ChangePupilClass(ctx context.Context, pupilID string, className string) error {

	// validate args
	errVld := valid.EmptyError()
	if len(pupilID) == 0 {
		errVld.Add("pupilID", "pupilID must be provided")
	}
	if len(className) == 0 {
		errVld.Add("class", "class must be provided")
	}
	if !errVld.IsEmpty() {
		return errVld
	}

	// get pupil
	pupil, err := s.repo.PupilByID(ctx, pupilID)
	if err != nil {
		return err
	}
	// parse className
	class, err := eventsvc.ClassFromClassName(className, time.Now())
	if err != nil {
		return valid.NewError("class", err.Error())
	}
	// if the pupil is already in the class, return
	if class == pupil.Class {
		return nil
	}
	// otherwise, change the pupil's class
	pupil.Class = class
	// update the pupil
	return s.repo.UpdatePupil(ctx, pupil)
}

// RemovePupils removes pupils using provided IDs and returns their IDs
func (s *service) RemovePupils(ctx context.Context, pupilIDs []string) error {
	if len(pupilIDs) == 0 {
		return valid.NewError("pupils", "no pupils ids were provided")
	}
	if len(pupilIDs) > MaxRemovePupils {
		return valid.NewError("pupils", fmt.Sprintf("no more than %d pupils can be deleted in one run",
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
		return errVld
	}

	return s.repo.RemovePupils(ctx, pupilIDs)
}

// PupilBio used in adding new pupils
type PupilBio struct {
	FirstName, LastName, ClassName string
}

// Pupil is a model of the pupil, adapted for this use case.
type Pupil struct {
	eventsvc.Pupil
	Class eventsvc.Class
}
