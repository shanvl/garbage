package garbage

import (
	"testing"
	"time"
)

func TestClass_NameFromDate(t *testing.T) {
	type fields struct {
		ID     ClassID
		Formed time.Time
		Letter string
	}
	type args struct {
		date time.Time
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    string
		wantErr bool
	}{
		{
			name: "class wasn't formed yet",
			fields: fields{
				ID:     "id",
				Formed: newDate(2020, 9, 1, 0),
				Letter: "Б",
			},
			args:    args{date: newDate(2020, 8, 31, 23)},
			want:    "",
			wantErr: true,
		},
		{
			name: "first day of the first class",
			fields: fields{
				ID:     "id",
				Formed: newDate(2020, 9, 1, 0),
				Letter: "Б",
			},
			args:    args{date: newDate(2020, 9, 1, 1)},
			want:    "1Б",
			wantErr: false,
		},
		{
			name: "date is in the next calender year, but in the same school year",
			fields: fields{
				ID:     "id",
				Formed: newDate(2021, 9, 1, 0),
				Letter: "Б",
			},
			args: args{
				date: newDate(2022, 2, 1, 0),
			},
			want:    "1Б",
			wantErr: false,
		},
		{
			name: "5,5 years after the class was formed",
			fields: fields{
				ID:     "id",
				Formed: newDate(2021, 9, 1, 0),
				Letter: "Б",
			},
			args: args{
				date: newDate(2026, 2, 1, 0),
			},
			want:    "5Б",
			wantErr: false,
		},
		{
			name: "10 years after the class was formed",
			fields: fields{
				ID:     "id",
				Formed: newDate(2020, 9, 1, 0),
				Letter: "Б",
			},
			args: args{
				date: newDate(2030, 9, 1, 0),
			},
			want:    "11Б",
			wantErr: false,
		},
		{
			name: "class is already graduated",
			fields: fields{
				ID:     "id",
				Formed: newDate(2020, 9, 1, 0),
				Letter: "Б",
			},
			args:    args{date: newDate(2032, 9, 1, 0)},
			want:    "",
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Class{
				ID:     tt.fields.ID,
				Formed: tt.fields.Formed,
				Letter: tt.fields.Letter,
			}
			got, err := c.NameFromDate(tt.args.date)
			if (err != nil) != tt.wantErr {
				t.Errorf("NameFromDate() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("NameFromDate() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func newDate(year int, month int, day int, hour int) time.Time {
	return time.Date(year, time.Month(month), day, hour, 0, 0, 0, time.UTC)
}
