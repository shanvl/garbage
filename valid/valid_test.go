package valid

import (
	"testing"
)

type validateFunction = func() (isValid bool, errorKey string, errorDescription string)

func Test_Check(t *testing.T) {
	tests := []struct {
		name              string
		validateFunctions []validateFunction
		wantErr           bool
	}{
		{
			name:              "no validate functions",
			validateFunctions: []validateFunction{},
			wantErr:           false,
		},
		{
			name: "2 validate functions with no errors",
			validateFunctions: []validateFunction{func() (bool, string, string) {
				return true, "", ""
			}, func() (bool, string, string) {
				return true, "", ""
			}},
			wantErr: false,
		},
		{
			name: "2 validate functions with 2 errors",
			validateFunctions: []validateFunction{func() (bool, string, string) {
				return false, "name", "wrong name"
			}, func() (bool, string, string) {
				return false, "age ", "wrong age"
			}},
			wantErr: true,
		},
		{
			name: "2 validate functions with 1 error",
			validateFunctions: []validateFunction{func() (bool, string, string) {
				return false, "name", "wrong name"
			}, func() (bool, string, string) {
				return true, "", ""
			}},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := Check(tt.validateFunctions...); (err != nil) != tt.wantErr {
				t.Errorf("Check() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
