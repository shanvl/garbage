package sorting

import (
	"testing"
)

func TestSort_IsValid(t *testing.T) {
	tests := []struct {
		name string
		s    By
		want bool
	}{
		{
			name: "valid input",
			s:    Plastic,
			want: true,
		},
		{
			name: "empty string",
			s:    "",
			want: false,
		},
		{
			name: "invalid input",
			s:    "invalid",
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.s.IsValid(); got != tt.want {
				t.Errorf("IsValid() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSort_IsForEventPupils(t *testing.T) {
	tests := []struct {
		name string
		s    By
		want bool
	}{
		{
			name: "valid input",
			s:    Plastic,
			want: true,
		},
		{
			name: "empty string",
			s:    "",
			want: false,
		},
		{
			name: "invalid sorting",
			s:    DateDes,
			want: false,
		},
		{
			name: "invalid input",
			s:    "invalid",
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.s.IsForEventPupils(); got != tt.want {
				t.Errorf("IsForEventPupils() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSort_IsForEventClasses(t *testing.T) {
	tests := []struct {
		name string
		s    By
		want bool
	}{
		{
			name: "valid input",
			s:    Plastic,
			want: true,
		},
		{
			name: "empty string",
			s:    "",
			want: false,
		},
		{
			name: "invalid sorting",
			s:    DateDes,
			want: false,
		},
		{
			name: "invalid input",
			s:    "invalid",
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.s.IsForEventClasses(); got != tt.want {
				t.Errorf("IsForEventClasses() = %v, want %v", got, tt.want)
			}
		})
	}
}
