package postgres

import (
	"testing"
	"time"
)

func Test_prepareTextSearchQuery(t *testing.T) {
	date := time.Date(2020, 10, 10, 10, 10, 10, 10, time.UTC)
	type args struct {
		q string
		t time.Time
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "empty string",
			args: args{
				q: "",
				t: time.Now(),
			},
			want: "",
		},
		{
			name: "string with no class name",
			args: args{
				q: "iv id 213bas",
				t: time.Now(),
			},
			want: "iv:* & id:* & 213bas:*",
		},
		{
			name: "string with a class name",
			args: args{
				q: "iv id 3B 213bas",
				t: date,
			},
			want: "iv:* & id:* & (3B:* | 2018b:*) & 213bas:*",
		},
		{
			name: "string with a class name without a letter",
			args: args{
				q: "iv id 3",
				t: date,
			},
			want: "iv:* & id:* & (3:* | 2018:*)",
		},
		{
			name: "string with invalid symbols",
			args: args{
				q: "iv i&d 3",
				t: date,
			},
			want: "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := prepareTextSearchQuery(tt.args.q, tt.args.t); got != tt.want {
				t.Errorf("prepareTextSearchQuery() = %v, want %v", got, tt.want)
			}
		})
	}
}
