package validation

import (
	"reflect"
	"testing"
)

func TestErrValidation_Error(t *testing.T) {
	tests := []struct {
		name   string
		fields map[string]string
		want   string
	}{
		{
			name:   "no fields",
			fields: nil,
			want:   "",
		},
		{
			name:   "1 field",
			fields: map[string]string{"name": "wrong name"},
			want:   "name: wrong name\n",
		},
		{
			name:   "2 fields",
			fields: map[string]string{"name": "wrong name", "age": "wrong age"},
			want:   "name: wrong name\nage: wrong age\n",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := &ErrValidation{
				fields: tt.fields,
			}
			if got := e.Error(); got != tt.want {
				t.Errorf("Error() = %v, want %v", got, tt.want)
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

// func TestErrValidation_add(t *testing.T) {
// 	type fields struct {
// 		fields map[string]string
// 	}
// 	type args struct {
// 		field string
// 		err   string
// 	}
// 	tests := []struct {
// 		name   string
// 		fields fields
// 		args   args
// 	}{
// 		// TODO: Add test cases.
// 	}
// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			e := &ErrValidation{
// 				fields: tt.fields.fields,
// 			}
// 		})
// 	}
// }
