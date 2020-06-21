package postgres_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/shanvl/garbage/internal/eventssvc"
	"github.com/shanvl/garbage/internal/eventssvc/eventing"
	"github.com/shanvl/garbage/internal/eventssvc/postgres"
	"github.com/shanvl/garbage/internal/eventssvc/sorting"
)

func TestEventingRepo_ChangePupilResources(t *testing.T) {
	r := postgres.NewEventingRepo(db)
	ctx := context.Background()
	pupilID, eventID := getPupilID(t), getEventID(t)
	removePupilResources(t, pupilID, eventID)
	resources := eventssvc.ResourceMap{
		eventssvc.Plastic: 10,
		eventssvc.Paper:   15,
	}
	type args struct {
		eventID   string
		pupilID   string
		resources eventssvc.ResourceMap
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "insert resources",
			args: args{
				eventID:   eventID,
				pupilID:   pupilID,
				resources: resources,
			},
			wantErr: false,
		},
		{
			name: "update resources",
			args: args{
				eventID:   eventID,
				pupilID:   pupilID,
				resources: resources,
			},
			wantErr: false,
		},
		{
			name: "foreign key error",
			args: args{
				eventID:   "some random id",
				pupilID:   "some random id",
				resources: resources,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := r.ChangePupilResources(ctx, tt.args.eventID, tt.args.pupilID,
				tt.args.resources)
			if (err != nil) != tt.wantErr {
				t.Errorf("ChangePupilResources() error = %v, wantErr %v", err, tt.wantErr)
			}
			if tt.name == "foreign key error" && !errors.Is(err, eventing.ErrNoEventPupil) {
				t.Errorf("ChangePupilResources() must return an instance of eventing."+
					"ErrNoEventPupil on foreign key error, returned %v", err)
			}
		})
	}
}
func TestEventingRepo_DeleteEvent(t *testing.T) {
	r := postgres.NewEventingRepo(db)
	ctx := context.Background()
	eventID := getEventID(t)
	type args struct {
		eventID string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "ok",
			args: args{
				eventID: eventID,
			},
			wantErr: false,
		},
		{
			name: "ok2",
			args: args{
				eventID: eventID,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := r.DeleteEvent(ctx, tt.args.eventID); (err != nil) != tt.wantErr {
				t.Errorf("DeleteEvent() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestEventingRepo_EventByID(t *testing.T) {
	r := postgres.NewEventingRepo(db)
	ctx := context.Background()
	type args struct {
		eventID string
	}
	eventID := getEventID(t)
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{
			name: "ok",
			args: args{
				eventID: eventID,
			},
			want:    eventID,
			wantErr: false,
		},
		{
			name: "no event",
			args: args{
				eventID: "noeventid",
			},
			want:    eventID,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := r.EventByID(ctx, tt.args.eventID)
			if (err != nil) != tt.wantErr {
				t.Errorf("EventByID() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if err == nil && got.ID != tt.want {
				t.Errorf("EventByID() got eventID = %v, want %v", got.ID, tt.want)
			}
			if tt.name == "no event" && !errors.Is(err, eventssvc.ErrUnknownEvent) {
				t.Errorf("EventByID() want eventssvc.ErrUnknownEvent, got %v", err)
			}
		})
	}
}

func TestEventingRepo_StoreEvent(t *testing.T) {
	r := postgres.NewEventingRepo(db)
	ctx := context.Background()
	type args struct {
		event *eventssvc.Event
	}
	event := &eventssvc.Event{
		ID:               "someid",
		Date:             time.Now().AddDate(0, 0, 5),
		Name:             "some name",
		ResourcesAllowed: []eventssvc.Resource{eventssvc.Gadgets, eventssvc.Paper},
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "ok",
			args: args{
				event,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := r.StoreEvent(ctx, tt.args.event); (err != nil) != tt.wantErr {
				t.Errorf("StoreEvent() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
		deleteEvent(t, event.ID)
	}
}

func TestEventingRepo_PupilByID(t *testing.T) {
	r := postgres.NewEventingRepo(db)
	ctx := context.Background()

	eID, deleteE := createEvent(t, &eventssvc.Event{
		ID:               "someid",
		Date:             time.Now().AddDate(0, 1, 0),
		Name:             "some name",
		ResourcesAllowed: []eventssvc.Resource{eventssvc.Gadgets},
	})
	defer deleteE()

	pID, deleteP := createPupil(t, &eventssvc.Pupil{
		ID:        "someid",
		FirstName: "fn",
		LastName:  "sn",
	}, eventssvc.Class{
		Letter:     "A",
		DateFormed: time.Now().AddDate(-2, 0, 0),
	})
	defer deleteP()

	oldPID, deleteOldP := createPupil(t, &eventssvc.Pupil{
		ID:        "someanotherid",
		FirstName: "fn",
		LastName:  "sn",
	}, eventssvc.Class{
		Letter:     "A",
		DateFormed: time.Now().AddDate(-13, 0, 0),
	})
	defer deleteOldP()

	type args struct {
		pupilID string
		eventID string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "ok",
			args: args{
				pupilID: pID,
				eventID: eID,
			},
			wantErr: false,
		},
		{
			name: "pupil is too old to participate in the event",
			args: args{
				pupilID: oldPID,
				eventID: eID,
			},
			wantErr: true,
		},
		{
			name: "no event",
			args: args{
				pupilID: pID,
				eventID: "noevent",
			},
			wantErr: true,
		},
		{
			name: "no pupil",
			args: args{
				pupilID: "nopupil",
				eventID: eID,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := r.PupilByID(ctx, tt.args.pupilID, tt.args.eventID)
			if (err != nil) != tt.wantErr {
				t.Errorf("PupilByID() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.name == "pupil is too old to participate in the event" && !errors.Is(err, eventing.ErrNoEventPupil) {
				t.Errorf("PupilByID() error = %v, want eventing.ErrNoEventPupil", err)
			}
			if tt.name == "no event" && !errors.Is(err, eventssvc.ErrUnknownEvent) {
				t.Errorf("PupilByID() error = %v, want eventssvc.ErrUnknownEvent", err)
			}
			if tt.name == "no pupil" && !errors.Is(err, eventssvc.ErrUnknownPupil) {
				t.Errorf("PupilByID() error = %v, want eventssvc.ErrUnknownPupil", err)
			}
		})
	}
}

func TestEventingRepo_EventPupils(t *testing.T) {
	r := postgres.NewEventingRepo(db)
	ctx := context.Background()
	eID := getEventID(t)

	type args struct {
		eventID string
		filters eventing.EventPupilFilters
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
			name: "without the text search",
			args: args{
				eventID: eID,
				filters: eventing.EventPupilFilters{
					NameAndClass: "",
				},
				sortBy: sorting.Plastic,
				amount: 10,
				skip:   0,
			},
			wantErr: false,
		},
		{
			name: "with the text search",
			args: args{
				eventID: eID,
				filters: eventing.EventPupilFilters{
					NameAndClass: "ro 7",
				},
				sortBy: sorting.Paper,
				amount: 50,
				skip:   0,
			},
			wantErr: false,
		},
		{
			name: "invalid value in the text search",
			args: args{
				eventID: eID,
				filters: eventing.EventPupilFilters{
					NameAndClass: "ro&7",
				},
				sortBy: sorting.Paper,
				amount: 50,
				skip:   0,
			},
			wantErr: false,
		},
		{
			name: "skip more than the total amount of the pupils in the db",
			args: args{
				eventID: eID,
				filters: eventing.EventPupilFilters{},
				sortBy:  sorting.Plastic,
				amount:  150,
				skip:    5000,
			},
			wantErr: false,
		},
		{
			name: "invalid event id",
			args: args{
				eventID: "wrongeventid",
				filters: eventing.EventPupilFilters{},
				sortBy:  sorting.Plastic,
				amount:  150,
				skip:    0,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, _, err := r.EventPupils(ctx, tt.args.eventID, tt.args.filters, tt.args.sortBy, tt.args.amount,
				tt.args.skip)

			if (err != nil) != tt.wantErr {
				t.Errorf("EventPupils() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestEventingRepo_EventClasses(t *testing.T) {
	r := postgres.NewEventingRepo(db)
	ctx := context.Background()
	eID := getEventID(t)

	type args struct {
		eventID string
		filters eventing.EventClassFilters
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
			name: "no classes specified",
			args: args{
				eventID: eID,
				filters: eventing.EventClassFilters{},
				sortBy:  sorting.Plastic,
				amount:  50,
				skip:    0,
			},
			wantErr: false,
		},
		{
			name: "a class is specified",
			args: args{
				eventID: eID,
				filters: eventing.EventClassFilters{
					Name: "3b",
				},
				sortBy: sorting.Plastic,
				amount: 50,
				skip:   0,
			},
			wantErr: false,
		},
		{
			name: "classes are specified",
			args: args{
				eventID: eID,
				filters: eventing.EventClassFilters{
					Name: "3",
				},
				sortBy: sorting.Plastic,
				amount: 50,
				skip:   0,
			},
			wantErr: false,
		},
		{
			name: "skip is more than the total amount of the classes in the db",
			args: args{
				eventID: eID,
				filters: eventing.EventClassFilters{},
				sortBy:  sorting.Plastic,
				amount:  50,
				skip:    999,
			},
			wantErr: false,
		},
		{
			name: "invalid class name",
			args: args{
				eventID: eID,
				filters: eventing.EventClassFilters{
					Name: "3b3",
				},
				sortBy: sorting.Plastic,
				amount: 50,
				skip:   0,
			},
			wantErr: false,
		},
		{
			name: "invalid event id",
			args: args{
				eventID: "wrongeventid",
				filters: eventing.EventClassFilters{
					Name: "",
				},
				sortBy: sorting.Plastic,
				amount: 50,
				skip:   0,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, _, err := r.EventClasses(ctx, tt.args.eventID, tt.args.filters,
				tt.args.sortBy, tt.args.amount, tt.args.skip)

			if (err != nil) != tt.wantErr {
				t.Errorf("EventClasses() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func createEvent(t *testing.T, event *eventssvc.Event) (string, func()) {
	t.Helper()
	q := "insert into event (id, name, date, resources_allowed)\n\tvalues ($1, $2, $3, $4::text[]::resource[]);"
	_, err := db.Exec(context.Background(), q, event.ID, event.Name, event.Date,
		eventssvc.ResourceSliceToStringSlice(event.ResourcesAllowed))
	if err != nil {
		t.Fatalf("prepare db: %v", err)
	}
	return event.ID, func() {
		q := "delete from event where id = $1"
		_, err := db.Exec(context.Background(), q, event.ID)
		if err != nil {
			t.Fatalf("clean db: %v", err)
		}
	}
}

func createPupil(t *testing.T, p *eventssvc.Pupil, c eventssvc.Class) (string, func()) {
	t.Helper()
	q := `
		insert into pupil (id, first_name, last_name, class_letter, class_date_formed)
		values ($1, $2, $3, $4, $5);
	`
	if _, err := db.Exec(context.Background(), q, p.ID, p.FirstName, p.LastName, c.Letter, c.DateFormed); err != nil {
		t.Fatalf("prepare db: %v", err)
	}
	return p.ID, func() {
		_, err := db.Exec(context.Background(), `delete from pupil where pupil.id = $1`, p.ID)

		if err != nil {
			t.Fatalf("clean db: %v", err)
		}
	}
}

func deleteEvent(t *testing.T, eventID string) {
	t.Helper()
	if _, err := db.Exec(context.Background(), "delete from event where id=$1", eventID); err != nil {
		t.Fatalf("prepare db: deleteEvent error: %v", err)
	}
}

func getPupilID(t *testing.T) string {
	t.Helper()
	var pupilID string
	if err := db.QueryRow(context.Background(), `select id from pupil`).Scan(&pupilID); err != nil {
		t.Fatalf("prepare db: %v", err)
	}
	return pupilID
}

func getEventID(t *testing.T) string {
	t.Helper()
	var eventID string
	if err := db.QueryRow(context.Background(), `select id from event`).Scan(&eventID); err != nil {
		t.Fatalf("prepare db: %v", err)
	}
	return eventID
}

func removePupilResources(t *testing.T, pupilID string, eventID string) {
	t.Helper()
	if _, err := db.Exec(context.Background(), `delete from resources where pupil_id = $1 and event_id = $2;`, pupilID,
		eventID); err != nil {
		t.Fatalf("prepare db: %v", err)
	}
}
