package garbage

import (
	"reflect"
	"testing"
	"time"
)

func TestClass_NameFromDate(t *testing.T) {
	type fields struct {
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
			c := Class{
				DateFormed: tt.fields.Formed,
				Letter:     tt.fields.Letter,
			}
			got, err := c.NameOnDate(tt.args.date)
			if (err != nil) != tt.wantErr {
				t.Errorf("NameOnDate() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("NameOnDate() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func newDate(year int, month int, day int, hour int) time.Time {
	return time.Date(year, time.Month(month), day, hour, 0, 0, 0, time.UTC)
}

func TestClassFromClassName(t *testing.T) {
	date := newDate(2010, 10, 10, 1)

	type args struct {
		className string
		date      time.Time
	}
	tests := []struct {
		name    string
		args    args
		want    Class
		wantErr bool
	}{
		{
			name: "valid class",
			args: args{
				className: "3B",
				date:      date,
			},
			want:    Class{"b", newDate(2008, 9, 1, 0)},
			wantErr: false,
		},
		{
			name: "class with non-alphanumeric chars",
			args: args{
				className: "  3- -* B  ",
				date:      date,
			},
			want:    Class{"b", newDate(2008, 9, 1, 0)},
			wantErr: false,
		},
		{
			name: "only class number",
			args: args{
				className: "3",
				date:      date,
			},
			want:    Class{},
			wantErr: true,
		},
		{
			name: "digit after letter in class",
			args: args{
				className: "3Б1",
				date:      date,
			},
			want:    Class{},
			wantErr: true,
		},
		{
			name: "2 letters in class",
			args: args{
				className: "10ББ",
				date:      date,
			},
			want:    Class{},
			wantErr: true,
		},
		{
			name: "0 in class",
			args: args{
				className: "0Б",
				date:      date,
			},
			want:    Class{},
			wantErr: true,
		},
		{
			name: "letter before digit in class",
			args: args{
				className: "Б10",
				date:      date,
			},
			want:    Class{},
			wantErr: true,
		},
		{
			name: "class number > 11",
			args: args{
				className: "12Б",
				date:      date,
			},
			want:    Class{},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ClassFromClassName(tt.args.className, tt.args.date)
			if (err != nil) != tt.wantErr {
				t.Errorf("ClassFromClassName() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("ClassFromClassName() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestParseClassName(t *testing.T) {
	date := newDate(2010, 10, 10, 1)
	type args struct {
		className string
		date      time.Time
	}
	tests := []struct {
		name           string
		args           args
		wantLetter     string
		wantDateFormed time.Time
		wantErr        bool
	}{
		{
			name: "3B",
			args: args{
				className: "3B",
				date:      date,
			},
			wantLetter:     "b",
			wantDateFormed: newDate(2008, 9, 1, 0),
			wantErr:        false,
		},
		{
			name: "3  B",
			args: args{
				className: "3  B",
				date:      date,
			},
			wantLetter:     "b",
			wantDateFormed: newDate(2008, 9, 1, 0),
			wantErr:        false,
		},
		{
			name: "B",
			args: args{
				className: "B",
				date:      date,
			},
			wantLetter:     "b",
			wantDateFormed: time.Time{},
			wantErr:        false,
		},
		{
			name: "3",
			args: args{
				className: "3",
				date:      date,
			},
			wantLetter:     "",
			wantDateFormed: newDate(2008, 9, 1, 0),
			wantErr:        false,
		},
		{
			name: "B3",
			args: args{
				className: "B3",
				date:      date,
			},
			wantLetter:     "",
			wantDateFormed: time.Time{},
			wantErr:        true,
		},
		{
			name: "3BB",
			args: args{
				className: "3BB",
				date:      date,
			},
			wantLetter:     "",
			wantDateFormed: time.Time{},
			wantErr:        true,
		},
		{
			name: "0B",
			args: args{
				className: "0B",
				date:      date,
			},
			wantLetter:     "",
			wantDateFormed: time.Time{},
			wantErr:        true,
		},
		{
			name: "empty",
			args: args{
				className: "",
				date:      date,
			},
			wantLetter:     "",
			wantDateFormed: time.Time{},
			wantErr:        true,
		},
		{
			name: "12B",
			args: args{
				className: "12B",
				date:      date,
			},
			wantLetter:     "",
			wantDateFormed: time.Time{},
			wantErr:        true,
		},
		{
			name: "3 BB",
			args: args{
				className: "3 BB",
				date:      date,
			},
			wantLetter:     "",
			wantDateFormed: time.Time{},
			wantErr:        true,
		},
		{
			name: "B3B",
			args: args{
				className: "B3B",
				date:      date,
			},
			wantLetter:     "",
			wantDateFormed: time.Time{},
			wantErr:        true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotLetter, gotDateFormed, err := ParseClassName(tt.args.className, tt.args.date)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseClassName() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotLetter != tt.wantLetter {
				t.Errorf("ParseClassName() gotLetter = %v, want %v", gotLetter, tt.wantLetter)
			}
			if !reflect.DeepEqual(gotDateFormed, tt.wantDateFormed) {
				t.Errorf("ParseClassName() gotDateFormed = %v, want %v", gotDateFormed, tt.wantDateFormed)
			}
		})
	}
}
