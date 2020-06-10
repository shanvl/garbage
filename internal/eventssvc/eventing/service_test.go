package eventing_test

import (
	"context"
	"errors"
	"reflect"
	"testing"
	"time"

	"github.com/shanvl/garbage/internal/eventssvc"
	"github.com/shanvl/garbage/internal/eventssvc/eventing"
	"github.com/shanvl/garbage/internal/eventssvc/mock"
	"github.com/shanvl/garbage/internal/eventssvc/sorting"
)

func Test_service_CreateEvent(t *testing.T) {
	var repository mock.EventingRepository
	repository.StoreEventFn = func(ctx context.Context, e *eventssvc.Event) (id eventssvc.EventID, err error) {
		return e.ID, nil
	}
	s := eventing.NewService(&repository)

	ctx := context.Background()

	type args struct {
		ctx              context.Context
		date             time.Time
		name             string
		resourcesAllowed []eventssvc.Resource
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
				resourcesAllowed: []eventssvc.Resource{"plastic", "gadgets"},
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
				resourcesAllowed: []eventssvc.Resource{"plastI", "gadgets"},
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
				resourcesAllowed: []eventssvc.Resource{"plastic", "gadgets"},
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
				resourcesAllowed: []eventssvc.Resource{"plastic", "gadgets"},
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
	const repoErrorEventID = "error"
	var repository mock.EventingRepository
	repository.DeleteEventFn = func(ctx context.Context, eventID eventssvc.EventID) error {
		if eventID == repoErrorEventID {
			return eventssvc.ErrUnknownEvent
		}
		return nil
	}
	s := eventing.NewService(&repository)

	ctx := context.Background()

	type args struct {
		ctx     context.Context
		eventID eventssvc.EventID
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "no eventID",
			args: args{
				ctx:     ctx,
				eventID: "",
			},
			wantErr: true,
		},
		{
			name: "no event with such eventID",
			args: args{
				ctx:     ctx,
				eventID: repoErrorEventID,
			},
			wantErr: true,
		},
		{
			name: "correct eventID",
			args: args{
				ctx:     ctx,
				eventID: "123",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := s.DeleteEvent(tt.args.ctx, tt.args.eventID)
			if (err != nil) != tt.wantErr {
				t.Errorf("DeleteEvent() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func Test_service_EventByID(t *testing.T) {
	var repository mock.EventingRepository
	repository.EventByIDFn = func(ctx context.Context, id eventssvc.EventID) (event *eventing.Event, err error) {
		if id == "not_found" {
			return nil, errors.New("not found")
		}
		if id == "error" {
			return nil, errors.New("some error")
		}
		return &eventing.Event{Event: eventssvc.Event{ID: id}}, nil
	}
	s := eventing.NewService(&repository)

	ctx := context.Background()

	type args struct {
		ctx     context.Context
		eventID eventssvc.EventID
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
			want:    &eventing.Event{Event: eventssvc.Event{ID: "123"}},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := s.EventByID(tt.args.ctx, tt.args.eventID)
			if (err != nil) != tt.wantErr {
				t.Errorf("EventByID() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("EventByID() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_service_ChangePupilResources(t *testing.T) {
	const (
		eventID                = "123"
		eventIDErrNoEventPupil = "errornoeventpupil"
		pupilID                = "123"
	)
	resourcesBrought := eventssvc.ResourceMap{eventssvc.Plastic: 22}
	resourcesAllowed := []eventssvc.Resource{eventssvc.Plastic, eventssvc.Gadgets}
	ctx := context.Background()

	var repository mock.EventingRepository
	repository.ChangePupilResourcesFn = func(ctx context.Context, eventID eventssvc.EventID, pupilID eventssvc.PupilID,
		resources eventssvc.ResourceMap) error {
		if eventID == eventIDErrNoEventPupil {
			return eventing.ErrNoEventPupil
		}
		return nil
	}
	repository.EventByIDFn = func(ctx context.Context, id eventssvc.EventID) (event *eventing.Event, err error) {
		return &eventing.Event{
				Event: eventssvc.Event{ID: id, ResourcesAllowed: resourcesAllowed},
			},
			nil
	}
	s := eventing.NewService(&repository)

	type args struct {
		ctx       context.Context
		eventID   eventssvc.EventID
		pupilID   eventssvc.PupilID
		resources eventssvc.ResourceMap
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "no eventID",
			args: args{
				ctx:       ctx,
				eventID:   "",
				pupilID:   pupilID,
				resources: resourcesBrought,
			},
			wantErr: true,
		},
		{
			name: "no pupilID",
			args: args{
				ctx:       ctx,
				eventID:   eventID,
				pupilID:   "",
				resources: resourcesBrought,
			},
			wantErr: true,
		},
		{
			name: "no resources",
			args: args{
				ctx:       ctx,
				eventID:   eventID,
				pupilID:   pupilID,
				resources: nil,
			},
			wantErr: true,
		},
		{
			name: "resource is not allowed",
			args: args{
				ctx:       ctx,
				eventID:   eventID,
				pupilID:   pupilID,
				resources: eventssvc.ResourceMap{eventssvc.Paper: 1},
			},
			wantErr: true,
		},
		{
			name: "one resource is allowed, another not",
			args: args{
				ctx:       ctx,
				eventID:   eventID,
				pupilID:   pupilID,
				resources: eventssvc.ResourceMap{eventssvc.Paper: 11, eventssvc.Plastic: 33},
			},
			wantErr: true,
		},
		{
			name: "ErrNoEventPupil",
			args: args{
				ctx:       ctx,
				eventID:   eventIDErrNoEventPupil,
				pupilID:   pupilID,
				resources: eventssvc.ResourceMap{eventssvc.Plastic: 11, eventssvc.Gadgets: 33},
			},
			wantErr: true,
		},
		{
			name: "add 2 resources",
			args: args{
				ctx:       ctx,
				eventID:   eventID,
				pupilID:   pupilID,
				resources: eventssvc.ResourceMap{eventssvc.Plastic: 11, eventssvc.Gadgets: 33},
			},
			wantErr: false,
		},
		{
			name: "subtract one resource, add another",
			args: args{
				ctx:       ctx,
				eventID:   eventID,
				pupilID:   pupilID,
				resources: eventssvc.ResourceMap{eventssvc.Plastic: -55, eventssvc.Gadgets: 33},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := s.ChangePupilResources(tt.args.ctx, tt.args.eventID, tt.args.pupilID, tt.args.resources)
			if (err != nil) != tt.wantErr {
				t.Errorf("ChangePupilResources() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.args.eventID == eventIDErrNoEventPupil && !errors.Is(err, eventssvc.ErrUnknownPupil) {
				t.Errorf("ChangePupilResources() on repo's ErrNoEventPupil didn't return eventssvc."+
					"ErrUnknownPupil: %v", err)
				return
			}
		})
	}
}

func Test_service_EventPupils(t *testing.T) {
	const (
		totalPupils = 75
		sortBy      = sorting.NameAsc
		eventID     = "123"
		amount      = 10
		skip        = 50
	)
	ctx := context.Background()

	var repository mock.EventingRepository
	repository.EventPupilsFn = func(ctx context.Context, eventID eventssvc.EventID,
		filters eventing.EventPupilFilters, sortBy sorting.By, amount int, skip int) (pupils []*eventing.Pupil,
		total int, err error) {

		if eventID == "not_found" {
			return nil, 0, errors.New("not found error")
		}
		if eventID == "error" {
			return nil, 0, errors.New("some error")
		}
		if amount < 0 {
			return make([]*eventing.Pupil, 0), totalPupils, nil
		}
		pupils = make([]*eventing.Pupil, amount)
		return pupils, totalPupils, nil
	}
	s := eventing.NewService(&repository)

	type args struct {
		ctx     context.Context
		eventID eventssvc.EventID
		filters eventing.EventPupilFilters
		sortBy  sorting.By
		amount  int
		skip    int
	}
	tests := []struct {
		name          string
		args          args
		wantPupilsLen int
		wantTotal     int
		wantErr       bool
	}{
		{
			name: "no eventID",
			args: args{
				ctx:     ctx,
				eventID: "",
				sortBy:  sortBy,
				amount:  amount,
				skip:    skip,
			},
			wantPupilsLen: 0,
			wantTotal:     0,
			wantErr:       true,
		},
		{
			name: "negative amount",
			args: args{
				ctx:     ctx,
				eventID: eventID,
				sortBy:  sortBy,
				amount:  -55,
				skip:    skip,
			},
			wantPupilsLen: eventing.DefaultAmount,
			wantTotal:     totalPupils,
			wantErr:       false,
		},
		{
			name: "negative skip",
			args: args{
				ctx:     ctx,
				eventID: eventID,
				sortBy:  sortBy,
				amount:  amount,
				skip:    -120,
			},
			wantPupilsLen: amount,
			wantTotal:     totalPupils,
			wantErr:       false,
		},
		{
			name: "invalid sortBy",
			args: args{
				ctx:     ctx,
				eventID: eventID,
				sortBy:  "invalid",
				amount:  amount,
				skip:    skip,
			},
			wantPupilsLen: amount,
			wantTotal:     totalPupils,
			wantErr:       false,
		},
		{
			name: "repo's internal error",
			args: args{
				ctx:     ctx,
				eventID: "error",
				sortBy:  sortBy,
				amount:  amount,
				skip:    skip,
			},
			wantPupilsLen: 0,
			wantTotal:     0,
			wantErr:       true,
		},
		{
			name: "no pupils found",
			args: args{
				ctx:     ctx,
				eventID: "not_found",
				sortBy:  sortBy,
				amount:  amount,
				skip:    skip,
			},
			wantPupilsLen: 0,
			wantTotal:     0,
			wantErr:       true,
		},
		{
			name: "get 10, skip 50",
			args: args{
				ctx:     ctx,
				eventID: eventID,
				sortBy:  sortBy,
				amount:  10,
				skip:    50,
			},
			wantPupilsLen: 10,
			wantTotal:     totalPupils,
			wantErr:       false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotPupils, gotTotal, err := s.EventPupils(tt.args.ctx, tt.args.eventID, tt.args.filters, tt.args.sortBy,
				tt.args.amount, tt.args.skip)
			gotPupilsLen := len(gotPupils)

			if (err != nil) != tt.wantErr {
				t.Errorf("EventPupils() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotPupilsLen != tt.wantPupilsLen {
				t.Errorf("EventPupils() gotPupilsLen = %v, want %v", gotPupilsLen, tt.wantPupilsLen)
			}
			if gotTotal != tt.wantTotal {
				t.Errorf("EventPupils() gotTotal = %v, want %v", gotTotal, tt.wantTotal)
			}
		})
	}
}

func Test_service_EventClasses(t *testing.T) {
	const (
		totalClasses = 75
		sortBy       = sorting.NameAsc
		eventID      = "123"
		amount       = 10
		skip         = 50
	)
	ctx := context.Background()

	var repository mock.EventingRepository
	repository.EventClassesFn = func(ctx context.Context, eventID eventssvc.EventID,
		filters eventing.EventClassFilters, sortBy sorting.By, amount int, skip int) (classes []*eventing.Class,
		total int, err error) {

		if eventID == "not_found" {
			return nil, 0, errors.New("not found error")
		}
		if eventID == "error" {
			return nil, 0, errors.New("some error")
		}
		if amount < 0 {
			return make([]*eventing.Class, 0), totalClasses, nil
		}
		classes = make([]*eventing.Class, amount)
		return classes, totalClasses, nil
	}
	s := eventing.NewService(&repository)

	type args struct {
		ctx     context.Context
		eventID eventssvc.EventID
		filters eventing.EventClassFilters
		sortBy  sorting.By
		amount  int
		skip    int
	}
	tests := []struct {
		name           string
		args           args
		wantClassesLen int
		wantTotal      int
		wantErr        bool
	}{
		{
			name: "no eventID",
			args: args{
				ctx:     ctx,
				eventID: "",
				sortBy:  sortBy,
				amount:  amount,
				skip:    skip,
			},
			wantClassesLen: 0,
			wantTotal:      0,
			wantErr:        true,
		},
		{
			name: "negative amount",
			args: args{
				ctx:     ctx,
				eventID: eventID,
				sortBy:  sortBy,
				amount:  -55,
				skip:    skip,
			},
			wantClassesLen: eventing.DefaultAmount,
			wantTotal:      totalClasses,
			wantErr:        false,
		},
		{
			name: "negative skip",
			args: args{
				ctx:     ctx,
				eventID: eventID,
				sortBy:  sortBy,
				amount:  amount,
				skip:    -120,
			},
			wantClassesLen: amount,
			wantTotal:      totalClasses,
			wantErr:        false,
		},
		{
			name: "invalid sortBy",
			args: args{
				ctx:     ctx,
				eventID: eventID,
				sortBy:  "invalid",
				amount:  amount,
				skip:    skip,
			},
			wantClassesLen: amount,
			wantTotal:      totalClasses,
			wantErr:        false,
		},
		{
			name: "repo's internal error",
			args: args{
				ctx:     ctx,
				eventID: "error",
				sortBy:  sortBy,
				amount:  amount,
				skip:    skip,
			},
			wantClassesLen: 0,
			wantTotal:      0,
			wantErr:        true,
		},
		{
			name: "no pupils found",
			args: args{
				ctx:     ctx,
				eventID: "not_found",
				sortBy:  sortBy,
				amount:  amount,
				skip:    skip,
			},
			wantClassesLen: 0,
			wantTotal:      0,
			wantErr:        true,
		},
		{
			name: "get 10, skip 50",
			args: args{
				ctx:     ctx,
				eventID: eventID,
				sortBy:  sortBy,
				amount:  10,
				skip:    50,
			},
			wantClassesLen: 10,
			wantTotal:      totalClasses,
			wantErr:        false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotClasses, gotTotal, err := s.EventClasses(tt.args.ctx, tt.args.eventID, tt.args.filters,
				tt.args.sortBy, tt.args.amount, tt.args.skip)
			gotClassesLen := len(gotClasses)

			if (err != nil) != tt.wantErr {
				t.Errorf("EventClasses() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotClassesLen != tt.wantClassesLen {
				t.Errorf("EventClasses() gotClassesLen = %v, want %v", gotClassesLen, tt.wantClassesLen)
			}
			if gotTotal != tt.wantTotal {
				t.Errorf("EventClasses() gotTotal = %v, want %v", gotTotal, tt.wantTotal)
			}
		})
	}
}

func Test_service_PupilByID(t *testing.T) {
	const pupilIDPupilNotFound = "not found"
	foundPupil := &eventing.Pupil{
		Pupil: eventssvc.Pupil{
			ID:        "123",
			FirstName: "FN",
			LastName:  "LN",
		},
		Class: "3B",
	}
	ctx := context.Background()

	var repository mock.EventingRepository
	repository.PupilByIDFn = func(ctx context.Context, pupilID eventssvc.PupilID,
		eventID eventssvc.EventID) (event *eventing.Pupil, err error) {

		if pupilID == pupilIDPupilNotFound {
			return nil, errors.New("not found")
		}
		return foundPupil, nil
	}
	s := eventing.NewService(&repository)

	type args struct {
		pupilID eventssvc.PupilID
		eventID eventssvc.EventID
	}
	tests := []struct {
		name    string
		args    args
		want    *eventing.Pupil
		wantErr bool
	}{
		{
			name: "no pupilID",
			args: args{
				pupilID: "",
				eventID: "eventID",
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "no eventID",
			args: args{
				pupilID: "pupilID",
				eventID: "",
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "no pupil was found",
			args: args{
				pupilID: pupilIDPupilNotFound,
				eventID: "eventID",
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "pupil was found",
			args: args{
				pupilID: foundPupil.ID,
				eventID: "eventID",
			},
			want:    foundPupil,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := s.PupilByID(ctx, tt.args.pupilID, tt.args.eventID)
			if (err != nil) != tt.wantErr {
				t.Errorf("PupilByID() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("PupilByID() got = %v, want %v", got, tt.want)
			}
		})
	}
}
