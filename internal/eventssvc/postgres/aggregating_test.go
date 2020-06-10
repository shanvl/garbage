package postgres_test

import (
	"context"
	"testing"
	"time"

	"github.com/shanvl/garbage/internal/eventssvc"
	"github.com/shanvl/garbage/internal/eventssvc/aggregating"
	"github.com/shanvl/garbage/internal/eventssvc/postgres"
	"github.com/shanvl/garbage/internal/eventssvc/sorting"
)

func TestAggregatingRepo_PupilByID(t *testing.T) {
	r := postgres.NewAggregatingRepo(db)
	ctx := context.Background()
	pID := getPupilID(t)

	type args struct {
		id            eventssvc.PupilID
		filters       aggregating.EventFilters
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
				filters:       aggregating.EventFilters{},
				eventsSorting: sorting.Plastic,
			},
			wantErr: false,
		},
		{
			name: "from is set",
			args: args{
				id: pID,
				filters: aggregating.EventFilters{
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
				filters: aggregating.EventFilters{
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
				filters: aggregating.EventFilters{
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
				filters: aggregating.EventFilters{
					From: time.Now().AddDate(50, 50, 50),
					To:   time.Now().AddDate(51, 51, 51),
				},
				eventsSorting: sorting.Plastic,
			},
			wantErr: false,
		},
		{
			name: "event name is set",
			args: args{
				id: pID,
				filters: aggregating.EventFilters{
					Name: "p",
				},
				eventsSorting: sorting.Plastic,
			},
			wantErr: false,
		},
		{
			name: "resources allowed are set",
			args: args{
				id: pID,
				filters: aggregating.EventFilters{
					ResourcesAllowed: []eventssvc.Resource{eventssvc.Gadgets, eventssvc.Paper, eventssvc.Plastic},
				},
				eventsSorting: sorting.Plastic,
			},
			wantErr: false,
		},
		{
			name: "all filters are set",
			args: args{
				id: pID,
				filters: aggregating.EventFilters{
					From:             newDate(2010, 1, 1),
					To:               newDate(2017, 1, 1),
					Name:             "p",
					ResourcesAllowed: []eventssvc.Resource{eventssvc.Gadgets, eventssvc.Paper, eventssvc.Plastic},
				},
				eventsSorting: sorting.Gadgets,
			},
			wantErr: false,
		},
		{
			name: "invalid pupil id",
			args: args{
				id:            "wrongid",
				filters:       aggregating.EventFilters{},
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

func TestAggregatingRepo_Events(t *testing.T) {
	r := postgres.NewAggregatingRepo(db)
	ctx := context.Background()

	type args struct {
		filters aggregating.EventFilters
		sortBy  sorting.By
		amount  int
		skip    int
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "no filters",
			args: args{
				filters: aggregating.EventFilters{},
				sortBy:  sorting.Plastic,
				amount:  150,
				skip:    0,
			},
			wantErr: false,
		},
		{
			name: "with a name",
			args: args{
				filters: aggregating.EventFilters{
					Name: "a b",
				},
				sortBy: sorting.Plastic,
				amount: 150,
				skip:   0,
			},
			wantErr: false,
		},
		{
			name: "invalid name",
			args: args{
				filters: aggregating.EventFilters{
					Name: "a&b",
				},
				sortBy: sorting.Plastic,
				amount: 150,
				skip:   0,
			},
			wantErr: false,
		},
		{
			name: "from is set",
			args: args{
				filters: aggregating.EventFilters{
					From: newDate(2018, 1, 1),
				},
				sortBy: sorting.Plastic,
				amount: 150,
				skip:   0,
			},
			wantErr: false,
		},
		{
			name: "to is set",
			args: args{
				filters: aggregating.EventFilters{
					To: newDate(2020, 1, 1),
				},
				sortBy: sorting.Plastic,
				amount: 150,
				skip:   0,
			},
			wantErr: false,
		},
		{
			name: "dates are set",
			args: args{
				filters: aggregating.EventFilters{
					From: newDate(2018, 1, 1),
					To:   newDate(2020, 1, 1),
				},
				sortBy: sorting.Plastic,
				amount: 150,
				skip:   0,
			},
			wantErr: false,
		},
		{
			name: "allowed resources are set",
			args: args{
				filters: aggregating.EventFilters{
					ResourcesAllowed: []eventssvc.Resource{eventssvc.Gadgets, eventssvc.Plastic, eventssvc.Paper},
				},
				sortBy: sorting.Plastic,
				amount: 150,
				skip:   0,
			},
			wantErr: false,
		},
		{
			name: "offset > entries",
			args: args{
				sortBy: sorting.Plastic,
				amount: 150,
				skip:   999,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, _, err := r.Events(ctx, tt.args.filters, tt.args.sortBy, tt.args.amount, tt.args.skip)
			if (err != nil) != tt.wantErr {
				t.Errorf("Events() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestAggregatingRepo_Pupils(t *testing.T) {
	r := postgres.NewAggregatingRepo(db)
	ctx := context.Background()

	type args struct {
		filters       aggregating.PupilFilters
		pupilsSorting sorting.By
		eventsSorting sorting.By
		amount        int
		skip          int
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "all filters",
			args: args{
				filters: aggregating.PupilFilters{
					EventFilters: aggregating.EventFilters{
						From:             newDate(2010, 1, 1),
						To:               newDate(2025, 5, 5),
						Name:             "p",
						ResourcesAllowed: []eventssvc.Resource{eventssvc.Gadgets},
					},
					NameAndClass: "a",
				},
				pupilsSorting: sorting.NameDes,
				eventsSorting: sorting.Gadgets,
				amount:        15,
				skip:          0,
			},
			wantErr: false,
		},
		{
			name: "no filters",
			args: args{
				filters:       aggregating.PupilFilters{},
				pupilsSorting: sorting.NameDes,
				eventsSorting: sorting.Gadgets,
				amount:        15,
				skip:          0,
			},
			wantErr: false,
		},
		{
			name: "no pupils",
			args: args{
				filters: aggregating.PupilFilters{
					NameAndClass: "zzzzzzznopupils",
				},
				pupilsSorting: sorting.NameDes,
				eventsSorting: sorting.Gadgets,
				amount:        15,
				skip:          0,
			},
			wantErr: false,
		},
		{
			name: "no events",
			args: args{
				filters: aggregating.PupilFilters{
					EventFilters: aggregating.EventFilters{
						Name: "zzzzzzznoevents",
					},
				},
				pupilsSorting: sorting.NameDes,
				eventsSorting: sorting.Gadgets,
				amount:        15,
				skip:          0,
			},
			wantErr: false,
		},
		{
			name: "offset > rows",
			args: args{
				filters:       aggregating.PupilFilters{},
				pupilsSorting: sorting.NameDes,
				eventsSorting: sorting.Gadgets,
				amount:        15,
				skip:          999,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, _, err := r.Pupils(ctx, tt.args.filters, tt.args.pupilsSorting, tt.args.eventsSorting,
				tt.args.amount, tt.args.skip)

			if (err != nil) != tt.wantErr {
				t.Errorf("Pupils() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestAggregatingRepo_Classes(t *testing.T) {
	r := postgres.NewAggregatingRepo(db)
	ctx := context.Background()

	type args struct {
		filters        aggregating.ClassFilters
		classesSorting sorting.By
		eventsSorting  sorting.By
		amount         int
		skip           int
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "all filters",
			args: args{
				filters: aggregating.ClassFilters{
					EventFilters: aggregating.EventFilters{
						From:             newDate(2010, 1, 1),
						To:               newDate(2025, 5, 5),
						Name:             "p",
						ResourcesAllowed: []eventssvc.Resource{eventssvc.Gadgets},
					},
					DateFormed: newDate(2009, 9, 1),
					Letter:     "a",
				},
				classesSorting: sorting.NameDes,
				eventsSorting:  sorting.Gadgets,
				amount:         15,
				skip:           0,
			},
			wantErr: false,
		},
		{
			name: "no filters",
			args: args{
				filters:        aggregating.ClassFilters{},
				classesSorting: sorting.NameDes,
				eventsSorting:  sorting.Gadgets,
				amount:         15,
				skip:           0,
			},
			wantErr: false,
		},
		{
			name: "no classes",
			args: args{
				filters: aggregating.ClassFilters{
					DateFormed: time.Now().AddDate(50, 1, 1),
				},
				classesSorting: sorting.NameDes,
				eventsSorting:  sorting.Gadgets,
				amount:         15,
				skip:           0,
			},
			wantErr: false,
		},
		{
			name: "no events",
			args: args{
				filters: aggregating.ClassFilters{
					EventFilters: aggregating.EventFilters{
						Name: "zzzzzzznoevents",
					},
				},
				classesSorting: sorting.NameDes,
				eventsSorting:  sorting.Gadgets,
				amount:         15,
				skip:           0,
			},
			wantErr: false,
		},
		{
			name: "offset > rows",
			args: args{
				filters:        aggregating.ClassFilters{},
				classesSorting: sorting.NameDes,
				eventsSorting:  sorting.Gadgets,
				amount:         15,
				skip:           999,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, _, err := r.Classes(ctx, tt.args.filters, tt.args.classesSorting, tt.args.eventsSorting,
				tt.args.amount, tt.args.skip)

			if (err != nil) != tt.wantErr {
				t.Errorf("Classes() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func newDate(year int, month int, day int) time.Time {
	return time.Date(year, time.Month(month), day, 0, 0, 0, 0, time.UTC)
}
