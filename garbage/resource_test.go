package garbage

import "testing"

func TestResource_IsValid(t *testing.T) {
	tests := []struct {
		name string
		r    Resource
		want bool
	}{
		{
			name: "known resource",
			r:    "plastic",
			want: true,
		},
		{
			name: "unknown resource",
			r:    "unknown",
			want: false,
		},
		{
			name: "empty string",
			r:    "",
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.r.IsKnown(); got != tt.want {
				t.Errorf("IsKnown() = %v, want %v", got, tt.want)
			}
		})
	}
}
