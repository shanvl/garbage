package valid

import (
	"reflect"
	"testing"
)

func TestErrValidation_Add(t *testing.T) {
	tests := []struct {
		name string
		args [][]string
		want map[string]string
	}{
		{
			name: "no args",
			args: [][]string{},
			want: map[string]string{},
		},
		{
			name: "2 args",
			args: [][]string{
				{"name", "wrong name"},
				{"age", "wrong age"},
			},
			want: map[string]string{"name": "wrong name", "age": "wrong age"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := EmptyError()
			for _, arg := range tt.args {
				e.Add(arg[0], arg[1])
			}
			if got := e.fields; !reflect.DeepEqual(got, tt.want) {
				t.Errorf("e.fields = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestErrValidation_Error(t *testing.T) {
	// We should consider in the tests that a map's values order is always random,
	// that's why there are 2 cases of "want" — they cover every possible combination of map's values.
	// Another way to solve this would be to sorting map's keys in the Error function itself,
	// before returning them as a string, but I don't see any practical value in it for a client
	tests := []struct {
		name   string
		fields map[string]string
		want   string
		want2  string
	}{
		{
			name:   "no fields",
			fields: nil,
			want:   "",
			want2:  "",
		},
		{
			name:   "1 field",
			fields: map[string]string{"name": "wrong name"},
			want:   "name: wrong name\n",
			want2:  "name: wrong name\n",
		},
		{
			name:   "2 fields",
			fields: map[string]string{"name": "wrong name", "age": "wrong age"},
			want:   "name: wrong name\nage: wrong age\n",
			want2:  "age: wrong age\nname: wrong name\n",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := &ErrValidation{
				fields: tt.fields,
			}
			if got := e.Error(); got != tt.want && got != tt.want2 {
				t.Errorf("Error() = %v, want %v, want2 %v", got, tt.want, tt.want2)
			}
		})
	}
}

func TestErrValidation_Fields(t *testing.T) {
	tests := []struct {
		name   string
		fields map[string]string
		want   map[string]string
	}{
		{
			name:   "no fields",
			fields: nil,
			want:   nil,
		},
		{
			name:   "2 fields",
			fields: map[string]string{"name": "wrong name", "age": "wrong age"},
			want:   map[string]string{"name": "wrong name", "age": "wrong age"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := &ErrValidation{
				fields: tt.fields,
			}
			if got := e.Fields(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Fields() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestErrValidation_IsEmpty(t *testing.T) {
	tests := []struct {
		name   string
		fields map[string]string
		want   bool
	}{
		{
			name:   "empty error with no fields",
			fields: map[string]string{},
			want:   true,
		},
		{
			name:   "2 fields",
			fields: map[string]string{"error": "error desc", "error2": "error desc"},
			want:   false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := &ErrValidation{
				fields: tt.fields,
			}
			if got := e.IsEmpty(); got != tt.want {
				t.Errorf("IsEmpty() = %v, want %v", got, tt.want)
			}
		})
	}
}
