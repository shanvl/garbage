package schooling_test

import (
	"context"
	"errors"
	"testing"

	"github.com/shanvl/garbage-events-service/garbage"
	"github.com/shanvl/garbage-events-service/mock"
	"github.com/shanvl/garbage-events-service/schooling"
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
	repo.AddPupilsFn = func(ctx context.Context, pupilsInfo []schooling.Pupil) ([]garbage.PupilID, error) {
		if len(pupilsInfo) > 0 && pupilsInfo[0].FirstName == "error" {
			return nil, errors.New("repo's error")
		}
		return make([]garbage.PupilID, len(pupilsInfo)), nil
	}
	s := schooling.NewService(&repo)

	tests := []struct {
		name    string
		pupils  []schooling.PupilInfo
		wantLen int
		wantErr bool
	}{
		{
			name:    "repo's error",
			pupils:  []schooling.PupilInfo{{"error", "ln", "8B"}},
			wantLen: 0,
			wantErr: true,
		},
		{
			name:    "empty name",
			pupils:  []schooling.PupilInfo{{"", "ln", "8B"}},
			wantLen: 0,
			wantErr: true,
		},
		{
			name:    "empty last name",
			pupils:  []schooling.PupilInfo{{"fn", "", "8B"}},
			wantLen: 0,
			wantErr: true,
		},
		{
			name:    "invalid className",
			pupils:  []schooling.PupilInfo{{"fn", "ln", "B8B"}},
			wantLen: 0,
			wantErr: true,
		},
		{
			name: "3 valid pupils",
			pupils: []schooling.PupilInfo{
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
			pupils: []schooling.PupilInfo{
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
				t.Errorf("AddPupils() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotLen != tt.wantLen {
				t.Errorf("AddPupils() gotLen = %v, wantLen %v", got, tt.wantLen)
			}
		})
	}
}
