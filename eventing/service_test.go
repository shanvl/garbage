package eventing

import (
	"testing"
	"time"

	"github.com/shanvl/garbage-events-service"
	"github.com/shanvl/garbage-events-service/mock"
	"github.com/shanvl/garbage-events-service/validation"
)

func Test_service_CreateEvent(t *testing.T) {
	repository := &mock.EventingRepository{
		StoreEventFn: func(e *garbage.Event) (id garbage.EventID, err error) {
			return e.ID, nil
		},
	}
	idGenerator := &mock.IDGenerator{
		GenerateEventIDFn: func() garbage.EventID {
			return "123"
		},
	}
	validator := validation.NewValidator()
	s := NewService(repository, idGenerator, validator)
	type args struct {
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
			got, err := s.CreateEvent(tt.args.date, tt.args.name, tt.args.resourcesAllowed)
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
