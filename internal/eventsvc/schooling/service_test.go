package schooling_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/shanvl/garbage/internal/eventsvc"
	"github.com/shanvl/garbage/internal/eventsvc/mock"
	"github.com/shanvl/garbage/internal/eventsvc/schooling"
)

func Test_service_RemovePupils(t *testing.T) {
	ctx := context.Background()
	ids := []string{"000", "001", "002"}

	var repo mock.SchoolingRepository
	repo.RemovePupilsFn = func(ctx context.Context, pupilIDs []string) error {
		if len(pupilIDs) > 0 && pupilIDs[0] == "error" {
			return errors.New("some error")
		}
		return nil
	}
	s := schooling.NewService(&repo)

	type args struct {
		pupilIDs []string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "no ids",
			args: args{
				pupilIDs: nil,
			},
			wantErr: true,
		},
		{
			name: "too much pupil ids",
			args: args{
				pupilIDs: make([]string, schooling.MaxRemovePupils+1),
			},
			wantErr: true,
		},
		{
			name: "one id is empty",
			args: args{
				pupilIDs: []string{"123", ""},
			},
			wantErr: true,
		},
		{
			name: "repo's error",
			args: args{
				pupilIDs: []string{"error"},
			},
			wantErr: true,
		},
		{
			name: "remove 10 pupils",
			args: args{
				pupilIDs: ids,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := s.RemovePupils(ctx, tt.args.pupilIDs)
			if (err != nil) != tt.wantErr {
				t.Errorf("RemovePupils() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func Test_service_AddPupils(t *testing.T) {
	ctx := context.Background()

	var repo mock.SchoolingRepository
	repo.StorePupilsFn = func(ctx context.Context, pupils []*schooling.Pupil) error {
		if len(pupils) > 0 && pupils[0].FirstName == "error" {
			return errors.New("repo's error")
		}
		return nil
	}
	s := schooling.NewService(&repo)

	tests := []struct {
		name    string
		pupils  []schooling.PupilBio
		wantLen int
		wantErr bool
	}{
		{
			name:    "no pupils",
			pupils:  nil,
			wantLen: 0,
			wantErr: true,
		},
		{
			name:    "too much pupils",
			pupils:  make([]schooling.PupilBio, schooling.MaxAddPupils+1),
			wantLen: 0,
			wantErr: true,
		},
		{
			name:    "repo.StorePupils error",
			pupils:  []schooling.PupilBio{{"error", "ln", "8B"}},
			wantLen: 0,
			wantErr: true,
		},
		{
			name:    "empty name",
			pupils:  []schooling.PupilBio{{"", "ln", "8B"}},
			wantLen: 0,
			wantErr: true,
		},
		{
			name: "first name is too long",
			pupils: []schooling.PupilBio{{"zzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzz", "ln",
				"8B"}},
			wantLen: 0,
			wantErr: true,
		},
		{
			name:    "empty last name",
			pupils:  []schooling.PupilBio{{"fn", "", "8B"}},
			wantLen: 0,
			wantErr: true,
		},
		{
			name: "last name is too long",
			pupils: []schooling.PupilBio{{"fn", "zzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzz",
				"8B"}},
			wantLen: 0,
			wantErr: true,
		},
		{
			name:    "invalid className",
			pupils:  []schooling.PupilBio{{"fn", "ln", "B8B"}},
			wantLen: 0,
			wantErr: true,
		},
		{
			name: "3 valid pupils",
			pupils: []schooling.PupilBio{
				{
					FirstName: "fn",
					LastName:  "ln",
					ClassName: "8B",
				}, {
					FirstName: "fn",
					LastName:  "ln",
					ClassName: "9B",
				}, {
					FirstName: "aa",
					LastName:  "bb",
					ClassName: "9B",
				}},
			wantLen: 3,
			wantErr: false,
		},
		{
			name: "no error on duplicates",
			pupils: []schooling.PupilBio{
				{
					FirstName: "fn",
					LastName:  "ln",
					ClassName: "8B",
				},
				{
					FirstName: "fn",
					LastName:  "ln",
					ClassName: "8B",
				},
				{
					FirstName: "aa",
					LastName:  "bb",
					ClassName: "8B",
				},
			},
			wantLen: 3,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := s.AddPupils(ctx, tt.pupils)
			gotLen := len(got)
			if (err != nil) != tt.wantErr {
				t.Errorf("StorePupils() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotLen != tt.wantLen {
				t.Errorf("StorePupils() gotLen = %v, wantLen %v", got, tt.wantLen)
			}
		})
	}
}

func Test_service_ChangePupilClass(t *testing.T) {
	const (
		pupilIDPupilNotFound           = "pupil not found"
		pupilIDStoreError              = "store error"
		pupilIDPupilWithRequestedClass = "111"
		requestedClassName             = "10B"
	)
	ctx := context.Background()

	var repo mock.SchoolingRepository
	repo.UpdatePupilFn = func(ctx context.Context, pupil *schooling.Pupil) error {
		if pupil.ID == pupilIDStoreError {
			return errors.New("repo's error")
		}
		return nil
	}
	repo.PupilByIDFn = func(ctx context.Context, pupilID string) (pupil *schooling.Pupil, err error) {
		if pupilID == pupilIDPupilNotFound {
			return nil, errors.New("pupil not found error")
		}
		return &schooling.Pupil{
			Pupil: eventsvc.Pupil{
				ID:        pupilID,
				FirstName: "FN",
				LastName:  "LN",
			},
			Class: eventsvc.Class{
				Letter:     "Y",
				DateFormed: time.Date(2010, 9, 1, 0, 0, 0, 0, time.UTC),
			},
		}, nil
	}
	s := schooling.NewService(&repo)

	type args struct {
		pupilID   string
		className string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "no pupilID",
			args: args{
				pupilID:   "",
				className: "className",
			},
			wantErr: true,
		},
		{
			name: "no className",
			args: args{
				pupilID:   "pupilID",
				className: "",
			},
			wantErr: true,
		},
		{
			name: "invalid className",
			args: args{
				pupilID:   "pupilID",
				className: "12B",
			},
			wantErr: true,
		},
		{
			name: "invalid className2",
			args: args{
				pupilID:   "pupilID",
				className: "10 BB",
			},
			wantErr: true,
		},
		{
			name: "no pupil with the given id",
			args: args{
				pupilID:   pupilIDPupilNotFound,
				className: "class name",
			},
			wantErr: true,
		},
		{
			name: "error on storing a pupil",
			args: args{
				pupilID:   pupilIDStoreError,
				className: "10B",
			},
			wantErr: true,
		},
		{
			name: "pupil is already in the given class",
			args: args{
				pupilID:   pupilIDPupilWithRequestedClass,
				className: requestedClassName,
			},
			wantErr: false,
		},
		{
			name: "swap a pupil's class for the found one",
			args: args{
				pupilID:   "pupilID",
				className: "10Y",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := s.ChangePupilClass(ctx, tt.args.pupilID, tt.args.className)
			if (err != nil) != tt.wantErr {
				t.Errorf("ChangePupilClass() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}
