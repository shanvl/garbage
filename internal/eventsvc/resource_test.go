package eventsvc

import (
	"reflect"
	"testing"
)

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
				ss: []string{"invalid resource"},
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "in upper case",
			args: args{
				ss: []string{"PLASTIC"},
			},
			want:    []Resource{Plastic},
			wantErr: false,
		},
		{
			name: "ok",
			args: args{
				[]string{"plastic", "gadgets"},
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

func TestResourceSliceToStringSlice(t *testing.T) {
	tests := []struct {
		name      string
		resources []Resource
		want      []string
	}{
		{
			name:      "no values",
			resources: []Resource{},
			want:      []string{},
		},
		{
			name:      "ok case",
			resources: []Resource{Plastic, Paper},
			want:      []string{"plastic", "paper"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ResourceSliceToStringSlice(tt.resources); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ResourceSliceToStringSlice() = %v, want %v", got, tt.want)
			}
		})
	}
}
