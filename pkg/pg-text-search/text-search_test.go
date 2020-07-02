package pgtextsearch

import (
	"testing"
)

func Test_PrepareQuery(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		{
			name:  "empty input",
			input: "",
			want:  "",
		},
		{
			name:  "2 words",
			input: "some input",
			want:  "some:* & input:*",
		},
		{
			name:  "5 words",
			input: "five words of sample input",
			want:  "five:* & words:* & of:* & sample:* & input:*",
		},
		{
			name:  "invalid input",
			input: "some & input",
			want:  "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := PrepareQuery(tt.input); got != tt.want {
				t.Errorf("prepareTextSearchQuery() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIsValidInput(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  bool
	}{
		{
			name:  "invalid input &",
			input: "some & input",
			want:  false,
		},
		{
			name:  "invalid input \"",
			input: "some input \"",
			want:  false,
		},
		{
			name:  "invalid input @",
			input: "some @input",
			want:  false,
		},
		{
			name:  "invalid input :",
			input: "some :input",
			want:  false,
		},
		{
			name:  "empty input",
			input: "",
			want:  true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsValidInput(tt.input); got != tt.want {
				t.Errorf("IsValidInput() = %v, want %v", got, tt.want)
			}
		})
	}
}
