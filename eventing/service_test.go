package eventing_test

import (
	"context"
	"errors"
	"reflect"
	"testing"
	"time"

	"github.com/shanvl/garbage-events-service"
	"github.com/shanvl/garbage-events-service/eventing"
	"github.com/shanvl/garbage-events-service/mock"
)

func Test_service_CreateEvent(t *testing.T) {
	var repository mock.EventingRepository
	repository.StoreEventFn = func(ctx context.Context, e *garbage.Event) (id garbage.EventID, err error) {
		return e.ID, nil
	}
	s := eventing.NewService(&repository)

	ctx := context.Background()

	type args struct {
		ctx              context.Context
		date             time.Time
		name             string
		resourcesAllowed []garbage.Resource
	}
	tests := []struct {
		name string
		args args
		// we check if the length of eventID returned is greater than 0. On error it will be 0
		idLenGT0 bool
		wantErr  bool
	}{
		{
			name: "date is in the past",
			args: args{
				ctx:              ctx,
				date:             time.Now().AddDate(0, 0, -1),
				name:             "some name",
				resourcesAllowed: []garbage.Resource{"plastic", "gadgets"},
			},
			idLenGT0: false,
			wantErr:  true,
		},
		{
			name: "wrong resources",
			args: args{
				ctx:              ctx,
				date:             time.Now().AddDate(0, 0, 1),
				name:             "some name",
				resourcesAllowed: []garbage.Resource{"plastI", "gadgets"},
			},
			idLenGT0: false,
			wantErr:  true,
		},
		{
			name: "no resources",
			args: args{
				ctx:              ctx,
				date:             time.Now().AddDate(0, 0, 1),
				name:             "",
				resourcesAllowed: nil,
			},
			idLenGT0: false,
			wantErr:  true,
		},
		{
			name: "no name but that's ok",
			args: args{
				ctx:              ctx,
				date:             time.Now().AddDate(0, 0, 1),
				name:             "",
				resourcesAllowed: []garbage.Resource{"plastic", "gadgets"},
			},
			idLenGT0: true,
			wantErr:  false,
		},
		{
			name: "ok case",
			args: args{
				ctx:              ctx,
				date:             time.Now().AddDate(0, 0, 1),
				name:             "some name",
				resourcesAllowed: []garbage.Resource{"plastic", "gadgets"},
			},
			idLenGT0: true,
			wantErr:  false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := s.CreateEvent(tt.args.ctx, tt.args.date, tt.args.name, tt.args.resourcesAllowed)
			gotLenGT0 := len(got) > 0
			if (err != nil) != tt.wantErr {
				t.Errorf("CreateEvent() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotLenGT0 != tt.idLenGT0 {
				t.Errorf("CreateEvent() got = %v, want %v", got, tt.idLenGT0)
			}
		})
	}
}

