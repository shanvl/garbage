package postgres_test

import (
	"context"
	"errors"
	"testing"

	"github.com/shanvl/garbage-events-service/internal/garbage"
	"github.com/shanvl/garbage-events-service/internal/postgres"
	"github.com/shanvl/garbage-events-service/internal/usecases/eventing"
)

func TestEventingRepo_ChangePupilResources(t *testing.T) {
	r := postgres.NewEventingRepo(db)
	ctx := context.Background()
	pupilID, eventID := getPupilID(t), getEventID(t)
	removePupilResources(t, pupilID, eventID)
	resources := map[garbage.Resource]int{
		garbage.Plastic: 10,
		garbage.Paper:   15,
	}
	type args struct {
		eventID   garbage.EventID
		pupilID   garbage.PupilID
		resources map[garbage.Resource]int
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

func getPupilID(t *testing.T) garbage.PupilID {
	t.Helper()
	var pupilID garbage.PupilID
	if err := db.QueryRow(`select id from pupil`).Scan(&pupilID); err != nil {
		t.Fatalf("prepare db: %v", err)
	}
	return pupilID
}

func getEventID(t *testing.T) garbage.EventID {
	t.Helper()
	var eventID garbage.EventID
	if err := db.QueryRow(`select id from event`).Scan(&eventID); err != nil {
		t.Fatalf("prepare db: %v", err)
	}
	return eventID
}

func removePupilResources(t *testing.T, pupilID garbage.PupilID, eventID garbage.EventID) {
	t.Helper()
	if _, err := db.Exec(`delete from resources where pupil_id = $1 and event_id = $2;`, pupilID,
		eventID); err != nil {
		t.Fatalf("prepare db: %v", err)
	}
}
