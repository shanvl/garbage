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

func TestServer_DeleteUser(t *testing.T) {
	ctx := context.Background()
	u := &authsvc.User{ID: "someid"}
	storeUser(t, u)
	defer deleteUserByID(t, u.ID)
	tests := []struct {
		name string
		req  *authv1pb.DeleteUserRequest
		code codes.Code
	}{
		{
			name: "no id",
			req:  &authv1pb.DeleteUserRequest{Id: ""},
			code: codes.InvalidArgument,
		},
		{
			name: "no user with such id",
			req:  &authv1pb.DeleteUserRequest{Id: "somerandomid"},
			code: codes.OK,
		},
		{
			name: "valid id",
			req:  &authv1pb.DeleteUserRequest{Id: u.ID},
			code: codes.OK,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			res, err := server.DeleteUser(ctx, tt.req)
			if tt.code == codes.OK {
				if err != nil {
					t.Errorf("DeleteUser() error == %v, wantErr == false", err)
				}
				if res == nil {
					t.Errorf("DeleteUser() res == nil, want != nil")
				}
				user := userByID(t, tt.req.GetId())
				if user != nil {
					t.Errorf("DeleteUser() user wasn't deleted %v", user)
				}
			} else {
				if err == nil {
					t.Errorf("DeleteUser() error == nil, wantErr == true")
				}
				if res != nil {
					t.Errorf("DeleteUser() res == %v, want == nil", res)
				}
				st, ok := status.FromError(err)
				if ok != true {
					t.Errorf("DeleteUser() couldn't get status from err %v", err)
				}
				if st.Code() != tt.code {
					t.Errorf("DeleteUser() err codes mismatch: code == %v, want == %v", st.Code(), tt.code)
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

func userByID(t *testing.T, id string) *authsvc.User {
	u, err := usersRepo.UserByID(context.Background(), id)
	if err != nil {
		if err == authsvc.ErrUnknownUser {
			return nil
		}
		t.Fatalf("couldn't get a user by id: %v", err)
	}
	return u
}
