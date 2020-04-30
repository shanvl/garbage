package schooling_test

import (
	"context"
	"errors"
	"testing"

	"github.com/shanvl/garbage-events-service/internal/garbage"
	"github.com/shanvl/garbage-events-service/internal/mock"
	"github.com/shanvl/garbage-events-service/internal/usecases/schooling"
)

func Test_service_RemovePupils(t *testing.T) {
	ctx := context.Background()
	ids := []garbage.PupilID{"000", "001", "002"}

	var repo mock.SchoolingRepository
	repo.RemovePupilsFn = func(ctx context.Context, pupilIDs []garbage.PupilID) ([]garbage.PupilID, error) {
		if len(pupilIDs) > 0 && pupilIDs[0] == "error" {
			return nil, errors.New("some error")
		}
		return ids, nil
	}
	s := schooling.NewService(&repo)

	type args struct {
		pupilIDs []garbage.PupilID
	}
	tests := []struct {
		name       string
		args       args
		wantIDsLen int
		wantErr    bool
	}{
		{
			name: "no ids",
			args: args{
				pupilIDs: nil,
			},
			wantIDsLen: 0,
			wantErr:    true,
		},
		{
			name: "too much pupil ids",
			args: args{
				pupilIDs: make([]garbage.PupilID, schooling.MaxRemovePupils+1),
			},
			wantIDsLen: 0,
			wantErr:    true,
		},
		{
			name: "one id is empty",
			args: args{
				pupilIDs: []garbage.PupilID{"123", ""},
			},
			wantIDsLen: 0,
			wantErr:    true,
		},
		{
			name: "repo's error",
			args: args{
				pupilIDs: []garbage.PupilID{"error"},
			},
			wantIDsLen: 0,
			wantErr:    true,
		},
		{
			name: "remove 10 pupils",
			args: args{
				pupilIDs: ids,
			},
			wantIDsLen: len(ids),
			wantErr:    false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := s.RemovePupils(ctx, tt.args.pupilIDs)
			gotLen := len(got)
			if (err != nil) != tt.wantErr {
				t.Errorf("RemovePupils() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotLen != tt.wantIDsLen {
				t.Errorf("RemovePupils() gotLen = %v, wantIDsLen %v", gotLen, tt.wantIDsLen)
			}
		})
	}
}

func Test_service_AddPupils(t *testing.T) {
	ctx := context.Background()

	var repo mock.SchoolingRepository
	repo.StorePupilsFn = func(ctx context.Context, pupils []*schooling.Pupil) ([]garbage.PupilID, error) {
		if len(pupils) > 0 && pupils[0].FirstName == "error" {
			return nil, errors.New("repo's error")
		}
		return make([]garbage.PupilID, len(pupils)), nil
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
			name:    "empty last name",
			pupils:  []schooling.PupilBio{{"fn", "", "8B"}},
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
	requestedClass := garbage.Class{
		YearFormed: 2010,
		Letter:     "B",
	}
	someClass := garbage.Class{
		Letter:     "Y",
		YearFormed: 2010,
	}
	pupilWithRequestedClass := &schooling.Pupil{
		Pupil: garbage.Pupil{
			ID:        pupilIDPupilWithRequestedClass,
			FirstName: "FN",
			LastName:  "LN",
		},
		Class: requestedClass,
	}
	ctx := context.Background()

	var repo mock.SchoolingRepository
	repo.StorePupilFn = func(ctx context.Context, pupil *schooling.Pupil) (garbage.PupilID, error) {
		if pupil.ID == pupilIDStoreError {
			return "", errors.New("repo's error")
		}
		return pupil.ID, nil
	}
	repo.PupilByIDFn = func(ctx context.Context, pupilID garbage.PupilID) (pupil *schooling.Pupil, err error) {
		if pupilID == pupilIDPupilNotFound {
			return nil, errors.New("pupil not found error")
		}
		if pupilID == pupilIDPupilWithRequestedClass {
			return pupilWithRequestedClass, nil
		}
		return &schooling.Pupil{
			Pupil: garbage.Pupil{
				ID:        pupilID,
				FirstName: "FN",
				LastName:  "LN",
			},
			Class: someClass,
		}, nil
	}
	s := schooling.NewService(&repo)

	type args struct {
		pupilID   garbage.PupilID
		className string
	}
	tests := []struct {
		name    string
		args    args
		want    garbage.PupilID
		wantErr bool
	}{
		{
			name: "no pupilID",
			args: args{
				pupilID:   "",
				className: "className",
			},
			want:    "",
			wantErr: true,
		},
		{
			name: "no className",
			args: args{
				pupilID:   "pupilID",
				className: "",
			},
			want:    "",
			wantErr: true,
		},
		{
			name: "no pupil with the given id",
			args: args{
				pupilID:   pupilIDPupilNotFound,
				className: "class name",
			},
			want:    "",
			wantErr: true,
		},
		{
			name: "error on storing a pupil",
			args: args{
				pupilID:   pupilIDStoreError,
				className: "10B",
			},
			want:    "",
			wantErr: true,
		},
		{
			name: "pupil is already in the given class",
			args: args{
				pupilID:   pupilIDPupilWithRequestedClass,
				className: requestedClassName,
			},
			want:    pupilIDPupilWithRequestedClass,
			wantErr: false,
		},
		{
			name: "swap a pupil's class for a found one",
			args: args{
				pupilID:   "pupilID",
				className: "10Y",
			},
			want:    "pupilID",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := s.ChangePupilClass(ctx, tt.args.pupilID, tt.args.className)
			if (err != nil) != tt.wantErr {
				t.Errorf("ChangePupilClass() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("ChangePupilClass() got = %v, want %v", got, tt.want)
			}
		})
	}
}
