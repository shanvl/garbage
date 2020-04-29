package postgres_test

import (
	"context"
	"reflect"
	"testing"

	"github.com/shanvl/garbage-events-service/internal/garbage"
	"github.com/shanvl/garbage-events-service/internal/postgres"
	"github.com/shanvl/garbage-events-service/internal/usecases/schooling"
)

func TestSchoolingRepo_StorePupil(t *testing.T) {
	var r = postgres.NewSchoolingRepo(db)
	var ctx = context.Background()
	type args struct {
		pupil *schooling.Pupil
	}
	tests := []struct {
		name    string
		args    args
		want    garbage.PupilID
		wantErr bool
	}{
		{
			name: "ok case",
			args: args{
				pupil: &schooling.Pupil{
					Pupil: garbage.Pupil{
						ID:        "aaa",
						FirstName: "fname",
						LastName:  "lname",
					},
					Class: garbage.Class{
						Letter:     "A",
						YearFormed: 2015,
					},
				},
			},
			want:    "aaa",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := r.StorePupil(ctx, tt.args.pupil)
			if (err != nil) != tt.wantErr {
				t.Errorf("StorePupil() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("StorePupil() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSchoolingRepo_PupilByID(t *testing.T) {
	var r = postgres.NewSchoolingRepo(db)
	var ctx = context.Background()
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
			name: "ok case",
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
				YearFormed: 2015,
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
				YearFormed: 2016,
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
				YearFormed: 2017,
			},
		},
	}
	t.Helper()
	stmt, err := db.Prepare(`
	insert into pupil (id, first_name, last_name, class_letter, class_year_formed)
	values ($1, $2, $3, $4, $5);`)
	if err != nil {
		t.Fatalf("prepare db: %v", err)
	}
	defer stmt.Close()
	for _, p := range pp {
		if _, err := stmt.Exec(p.ID, p.FirstName, p.LastName, p.Class.Letter, p.Class.YearFormed); err != nil {
			t.Fatalf("prepare db: %v", err)
		}
	}
	return pp, func() {
		_, err := db.Exec(`delete from pupil where pupil.id in ($1, $2, $3)`, pp[0].ID, pp[1].ID, pp[2].ID)
		if err != nil {
			t.Fatalf("clean db: %v", err)
		}
	}
}
