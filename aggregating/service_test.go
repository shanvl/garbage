package aggregating_test

import (
	"context"
	"errors"
	"reflect"
	"testing"

	"github.com/shanvl/garbage-events-service/aggregating"
	"github.com/shanvl/garbage-events-service/mock"
	"github.com/shanvl/garbage-events-service/sorting"
)

func Test_service_Classes(t *testing.T) {
	const (
		errorLetter  = "e"
		totalClasses = 25
	)
	classes := []*aggregating.Class{{}, {}, {}}
	ctx := context.Background()

	var repo mock.AggregatingRepository
	repo.ClassesFn = func(ctx context.Context, filters aggregating.ClassesFilters,
		classesSorting, eventsSorting sorting.By, amount, skip int) ([]*aggregating.Class, int, error) {

		if filters.Letter == errorLetter {
			return nil, 0, errors.New("some error")
		}
		return classes, totalClasses, nil
	}
	s := aggregating.NewService(&repo)

	type args struct {
		filters        aggregating.ClassesFilters
		classesSorting sorting.By
		eventsSorting  sorting.By
		amount, skip   int
	}
	tests := []struct {
		name        string
		args        args
		wantClasses []*aggregating.Class
		wantTotal   int
		wantErr     bool
	}{
		{
			name: "repo's error",
			args: args{
				aggregating.ClassesFilters{Letter: errorLetter},
				sorting.NameAsc,
				sorting.DateDes,
				25,
				0,
			},
			wantClasses: nil,
			wantTotal:   0,
			wantErr:     true,
		},
		{
			name: "negative amount",
			args: args{
				aggregating.ClassesFilters{},
				sorting.NameAsc,
				sorting.DateDes,
				-10,
				0,
			},
			wantClasses: classes,
			wantTotal:   totalClasses,
			wantErr:     false,
		},
		{
			name: "negative skip",
			args: args{
				aggregating.ClassesFilters{},
				sorting.NameAsc,
				sorting.DateDes,
				25,
				-50,
			},
			wantClasses: classes,
			wantTotal:   totalClasses,
			wantErr:     false,
		},
		{
			name: "invalid classes sorting",
			args: args{
				aggregating.ClassesFilters{},
				"invalid sorting",
				sorting.DateDes,
				25,
				0,
			},
			wantClasses: classes,
			wantTotal:   totalClasses,
			wantErr:     false,
		},
		{
			name: "invalid events sorting",
			args: args{
				aggregating.ClassesFilters{},
				sorting.DateDes,
				"invalid sorting",
				25,
				0,
			},
			wantClasses: classes,
			wantTotal:   totalClasses,
			wantErr:     false,
		},
		{
			name: "no errors",
			args: args{
				aggregating.ClassesFilters{},
				sorting.NameAsc,
				sorting.DateDes,
				25,
				0,
			},
			wantClasses: classes,
			wantTotal:   totalClasses,
			wantErr:     false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotClasses, gotTotal, err := s.Classes(ctx, tt.args.filters, tt.args.classesSorting,
				tt.args.eventsSorting, tt.args.amount, tt.args.skip)
			if (err != nil) != tt.wantErr {
				t.Errorf("Classes() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotClasses, tt.wantClasses) {
				t.Errorf("Classes() gotClasses = %v, want %v", gotClasses, tt.wantClasses)
			}
			if gotTotal != tt.wantTotal {
				t.Errorf("Classes() gotTotal = %v, want %v", gotTotal, tt.wantTotal)
			}
		})
	}
}
