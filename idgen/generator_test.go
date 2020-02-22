package idgen

import (
	"testing"
)

func Test_idGen_GenerateEventID(t *testing.T) {
	tests := []struct {
		name    string
		wantLen int
		wantErr bool
	}{
		{
			name:    "the output is of desirable length",
			wantLen: 14,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			i := NewIDGenerator()
			got, err := i.GenerateEventID()
			length := len(got)
			if (err != nil) != tt.wantErr {
				t.Errorf("GenerateEventID() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if length != tt.wantLen {
				t.Errorf("GenerateEventID() len(got) = %v, want %v", got, tt.wantLen)
			}
		})
	}
}
