package grpc_test

import (
	"context"
	"testing"

	authv1pb "github.com/shanvl/garbage/api/auth/v1/pb"
	"github.com/shanvl/garbage/internal/authsvc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func TestServer_CreateUser(t *testing.T) {
	ctx := context.Background()
	const (
		email = "email"
	)
	tests := []struct {
		name string
		req  *authv1pb.CreateUserRequest
		code codes.Code
	}{
		{
			name: "no email",
			req:  &authv1pb.CreateUserRequest{},
			code: codes.InvalidArgument,
		},
		{
			name: "ok",
			req:  &authv1pb.CreateUserRequest{Email: email},
			code: codes.OK,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			res, err := server.CreateUser(ctx, tt.req)
			if tt.code == codes.OK {
				if err != nil {
					t.Errorf("CreateUser() error == %v, wantErr == false", err)
				}
				if res == nil {
					t.Errorf("CreateUser() res == nil, want != nil")
				}
				if userByEmail(t, tt.req.GetEmail()) == nil {
					t.Errorf("CreateUser() couldn't find created user")
				}
				// uniqueness of the email
				res1, err := server.CreateUser(ctx, tt.req)
				if err == nil {
					t.Errorf("CreateUser() no unique email error")
					deleteUserByID(t, res1.GetId())
				}
				deleteUserByID(t, res.GetId())
			} else {
				if err == nil {
					t.Errorf("CreateUser() error == nil, wantErr == true")
				}
				if res != nil {
					t.Errorf("CreateUser() res == %v, want == nil", res)
				}
				st, ok := status.FromError(err)
				if ok != true {
					t.Errorf("CreateUser() couldn't get status from err %v", err)
				}
				if st.Code() != tt.code {
					t.Errorf("CreateUser() err codes mismatch: code == %v, want == %v", st.Code(), tt.code)
				}
			}
		})
	}
}

func userByEmail(t *testing.T, email string) *authsvc.User {
	u, err := authentRepo.UserByEmail(context.Background(), email)
	if err != nil {
		t.Fatalf("couldn't get a user by email: %v", err)
	}
	return u
}
