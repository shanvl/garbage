package sorting

import (
	"testing"
)

func TestSort_IsDate(t *testing.T) {
	tests := []struct {
		name string
		s    By
		want bool
	}{
		{
			name: "by date",
			s:    DateDes,
			want: true,
		},
		{
			name: "by plastic",
			s:    Plastic,
			want: false,
		},
		{
			name: "by name",
			s:    NameAsc,
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.s.IsDate(); got != tt.want {
				t.Errorf("IsDate() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSort_IsName(t *testing.T) {
	tests := []struct {
		name string
		s    By
		want bool
	}{
		{
			name: "by name",
			s:    NameAsc,
			want: true,
		},
		{
			name: "by plastic",
			s:    Plastic,
			want: false,
		},
		{
			name: "by date",
			s:    DateAsc,
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.s.IsName(); got != tt.want {
				t.Errorf("IsName() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSort_IsResources(t *testing.T) {
	tests := []struct {
		name string
		s    By
		want bool
	}{
		{
			name: "by gadgets",
			s:    Gadgets,
			want: true,
		},
		{
			name: "by date",
			s:    DateDes,
			want: false,
		},
		{
			name: "by name",
			s:    NameAsc,
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.s.IsResources(); got != tt.want {
				t.Errorf("IsResources() = %v, want %v", got, tt.want)
			}
		})
	}
}
