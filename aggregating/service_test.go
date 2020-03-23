package aggregating_test

import (
	"context"
	"errors"
	"reflect"
	"testing"

	"github.com/shanvl/garbage-events-service/aggregating"
	"github.com/shanvl/garbage-events-service/garbage"
	"github.com/shanvl/garbage-events-service/mock"
	"github.com/shanvl/garbage-events-service/sorting"
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

func Test_service_ClassByID(t *testing.T) {
	const (
		repoError = "e"
	)
	class := &aggregating.Class{}
	ctx := context.Background()

	var repo mock.AggregatingRepository
	repo.ClassByIDFn = func(ctx context.Context, id garbage.ClassID, filters aggregating.EventsByDateFilter,
		eventsSorting sorting.By) (*aggregating.Class, error) {

		if id == repoError {
			return nil, errors.New("some error")
		}
		return class, nil
	}
	s := aggregating.NewService(&repo)

	type args struct {
		id            garbage.ClassID
		filters       aggregating.EventsByDateFilter
		eventsSorting sorting.By
	}
	tests := []struct {
		name    string
		args    args
		want    *aggregating.Class
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
			name: "no eventID",
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
			want:    class,
			wantErr: false,
		},
		{
			name: "ok args",
			args: args{
				"id",
				aggregating.EventsByDateFilter{},
				sorting.DateDes,
			},
			want:    class,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotClass, err := s.ClassByID(ctx, tt.args.id, tt.args.filters, tt.args.eventsSorting)
			if (err != nil) != tt.wantErr {
				t.Errorf("ClassByID() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotClass, tt.want) {
				t.Errorf("ClassByID() gotClass = %v, want %v", gotClass, tt.want)
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
