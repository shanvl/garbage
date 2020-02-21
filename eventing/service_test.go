package eventing

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/shanvl/garbage-events-service"
	"github.com/shanvl/garbage-events-service/mock"
	"github.com/shanvl/garbage-events-service/validation"
)

func Test_service_CreateEvent(t *testing.T) {
	var repository mock.EventingRepository
	repository.StoreEventFn = func(ctx context.Context, e *garbage.Event) (id garbage.EventID, err error) {
		return e.ID, nil
	}
	var idGenerator mock.IDGenerator
	idGenerator.GenerateEventIDFn = func() garbage.EventID {
		return "123"
	}
	validator := validation.NewValidator()
	s := NewService(&repository, &idGenerator, validator)

	ctx := context.Background()

	type args struct {
		ctx              context.Context
		date             time.Time
		name             string
		resourcesAllowed []garbage.Resource
	}
	tests := []struct {
		name    string
		args    args
		want    garbage.EventID
		wantErr bool
	}{
		{
			name: "date is in the past",
			args: args{
				ctx:              ctx,
				date:             time.Now().AddDate(0, 0, -1),
				name:             "some name",
				resourcesAllowed: []garbage.Resource{"plastic", "gadgets"},
			},
			want:    "",
			wantErr: true,
		},
		{
			name: "wrong resources",
			args: args{
				ctx:              ctx,
				date:             time.Now().AddDate(0, 0, 1),
				name:             "some name",
				resourcesAllowed: []garbage.Resource{"plastI", "gadgets"},
			},
			want:    "",
			wantErr: true,
		},
		{
			name: "no resources",
			args: args{
				ctx:              ctx,
				date:             time.Now().AddDate(0, 0, 1),
				name:             "",
				resourcesAllowed: nil,
			},
			want:    "",
			wantErr: true,
		},
		{
			name: "no name but that's ok",
			args: args{
				ctx:              ctx,
				date:             time.Now().AddDate(0, 0, 1),
				name:             "",
				resourcesAllowed: []garbage.Resource{"plastic", "gadgets"},
			},
			want:    "123",
			wantErr: false,
		},
		{
			name: "👍",
			args: args{
				ctx:              ctx,
				date:             time.Now().AddDate(0, 0, 1),
				name:             "some name",
				resourcesAllowed: []garbage.Resource{"plastic", "gadgets"},
			},
			want:    "123",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := s.CreateEvent(tt.args.ctx, tt.args.date, tt.args.name, tt.args.resourcesAllowed)
			if (err != nil) != tt.wantErr {
				t.Errorf("CreateEvent() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("CreateEvent() got = %v, want %v", got, tt.want)
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
	var idGenerator mock.IDGenerator
	idGenerator.GenerateEventIDFn = func() garbage.EventID {
		return "123"
	}
	validator := validation.NewValidator()
	s := NewService(&repository, &idGenerator, validator)

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
			name: "👍",
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
