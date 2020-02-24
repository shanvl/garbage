package eventing

import (
	"testing"
)

func TestSort_IsValid(t *testing.T) {
	tests := []struct {
		name string
		s    SortBy
		want bool
	}{
		{
			name: "valid input",
			s:    "plastic",
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
