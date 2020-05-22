package postgres_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/shanvl/garbage-events-service/internal/garbage"
	"github.com/shanvl/garbage-events-service/internal/postgres"
	"github.com/shanvl/garbage-events-service/internal/sorting"
	"github.com/shanvl/garbage-events-service/internal/usecases/eventing"
)

func TestEventingRepo_ChangePupilResources(t *testing.T) {
	r := postgres.NewEventingRepo(db)
	ctx := context.Background()
	pupilID, eventID := getPupilID(t), getEventID(t)
	removePupilResources(t, pupilID, eventID)
	resources := garbage.ResourceMap{
		garbage.Plastic: 10,
		garbage.Paper:   15,
	}
	type args struct {
		eventID   garbage.EventID
		pupilID   garbage.PupilID
		resources garbage.ResourceMap
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
		eventID garbage.EventID
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
		eventID garbage.EventID
	}
	eventID := getEventID(t)
	tests := []struct {
		name    string
		args    args
		want    garbage.EventID
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
			if tt.name == "no event" && !errors.Is(err, garbage.ErrUnknownEvent) {
				t.Errorf("EventByID() want garbage.ErrUnknownEvent, got %v", err)
			}
		})
	}
}

func TestEventingRepo_StoreEvent(t *testing.T) {
	r := postgres.NewEventingRepo(db)
	ctx := context.Background()
	type args struct {
		event *garbage.Event
	}
	event := &garbage.Event{
		ID:               "someid",
		Date:             time.Now().AddDate(0, 0, 5),
		Name:             "some name",
		ResourcesAllowed: []garbage.Resource{garbage.Gadgets, garbage.Paper},
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "invalid resource",
			args: args{
				&garbage.Event{
					ID:               "someid",
					Date:             time.Now().AddDate(0, 0, 5),
					Name:             "some name",
					ResourcesAllowed: []garbage.Resource{garbage.Gadgets, "invalid resource"},
				},
			},
			wantErr: true,
		},
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

	eID, deleteE := createEvent(t, &garbage.Event{
		ID:               "someid",
		Date:             time.Now().AddDate(0, 1, 0),
		Name:             "some name",
		ResourcesAllowed: []garbage.Resource{garbage.Gadgets},
	})
	defer deleteE()

	pID, deleteP := createPupil(t, &garbage.Pupil{
		ID:        "someid",
		FirstName: "fn",
		LastName:  "sn",
	}, garbage.Class{
		Letter:     "A",
		DateFormed: time.Now().AddDate(-2, 0, 0),
	})
	defer deleteP()

	oldPID, deleteOldP := createPupil(t, &garbage.Pupil{
		ID:        "someanotherid",
		FirstName: "fn",
		LastName:  "sn",
	}, garbage.Class{
		Letter:     "A",
		DateFormed: time.Now().AddDate(-13, 0, 0),
	})
	defer deleteOldP()

	type args struct {
		pupilID garbage.PupilID
		eventID garbage.EventID
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
			if tt.name == "no event" && !errors.Is(err, garbage.ErrUnknownEvent) {
				t.Errorf("PupilByID() error = %v, want garbage.ErrUnknownEvent", err)
			}
			if tt.name == "no pupil" && !errors.Is(err, garbage.ErrUnknownPupil) {
				t.Errorf("PupilByID() error = %v, want garbage.ErrUnknownPupil", err)
			}
		})
	}
}

func TestEventingRepo_EventPupils(t *testing.T) {
	r := postgres.NewEventingRepo(db)
	ctx := context.Background()
	eID := getEventID(t)

	type args struct {
		eventID garbage.EventID
		filters eventing.EventPupilsFilters
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
			name: "without text search",
			args: args{
				eventID: eID,
				filters: eventing.EventPupilsFilters{
					NameAndClass: "",
				},
				sortBy: "",
				amount: 150,
				skip:   0,
			},
			wantErr: false,
		},
		{
			name: "with text search",
			args: args{
				eventID: eID,
				filters: eventing.EventPupilsFilters{
					NameAndClass: "a 3",
				},
				sortBy: "",
				amount: 150,
				skip:   0,
			},
			wantErr: false,
		},
		{
			name: "skip more than the total amount of the pupils in the db",
			args: args{
				eventID: eID,
				filters: eventing.EventPupilsFilters{
					NameAndClass: "",
				},
				sortBy: "",
				amount: 150,
				skip:   5000,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, gotTotal, err := r.EventPupils(ctx, tt.args.eventID, tt.args.filters, tt.args.sortBy, tt.args.amount,
				tt.args.skip)

			if (err != nil) != tt.wantErr || gotTotal == 0 {
				t.Errorf("EventPupils() error = %v, wantErr %v, gotTotal %d, ", err, tt.wantErr, gotTotal)
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
		eventID garbage.EventID
		filters eventing.EventClassesFilters
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
				filters: eventing.EventClassesFilters{},
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
				filters: eventing.EventClassesFilters{
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
				filters: eventing.EventClassesFilters{
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
				filters: eventing.EventClassesFilters{},
				sortBy:  sorting.Plastic,
				amount:  50,
				skip:    999,
			},
			wantErr: false,
		},
		{
			name: "ivalid class name",
			args: args{
				eventID: eID,
				filters: eventing.EventClassesFilters{
					Name: "3b3",
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

func createEvent(t *testing.T, event *garbage.Event) (garbage.EventID, func()) {
	t.Helper()
	q := "insert into event (id, name, date, resources_allowed)\n\tvalues ($1, $2, $3, $4::text[]::resource[]);"
	_, err := db.Exec(context.Background(), q, event.ID, event.Name, event.Date,
		garbage.ResourceSliceToStringSlice(event.ResourcesAllowed))
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

func createPupil(t *testing.T, p *garbage.Pupil, c garbage.Class) (garbage.PupilID, func()) {
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

func deleteEvent(t *testing.T, eventID garbage.EventID) {
	t.Helper()
	if _, err := db.Exec(context.Background(), "delete from event where id=$1", eventID); err != nil {
		t.Fatalf("prepare db: deleteEvent error: %v", err)
	}
}

func getPupilID(t *testing.T) garbage.PupilID {
	t.Helper()
	var pupilID garbage.PupilID
	if err := db.QueryRow(context.Background(), `select id from pupil`).Scan(&pupilID); err != nil {
		t.Fatalf("prepare db: %v", err)
	}
	return pupilID
}

func getEventID(t *testing.T) garbage.EventID {
	t.Helper()
	var eventID garbage.EventID
	if err := db.QueryRow(context.Background(), `select id from event`).Scan(&eventID); err != nil {
		t.Fatalf("prepare db: %v", err)
	}
	return eventID
}

func removePupilResources(t *testing.T, pupilID garbage.PupilID, eventID garbage.EventID) {
	t.Helper()
	if _, err := db.Exec(context.Background(), `delete from resources where pupil_id = $1 and event_id = $2;`, pupilID,
		eventID); err != nil {
		t.Fatalf("prepare db: %v", err)
	}
}
