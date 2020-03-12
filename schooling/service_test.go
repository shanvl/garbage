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
