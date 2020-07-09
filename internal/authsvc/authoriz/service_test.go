package authoriz_test

import (
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
	tm.VerifyFn = func(token string) (*authsvc.UserClaims, error) {
		if token == invalidToken {
			return nil, errors.New("error")
		}
		return &authsvc.UserClaims{Role: token}, nil
	}
	s := authoriz.NewService(tm, protectedRPC)
	type args struct {
		accessToken string
		rpcName     string
	}
	tests := []struct {
		name       string
		args       args
		wantClaims *authsvc.UserClaims
		wantErr    bool
	}{
		{
			name: "invalid token",
			args: args{
				accessToken: invalidToken,
				rpcName:     protectedRPCName,
			},
			wantClaims: nil,
			wantErr:    true,
		},
		{
			name: "invalid role",
			args: args{
				accessToken: "member",
				rpcName:     protectedRPCName,
			},
			wantClaims: nil,
			wantErr:    true,
		},
		{
			name: "unprotected RPC",
			args: args{
				accessToken: "member",
				rpcName:     "unprotected rpc",
			},
			wantClaims: &authsvc.UserClaims{Role: "member"},
			wantErr:    false,
		},
		{
			name: "protected RPC admin",
			args: args{
				accessToken: "admin",
				rpcName:     protectedRPCName,
			},
			wantClaims: &authsvc.UserClaims{Role: "admin"},
			wantErr:    false,
		},
		{
			name: "protected RPC root",
			args: args{
				accessToken: "root",
				rpcName:     protectedRPCName,
			},
			wantClaims: &authsvc.UserClaims{Role: "root"},
			wantErr:    false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := s.Authorize(tt.args.accessToken, tt.args.rpcName)
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
