package grpc

import (
	"reflect"
	"testing"

	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func Test_errWithDetails(t *testing.T) {
	type args struct {
		code    codes.Code
		message string
		details map[string]string
	}
	type want struct {
		code      codes.Code
		message   string
		keyValues map[string]string
	}
	tests := []struct {
		name    string
		args    args
		wantErr want
	}{
		{
			name: "2 field map",
			args: args{
				code:    codes.InvalidArgument,
				message: "error message",
				details: map[string]string{"errField1": "errMessage1", "errField2": "errMessage2"},
			},
			wantErr: want{
				code:      codes.InvalidArgument,
				message:   "error message",
				keyValues: map[string]string{"errField1": "errMessage1", "errField2": "errMessage2"},
			},
		},
		{
			name: "no map",
			args: args{
				code:    codes.InvalidArgument,
				message: "error message",
				details: nil,
			},
			wantErr: want{
				code:      codes.InvalidArgument,
				message:   "error message",
				keyValues: map[string]string{},
			},
		},
		{
			name: "no message, no details",
			args: args{
				code:    codes.InvalidArgument,
				message: "",
				details: nil,
			},
			wantErr: want{
				code:      codes.InvalidArgument,
				message:   "",
				keyValues: map[string]string{},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := errWithDetails(tt.args.code, tt.args.message, tt.args.details)
			w := want{}
			st := status.Convert(err)
			w.code = st.Code()
			w.message = st.Message()
			w.keyValues = map[string]string{}
			for _, detail := range st.Details() {
				if d, ok := detail.(*errdetails.BadRequest); ok {
					for _, violation := range d.GetFieldViolations() {
						w.keyValues[violation.GetField()] = violation.GetDescription()
					}
				}
			}
			if !reflect.DeepEqual(w, tt.wantErr) {
				t.Errorf("errWithDetails error: got: %v, want: %v\n", w, tt.wantErr)
				return
			}
		})
	}
}
