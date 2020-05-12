package garbage

import (
	"reflect"
	"testing"
)

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

func TestStringSliceToResourceSlice(t *testing.T) {
	type args struct {
		ss []string
	}
	tests := []struct {
		name    string
		args    args
		want    []Resource
		wantErr bool
	}{
		{
			name: "invalid resource",
			args: args{
				ss: []string{string(Plastic), "invalid resource"},
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "ok",
			args: args{
				[]string{string(Plastic), string(Gadgets)},
			},
			want:    []Resource{Plastic, Gadgets},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := StringSliceToResourceSlice(tt.args.ss)
			if (err != nil) != tt.wantErr {
				t.Errorf("StringSliceToResourceSlice() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("StringSliceToResourceSlice() got = %v, want %v", got, tt.want)
			}
		})
	}
}
