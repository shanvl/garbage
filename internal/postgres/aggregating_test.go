package postgres_test

import (
	"context"
	"testing"
	"time"

	"github.com/shanvl/garbage-events-service/internal/garbage"
	"github.com/shanvl/garbage-events-service/internal/postgres"
	"github.com/shanvl/garbage-events-service/internal/sorting"
	"github.com/shanvl/garbage-events-service/internal/usecases/aggregating"
)

func TestAggregatingRepo_PupilByID(t *testing.T) {
	r := postgres.NewAggregatingRepo(db)
	ctx := context.Background()
	pID := getPupilID(t)

	type args struct {
		id            garbage.PupilID
		filters       aggregating.EventsByDateFilter
		eventsSorting sorting.By
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "no dates",
			args: args{
				id:            pID,
				filters:       aggregating.EventsByDateFilter{},
				eventsSorting: sorting.Plastic,
			},
			wantErr: false,
		},
		{
			name: "from is set",
			args: args{
				id: pID,
				filters: aggregating.EventsByDateFilter{
					From: newDate(2018, 12, 12),
					To:   time.Time{},
				},
				eventsSorting: sorting.Plastic,
			},
			wantErr: false,
		},
		{
			name: "to is set",
			args: args{
				id: pID,
				filters: aggregating.EventsByDateFilter{
					From: time.Time{},
					To:   newDate(2018, 12, 12),
				},
				eventsSorting: sorting.Plastic,
			},
			wantErr: false,
		},
		{
			name: "dates are set",
			args: args{
				id: pID,
				filters: aggregating.EventsByDateFilter{
					From: newDate(2018, 12, 12),
					To:   newDate(2020, 12, 12),
				},
				eventsSorting: sorting.Plastic,
			},
			wantErr: false,
		},
		{
			name: "dates are in the future",
			args: args{
				id: pID,
				filters: aggregating.EventsByDateFilter{
					From: time.Now().AddDate(50, 50, 50),
					To:   time.Now().AddDate(51, 51, 51),
				},
				eventsSorting: sorting.Plastic,
			},
			wantErr: false,
		},
		{
			name: "invalid pupil id",
			args: args{
				id:            "wrongid",
				filters:       aggregating.EventsByDateFilter{},
				eventsSorting: sorting.Plastic,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := r.PupilByID(ctx, tt.args.id, tt.args.filters, tt.args.eventsSorting)
			if (err != nil) != tt.wantErr {
				t.Errorf("PupilByID() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func newDate(year int, month int, day int) time.Time {
	return time.Date(year, time.Month(month), day, 0, 0, 0, 0, time.UTC)
}
