package garbage

import (
	"reflect"
	"testing"
	"time"
)

func TestNewEvent(t *testing.T) {
	eventDate := time.Date(2020, 10, 10, 0, 0, 0, 0, time.UTC)
	resourcesAllowed := []Resource{"plastic", "gadgets"}
	resourcesBrought := make(map[Resource]int)
	type args struct {
		id               EventID
		date             time.Time
		name             string
		resourcesAllowed []Resource
		resourcesBrought map[Resource]int
	}
	tests := []struct {
		name string
		args args
		want *Event
	}{
		{
			name: "with no name provided",
			args: args{
				id:               "123",
				date:             eventDate,
				name:             "",
				resourcesAllowed: resourcesAllowed,
				resourcesBrought: resourcesBrought,
			},
			want: &Event{
				ID:               "123",
				Date:             eventDate,
				Name:             "10-10-2020",
				ResourcesAllowed: resourcesAllowed,
				ResourcesBrought: resourcesBrought,
			},
		},
		{
			name: "with a name",
			args: args{
				id:               "123",
				date:             eventDate,
				name:             "some name",
				resourcesAllowed: resourcesAllowed,
				resourcesBrought: resourcesBrought,
			},
			want: &Event{
				ID:               "123",
				Date:             eventDate,
				Name:             "some name",
				ResourcesAllowed: resourcesAllowed,
				ResourcesBrought: resourcesBrought,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewEvent(tt.args.id, tt.args.date, tt.args.name, tt.args.resourcesAllowed); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewEvent() = %v, want %v", got, tt.want)
			}
		})
	}
}
