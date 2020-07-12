package grpc_test

import (
	"context"
	"testing"

	authv1pb "github.com/shanvl/garbage/api/auth/v1/pb"
	"github.com/shanvl/garbage/internal/authsvc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func TestServer_Authorize(t *testing.T) {
	ctx := context.Background()
	accessToken := generateAccessToken(t, "clientid", "userid", authsvc.Member)
	tests := []struct {
		name string
		req  *authv1pb.AuthorizeRequest
		code codes.Code
	}{
		{
			name: "no token",
			req: &authv1pb.AuthorizeRequest{
				Method: "/shanvl.garbage.auth.v1.AuthService/Logout",
				Token:  "",
			},
			code: codes.InvalidArgument,
		},
		{
			name: "no method",
			req: &authv1pb.AuthorizeRequest{
				Method: "",
				Token:  accessToken,
			},
			code: codes.InvalidArgument,
		},
		{
			name: "no permission",
			req: &authv1pb.AuthorizeRequest{
				Method: "/shanvl.garbage.auth.v1.AuthService/ActivateUser",
				Token:  accessToken,
			},
			code: codes.PermissionDenied,
		},
		{
			name: "protected, ok",
			req: &authv1pb.AuthorizeRequest{
				Method: "/shanvl.garbage.auth.v1.AuthService/Logout",
				Token:  accessToken,
			},
			code: codes.OK,
		},
		{
			name: "unprotected, ok",
			req: &authv1pb.AuthorizeRequest{
				Method: "/shanvl.garbage.auth.v1.AuthService/Login",
				Token:  accessToken,
			},
			code: codes.OK,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			res, err := server.Authorize(ctx, tt.req)
			if tt.code == codes.OK {
				if err != nil {
					t.Errorf("Authorize() error == %v, wantErr == false", err)
				}
				if res == nil {
					t.Errorf("Authorize() res == nil, want != nil")
				}
			} else {
				if err == nil {
					t.Errorf("Authorize() error == nil, wantErr == true")
				}
				if res != nil {
					t.Errorf("Authorize() res == %v, want == nil", res)
				}
				st, ok := status.FromError(err)
				if ok != true {
					t.Errorf("Authorize() couldn't get status from err %v", err)
				}
				if st.Code() != tt.code {
					t.Errorf("Authorize() err codes mismatch: code == %v, want == %v", st.Code(), tt.code)
				}
			}
		})
	}
}

func generateAccessToken(t *testing.T, clientID, userID string, role authsvc.Role) string {
	token, err := tokenManager.Generate(authsvc.Access, clientID, userID, role)
	if err != nil {
		t.Fatalf("couldn't generate access token: %v", err)
	}
	return token
}
