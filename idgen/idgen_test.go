package idgen

import (
	"testing"
)

func Test_gen(t *testing.T) {
	tests := []struct {
		name       string
		length     int
		wantLength int
		wantErr    bool
	}{
		{
			name:       "len is 10",
			length:     10,
			wantLength: 10,
			wantErr:    false,
		},
		{
			name:       "len is 0",
			length:     0,
			wantLength: defLen,
			wantErr:    false,
		},
		{
			name:       "len is negative",
			length:     -10,
			wantLength: defLen,
			wantErr:    false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := gen(tt.length)
			gotLength := len(got)
			if (err != nil) != tt.wantErr {
				t.Errorf("gen() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotLength != tt.wantLength {
				t.Errorf("gen() gotLength = %v, wantLength %v", got, tt.wantLength)
			}
		})
	}
}
