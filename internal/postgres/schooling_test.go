package postgres_test

import (
	"context"
	"reflect"
	"testing"
	"time"

	"github.com/shanvl/garbage-events-service/internal/garbage"
	"github.com/shanvl/garbage-events-service/internal/postgres"
	"github.com/shanvl/garbage-events-service/internal/usecases/schooling"
)

func TestSchoolingRepo_StorePupil(t *testing.T) {
	r := postgres.NewSchoolingRepo(db)
	ctx := context.Background()
	type args struct {
		pupil *schooling.Pupil
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "ok",
			args: args{
				pupil: &schooling.Pupil{
					Pupil: garbage.Pupil{
						ID:        "111",
						FirstName: "fname",
						LastName:  "lname",
					},
					Class: garbage.Class{
						Letter:     "A",
						DateFormed: classDateFromYear(2015),
					},
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := r.StorePupil(ctx, tt.args.pupil)
			if (err != nil) != tt.wantErr {
				t.Errorf("StorePupil() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if err := r.RemovePupils(ctx, []garbage.PupilID{tt.args.pupil.ID}); err != nil {
				t.Fatalf("wasn't able to clean the db after StorePupil(): %v", err)
			}
		})
	}
}

func TestSchoolingRepo_PupilByID(t *testing.T) {
	r := postgres.NewSchoolingRepo(db)
	ctx := context.Background()
	pp, cleanDB := seedPupils(t)
	defer cleanDB()
	type args struct {
		pupilID garbage.PupilID
	}
	tests := []struct {
		name    string
		args    args
		want    *schooling.Pupil
		wantErr bool
	}{
		{
			name: "ok",
			args: args{
				pupilID: pp[0].ID,
			},
			want:    pp[0],
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := r.PupilByID(ctx, tt.args.pupilID)
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

func TestSchoolingRepo_StorePupils(t *testing.T) {
	r := postgres.NewSchoolingRepo(db)
	ctx := context.Background()
	pp := []*schooling.Pupil{
		{
			Pupil: garbage.Pupil{
				ID:        "111",
				FirstName: "fn1",
				LastName:  "ln1",
			},
			Class: garbage.Class{
				Letter:     "a",
				DateFormed: classDateFromYear(2015),
			},
		},
		{
			Pupil: garbage.Pupil{
				ID:        "222",
				FirstName: "fn2",
				LastName:  "ln2",
			},
			Class: garbage.Class{
				Letter:     "b",
				DateFormed: classDateFromYear(2016),
			},
		},
		{
			Pupil: garbage.Pupil{
				ID:        "333",
				FirstName: "fn3",
				LastName:  "ln3",
			},
			Class: garbage.Class{
				Letter:     "c",
				DateFormed: classDateFromYear(2017),
			},
		},
	}
	tests := []struct {
		name    string
		pupils  []*schooling.Pupil
		wantErr bool
	}{
		{
			name:    "ok",
			pupils:  pp,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := r.StorePupils(ctx, tt.pupils)
			if (err != nil) != tt.wantErr {
				t.Errorf("StorePupils() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if err := r.RemovePupils(ctx, []garbage.PupilID{tt.pupils[0].ID, tt.pupils[1].ID,
				tt.pupils[2].ID}); err != nil {
				t.Fatalf("wasn't able to clean db after StorePupils(): %v", err)
			}
		})
	}
}

func TestSchoolingRepo_RemovePupils(t *testing.T) {
	r := postgres.NewSchoolingRepo(db)
	ctx := context.Background()
	pupils, cleanDB := seedPupils(t)
	defer cleanDB()
	ppIDs := make([]garbage.PupilID, len(pupils))
	for i, p := range pupils {
		ppIDs[i] = p.ID
	}
	tests := []struct {
		name     string
		pupilIDs []garbage.PupilID
		wantErr  bool
	}{
		{
			name:     "ok",
			pupilIDs: ppIDs,
			wantErr:  false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := r.RemovePupils(ctx, tt.pupilIDs)
			if (err != nil) != tt.wantErr {
				t.Errorf("RemovePupils() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func seedPupils(t *testing.T) ([]*schooling.Pupil, func()) {
	pp := []*schooling.Pupil{
		{
			Pupil: garbage.Pupil{
				ID:        "111",
				FirstName: "fn1",
				LastName:  "ln1",
			},
			Class: garbage.Class{
				Letter:     "a",
				DateFormed: classDateFromYear(2015),
			},
		},
		{
			Pupil: garbage.Pupil{
				ID:        "222",
				FirstName: "fn2",
				LastName:  "ln2",
			},
			Class: garbage.Class{
				Letter:     "b",
				DateFormed: classDateFromYear(2016),
			},
		},
		{
			Pupil: garbage.Pupil{
				ID:        "333",
				FirstName: "fn3",
				LastName:  "ln3",
			},
			Class: garbage.Class{
				Letter:     "c",
				DateFormed: classDateFromYear(2017),
			},
		},
	}
	t.Helper()
	q := `
		insert into pupil (id, first_name, last_name, class_letter, class_date_formed)
		values ($1, $2, $3, $4, $5);
	`
	for _, p := range pp {
		if _, err := db.Exec(context.Background(), q, p.ID, p.FirstName, p.LastName, p.Class.Letter,
			p.Class.DateFormed); err != nil {

			t.Fatalf("prepare db: %v", err)
		}
	}
	return pp, func() {
		_, err := db.Exec(context.Background(), `delete from pupil where pupil.id in ($1, $2, $3)`, pp[0].ID,
			pp[1].ID, pp[2].ID)

		if err != nil {
			t.Fatalf("clean db: %v", err)
		}
	}
}

func classDateFromYear(year int) time.Time {
	return time.Date(year, 9, 1, 0, 0, 0, 0, time.UTC)
}
