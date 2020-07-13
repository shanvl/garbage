package authoriz_test

import (
	"context"
	"errors"
	"reflect"
	"testing"

	"github.com/shanvl/garbage/internal/authsvc"
	"github.com/shanvl/garbage/internal/authsvc/authoriz"
	"github.com/shanvl/garbage/internal/authsvc/mock"
)

func Test_service_Authorize(t *testing.T) {
	t.Parallel()
	const invalidToken = "tmerror"
	const protectedRPCName = "somename"
	protectedRPC := map[string][]authsvc.Role{
		protectedRPCName: {authsvc.Root, authsvc.Admin},
	}
	tm := &mock.TokenManager{}
	tm.VerifyFn = func(token string) (authsvc.UserClaims, error) {
		if token == invalidToken {
			return authsvc.UserClaims{}, errors.New("error")
		}
		return authsvc.UserClaims{Role: token}, nil
	}
	s := authoriz.NewService(tm, protectedRPC)
	type args struct {
		accessToken string
		method      string
	}
	tests := []struct {
		name       string
		args       args
		wantClaims authsvc.UserClaims
		wantErr    bool
	}{
		{
			name: "no token",
			args: args{
				accessToken: "",
				method:      protectedRPCName,
			},
			wantClaims: authsvc.UserClaims{},
			wantErr:    true,
		},
		{
			name: "no method",
			args: args{
				accessToken: "member",
				method:      "",
			},
			wantClaims: authsvc.UserClaims{},
			wantErr:    true,
		},
		{
			name: "invalid token",
			args: args{
				accessToken: invalidToken,
				method:      protectedRPCName,
			},
			wantClaims: authsvc.UserClaims{},
			wantErr:    true,
		},
		{
			name: "invalid role",
			args: args{
				accessToken: "member",
				method:      protectedRPCName,
			},
			wantClaims: authsvc.UserClaims{},
			wantErr:    true,
		},
		{
			name: "unprotected RPC",
			args: args{
				accessToken: "member",
				method:      "unprotected rpc",
			},
			wantClaims: authsvc.UserClaims{},
			wantErr:    false,
		},
		{
			name: "protected RPC admin",
			args: args{
				accessToken: "admin",
				method:      protectedRPCName,
			},
			wantClaims: authsvc.UserClaims{Role: "admin"},
			wantErr:    false,
		},
		{
			name: "protected RPC root",
			args: args{
				accessToken: "root",
				method:      protectedRPCName,
			},
			wantClaims: authsvc.UserClaims{Role: "root"},
			wantErr:    false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := s.Authorize(context.Background(), tt.args.accessToken, tt.args.method)
			if (err != nil) != tt.wantErr {
				t.Errorf("Authorize() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.wantClaims) {
				t.Errorf("Authorize() got = %v, want %v", got, tt.wantClaims)
			}
		})
	}
}