func Test_service_DeleteEvent(t *testing.T) {
	var repository mock.EventingRepository
	repository.DeleteEventFn = func(ctx context.Context, eventID garbage.EventID) (id garbage.EventID, err error) {
		if eventID == "not_found" {
			return "", errors.New("repo's not found error")
		}
		return eventID, nil
	}
	s := eventing.NewService(&repository)

	ctx := context.Background()

	type args struct {
		ctx     context.Context
		eventID garbage.EventID
	}
	tests := []struct {
		name    string
		args    args
		want    garbage.EventID
		wantErr bool
	}{
		{
			name: "no eventID",
			args: args{
				ctx:     ctx,
				eventID: "",
			},
			want:    "",
			wantErr: true,
		},
		{
			name: "no event with such eventID",
			args: args{
				ctx:     ctx,
				eventID: "not_found",
			},
			want:    "",
			wantErr: true,
		},
		{
			name: "correct eventID",
			args: args{
				ctx:     ctx,
				eventID: "123",
			},
			want:    "123",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := s.DeleteEvent(tt.args.ctx, tt.args.eventID)
			if (err != nil) != tt.wantErr {
				t.Errorf("DeleteEvent() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("DeleteEvent() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_service_Events(t *testing.T) {
	const (
		totalEvents = 55
		sortBy      = eventing.DateDesc
		name        = "some name"
		amount      = 10
		skip        = 50
	)
	date := time.Now().AddDate(0, -1, 0)
	ctx := context.Background()

	var repository mock.EventingRepository
	repository.EventsFn = func(ctx context.Context, name string, date time.Time, sortBy eventing.SortBy, amount int,
		skip int) (events []*eventing.Event, total int, err error) {
		if name == "not_found" {
			return nil, 0, nil
		}
		if name == "error" {
			return nil, 0, errors.New("some error")
		}
		if amount < 0 {
			return make([]*eventing.Event, 0), totalEvents, nil
		}
		events = make([]*eventing.Event, amount)
		return events, totalEvents, nil
	}
	s := eventing.NewService(&repository)

	type args struct {
		ctx    context.Context
		name   string
		date   time.Time
		sortBy eventing.SortBy
		amount int
		skip   int
	}
	tests := []struct {
		name          string
		args          args
		wantEventsLen int
		wantTotal     int
		wantErr       bool
	}{
		{
			name: "empty name",
			args: args{
				ctx:    ctx,
				name:   "",
				date:   date,
				sortBy: sortBy,
				amount: amount,
				skip:   skip,
			},
			wantEventsLen: amount,
			wantTotal:     totalEvents,
			wantErr:       false,
		},
		{
			name: "zero date",
			args: args{
				ctx:    ctx,
				name:   name,
				date:   time.Time{},
				sortBy: sortBy,
				amount: amount,
				skip:   skip,
			},
			wantEventsLen: amount,
			wantTotal:     totalEvents,
			wantErr:       false,
		},
		{
			name: "negative amount",
			args: args{
				ctx:    ctx,
				name:   name,
				date:   date,
				sortBy: sortBy,
				amount: -55,
				skip:   skip,
			},
			wantEventsLen: 0,
			wantTotal:     totalEvents,
			wantErr:       false,
		},
		{
			name: "negative skip",
			args: args{
				ctx:    ctx,
				name:   name,
				date:   time.Time{},
				sortBy: sortBy,
				amount: amount,
				skip:   -55,
			},
			wantEventsLen: amount,
			wantTotal:     totalEvents,
			wantErr:       false,
		},
		{
			name: "invalid sortBy",
			args: args{
				ctx:    ctx,
				name:   name,
				date:   date,
				sortBy: "invalid",
				amount: amount,
				skip:   skip,
			},
			wantEventsLen: amount,
			wantTotal:     totalEvents,
			wantErr:       false,
		},
		{
			name: "repo's internal error",
			args: args{
				ctx:    ctx,
				name:   "error",
				date:   date,
				sortBy: sortBy,
				amount: amount,
				skip:   skip,
			},
			wantEventsLen: 0,
			wantTotal:     0,
			wantErr:       true,
		},
		{
			name: "not found",
			args: args{
				ctx:    ctx,
				name:   "not_found",
				date:   date,
				sortBy: sortBy,
				amount: amount,
				skip:   skip,
			},
			wantEventsLen: 0,
			wantTotal:     0,
			wantErr:       false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotEvents, gotTotal, err := s.Events(tt.args.ctx, tt.args.name, tt.args.date, tt.args.sortBy,
				tt.args.amount, tt.args.skip)
			if (err != nil) != tt.wantErr {
				t.Errorf("Events() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if len(gotEvents) != tt.wantEventsLen {
				t.Errorf("Events() gotEventsLen = %v, want %v", gotEvents, tt.wantEventsLen)
			}
			if gotTotal != tt.wantTotal {
				t.Errorf("Events() gotTotal = %v, want %v", gotTotal, tt.wantTotal)
			}
		})
	}
}

func Test_service_Event(t *testing.T) {
	var repository mock.EventingRepository
	repository.EventFn = func(ctx context.Context, id garbage.EventID) (event *eventing.Event, err error) {
		if id == "not_found" {
			return nil, errors.New("not found")
		}
		if id == "error" {
			return nil, errors.New("some error")
		}
		return &eventing.Event{Event: garbage.Event{ID: id}}, nil
	}
	ctx := context.Background()

	type args struct {
		ctx     context.Context
		eventID garbage.EventID
	}
	tests := []struct {
		name    string
		args    args
		want    *eventing.Event
		wantErr bool
	}{
		{
			name: "empty eventID",
			args: args{
				ctx:     ctx,
				eventID: "",
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "not found",
			args: args{
				ctx:     ctx,
				eventID: "not_found",
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "error",
			args: args{
				ctx:     ctx,
				eventID: "error",
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "event is found",
			args: args{
				ctx:     ctx,
				eventID: "123",
			},
			want:    &eventing.Event{Event: garbage.Event{ID: "123"}},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := eventing.NewService(&repository)
			got, err := s.Event(tt.args.ctx, tt.args.eventID)
			if (err != nil) != tt.wantErr {
				t.Errorf("Event() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Event() got = %v, want %v", got, tt.want)
			}
		})
	}
}
