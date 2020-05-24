package aggregating_test

import (
	"context"
	"errors"
	"reflect"
	"testing"

	"github.com/shanvl/garbage-events-service/internal/garbage"
	"github.com/shanvl/garbage-events-service/internal/mock"
	"github.com/shanvl/garbage-events-service/internal/sorting"
	"github.com/shanvl/garbage-events-service/internal/usecases/aggregating"
)

func Test_service_Classes(t *testing.T) {
	const (
		repoError    = "e"
		totalClasses = 25
	)
	classes := []*aggregating.Class{{}, {}, {}}
	ctx := context.Background()

	var repo mock.AggregatingRepository
	repo.ClassesFn = func(ctx context.Context, filters aggregating.ClassesFilters,
		classesSorting, eventsSorting sorting.By, amount, skip int) ([]*aggregating.Class, int, error) {

		if filters.Letter == repoError {
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
				aggregating.ClassesFilters{Letter: repoError},
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
			name: "ok args",
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

func Test_service_Events(t *testing.T) {
	const (
		totalEvents = 55
		sortBy      = sorting.DateDes
		amount      = 3
		skip        = 50
		repoError   = "e"
	)
	events := []*aggregating.Event{{}, {}, {}}
	ctx := context.Background()

	var repository mock.AggregatingRepository
	repository.EventsFn = func(ctx context.Context, filters aggregating.EventsFilters, sortBy sorting.By, amount,
		skip int) ([]*aggregating.Event, int, error) {

		if filters.Name == repoError {
			return nil, 0, errors.New("some error")
		}

		return events, totalEvents, nil
	}
	s := aggregating.NewService(&repository)

	type args struct {
		ctx     context.Context
		filters aggregating.EventsFilters
		sortBy  sorting.By
		amount  int
		skip    int
	}
	tests := []struct {
		name       string
		args       args
		wantEvents []*aggregating.Event
		wantTotal  int
		wantErr    bool
	}{
		{
			name: "negative amount",
			args: args{
				ctx:     ctx,
				filters: aggregating.EventsFilters{},
				sortBy:  sortBy,
				amount:  -55,
				skip:    skip,
			},
			wantEvents: events,
			wantTotal:  totalEvents,
			wantErr:    false,
		},
		{
			name: "negative skip",
			args: args{
				ctx:     ctx,
				filters: aggregating.EventsFilters{},
				sortBy:  sortBy,
				amount:  amount,
				skip:    -55,
			},
			wantEvents: events,
			wantTotal:  totalEvents,
			wantErr:    false,
		},
		{
			name: "invalid sortBy",
			args: args{
				ctx:     ctx,
				filters: aggregating.EventsFilters{},
				sortBy:  "invalid",
				amount:  amount,
				skip:    skip,
			},
			wantEvents: events,
			wantTotal:  totalEvents,
			wantErr:    false,
		},
		{
			name: "ok args",
			args: args{
				ctx:     ctx,
				filters: aggregating.EventsFilters{},
				sortBy:  sortBy,
				amount:  amount,
				skip:    skip,
			},
			wantEvents: events,
			wantTotal:  totalEvents,
			wantErr:    false,
		},
		{
			name: "repo's internal error",
			args: args{
				ctx:     ctx,
				filters: aggregating.EventsFilters{Name: repoError},
				sortBy:  sortBy,
				amount:  amount,
				skip:    skip,
			},
			wantEvents: nil,
			wantTotal:  0,
			wantErr:    true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotEvents, gotTotal, err := s.Events(tt.args.ctx, tt.args.filters, tt.args.sortBy, tt.args.amount,
				tt.args.skip)
			if (err != nil) != tt.wantErr {
				t.Errorf("Events() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotEvents, tt.wantEvents) {
				t.Errorf("Events() gotEvents = %v, want %v", gotEvents, tt.wantEvents)
			}
			if gotTotal != tt.wantTotal {
				t.Errorf("Events() gotTotal = %v, want %v", gotTotal, tt.wantTotal)
			}
		})
	}
}

func Test_service_Pupils(t *testing.T) {
	const (
		repoError   = "e"
		totalPupils = 25
	)
	classes := []*aggregating.Pupil{{}, {}, {}}
	ctx := context.Background()

	var repo mock.AggregatingRepository
	repo.PupilsFn = func(ctx context.Context, filters aggregating.PupilsFilters,
		pupilsSorting, eventsSorting sorting.By, amount, skip int) ([]*aggregating.Pupil, int, error) {

		if filters.Name == repoError {
			return nil, 0, errors.New("some error")
		}
		return classes, totalPupils, nil
	}
	s := aggregating.NewService(&repo)

	type args struct {
		filters       aggregating.PupilsFilters
		pupilsSorting sorting.By
		eventsSorting sorting.By
		amount, skip  int
	}
	tests := []struct {
		name       string
		args       args
		wantPupils []*aggregating.Pupil
		wantTotal  int
		wantErr    bool
	}{
		{
			name: "repo's error",
			args: args{
				aggregating.PupilsFilters{Name: repoError},
				sorting.NameAsc,
				sorting.DateDes,
				25,
				0,
			},
			wantPupils: nil,
			wantTotal:  0,
			wantErr:    true,
		},
		{
			name: "negative amount",
			args: args{
				aggregating.PupilsFilters{},
				sorting.NameAsc,
				sorting.DateDes,
				-10,
				0,
			},
			wantPupils: classes,
			wantTotal:  totalPupils,
			wantErr:    false,
		},
		{
			name: "negative skip",
			args: args{
				aggregating.PupilsFilters{},
				sorting.NameAsc,
				sorting.DateDes,
				25,
				-50,
			},
			wantPupils: classes,
			wantTotal:  totalPupils,
			wantErr:    false,
		},
		{
			name: "invalid pupils sorting",
			args: args{
				aggregating.PupilsFilters{},
				"invalid sorting",
				sorting.DateDes,
				25,
				0,
			},
			wantPupils: classes,
			wantTotal:  totalPupils,
			wantErr:    false,
		},
		{
			name: "invalid events sorting",
			args: args{
				aggregating.PupilsFilters{},
				sorting.DateDes,
				"invalid sorting",
				25,
				0,
			},
			wantPupils: classes,
			wantTotal:  totalPupils,
			wantErr:    false,
		},
		{
			name: "ok args",
			args: args{
				aggregating.PupilsFilters{},
				sorting.NameAsc,
				sorting.DateDes,
				25,
				0,
			},
			wantPupils: classes,
			wantTotal:  totalPupils,
			wantErr:    false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotPupils, gotTotal, err := s.Pupils(ctx, tt.args.filters, tt.args.pupilsSorting,
				tt.args.eventsSorting, tt.args.amount, tt.args.skip)
			if (err != nil) != tt.wantErr {
				t.Errorf("Pupils() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotPupils, tt.wantPupils) {
				t.Errorf("Pupils() gotPupils = %v, want %v", gotPupils, tt.wantPupils)
			}
			if gotTotal != tt.wantTotal {
				t.Errorf("Pupils() gotTotal = %v, want %v", gotTotal, tt.wantTotal)
			}
		})
	}
}

func Test_service_PupilByID(t *testing.T) {
	const (
		repoError = "e"
	)
	pupil := &aggregating.Pupil{}
	ctx := context.Background()

	var repo mock.AggregatingRepository
	repo.PupilByIDFn = func(ctx context.Context, id garbage.PupilID, filters aggregating.EventsByDateFilter,
		eventsSorting sorting.By) (*aggregating.Pupil, error) {

		if id == repoError {
			return nil, errors.New("some error")
		}
		return pupil, nil
	}
	s := aggregating.NewService(&repo)

	type args struct {
		id            garbage.PupilID
		filters       aggregating.EventsByDateFilter
		eventsSorting sorting.By
	}
	tests := []struct {
		name    string
		args    args
		want    *aggregating.Pupil
		wantErr bool
	}{
		{
			name: "repo's error",
			args: args{
				repoError,
				aggregating.EventsByDateFilter{},
				repoError,
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "no pupilID",
			args: args{
				"",
				aggregating.EventsByDateFilter{},
				repoError,
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "invalid events sorting",
			args: args{
				"id",
				aggregating.EventsByDateFilter{},
				"invalid sorting",
			},
			want:    pupil,
			wantErr: false,
		},
		{
			name: "ok args",
			args: args{
				"id",
				aggregating.EventsByDateFilter{},
				sorting.DateDes,
			},
			want:    pupil,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotPupil, err := s.PupilByID(ctx, tt.args.id, tt.args.filters, tt.args.eventsSorting)
			if (err != nil) != tt.wantErr {
				t.Errorf("PupilByID() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotPupil, tt.want) {
				t.Errorf("PupilByID() gotPupil = %v, want %v", gotPupil, tt.want)
			}
		})
	}
}
