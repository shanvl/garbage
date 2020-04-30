package garbage

import (
	"testing"
	"time"
)

func TestClass_NameFromDate(t *testing.T) {
	type fields struct {
		Formed int
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
				Formed: 2020,
				Letter: "Б",
			},
			args:    args{date: newDate(2020, 8, 31, 23)},
			want:    "",
			wantErr: true,
		},
		{
			name: "first day of the first class",
			fields: fields{
				Formed: 2020,
				Letter: "Б",
			},
			args:    args{date: newDate(2020, 9, 1, 1)},
			want:    "1Б",
			wantErr: false,
		},
		{
			name: "date is in the next calender year, but in the same school year",
			fields: fields{
				Formed: 2021,
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
				Formed: 2021,
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
				Formed: 2020,
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
				Formed: 2020,
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
				YearFormed: tt.fields.Formed,
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

func TestParseClassName(t *testing.T) {
	const className = "3Б"
	date := newDate(2010, 10, 10, 1)

	type args struct {
		className string
		date      time.Time
	}
	tests := []struct {
		name       string
		args       args
		wantLetter string
		wantFormed int
		wantErr    bool
	}{
		{
			name: "valid class",
			args: args{
				className: className,
				date:      date,
			},
			wantLetter: "Б",
			wantFormed: 2008,
			wantErr:    false,
		},
		{
			name: "class with non-alphanumeric chars",
			args: args{
				className: "  3- -* Б  ",
				date:      date,
			},
			wantLetter: "Б",
			wantFormed: 2008,
			wantErr:    false,
		},
		{
			name: "digit after letter in class",
			args: args{
				className: "3Б1",
				date:      date,
			},
			wantLetter: "",
			wantFormed: 0,
			wantErr:    true,
		},
		{
			name: "2 letters in class",
			args: args{
				className: "10ББ",
				date:      date,
			},
			wantLetter: "",
			wantFormed: 0,
			wantErr:    true,
		},
		{
			name: "0 in class",
			args: args{
				className: "0Б",
				date:      date,
			},
			wantLetter: "",
			wantFormed: 0,
			wantErr:    true,
		},
		{
			name: "letter before digit in class",
			args: args{
				className: "Б10",
				date:      date,
			},
			wantLetter: "",
			wantFormed: 0,
			wantErr:    true,
		},
		{
			name: "class number > 11",
			args: args{
				className: "12Б",
				date:      date,
			},
			wantLetter: "",
			wantFormed: 0,
			wantErr:    true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotLetter, gotFormed, err := ParseClassName(tt.args.className, tt.args.date)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseClassName() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotLetter != tt.wantLetter {
				t.Errorf("ParseClassName() gotLetter = %v, want %v", gotLetter, tt.wantLetter)
			}
			if gotFormed != tt.wantFormed {
				t.Errorf("ParseClassName() gotFormed = %v, want %v", gotFormed, tt.wantFormed)
			}
		})
	}
}
