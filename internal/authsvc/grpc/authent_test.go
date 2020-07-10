package grpc_test

import (
	"context"
	"testing"

	authv1pb "github.com/shanvl/garbage/api/auth/v1/pb"
	"github.com/shanvl/garbage/internal/authsvc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func TestServer_Login(t *testing.T) {
	ctx := context.Background()
	const (
		email    = "email"
		password = "password"
	)
	u := newUser(t, "someid", email, password, authsvc.Member)
	storeUser(t, u)
	defer deleteUserByID(t, u.ID)
	tests := []struct {
		name string
		req  *authv1pb.LoginRequest
		code codes.Code
	}{
		{
			name: "no email",
			req:  &authv1pb.LoginRequest{Password: password},
			code: codes.InvalidArgument,
		},
		{
			name: "no password",
			req:  &authv1pb.LoginRequest{Email: email},
			code: codes.InvalidArgument,
		},
		{
			name: "invalid password",
			req:  &authv1pb.LoginRequest{Email: email, Password: "invalid"},
			code: codes.Unauthenticated,
		},
		{
			name: "ok",
			req:  &authv1pb.LoginRequest{Email: email, Password: password},
			code: codes.OK,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			res, err := server.Login(ctx, tt.req)
			if tt.code == codes.OK {
				if err != nil {
					t.Errorf("Login() error == %v, wantErr == false", err)
				}
				if res == nil {
					t.Errorf("Login() res == nil, want != nil")
				}
			} else {
				if err == nil {
					t.Errorf("Login() error == nil, wantErr == true")
				}
				if res != nil {
					t.Errorf("Login() res == %v, want == nil", res)
				}
				st, ok := status.FromError(err)
				if ok != true {
					t.Errorf("Login() couldn't get status from err %v", err)
				}
				if st.Code() != tt.code {
					t.Errorf("Login() err codes mismatch: code == %v, want == %v", st.Code(), tt.code)
				}
			}
		})
	}
}

func storeUser(t *testing.T, u *authsvc.User) {
	t.Helper()
	err := usersRepo.StoreUser(context.Background(), u)
	if err != nil {
		t.Fatalf("couldn't store user: %v", err)
	}
}

func deleteUserByID(t *testing.T, id string) {
	t.Helper()
	err := usersRepo.DeleteUser(context.Background(), id)
	if err != nil {
		t.Fatalf("couldn't delete user: %v", err)
	}
}

func newUser(t *testing.T, id, email, password string, role authsvc.Role) *authsvc.User {
	t.Helper()
	u := &authsvc.User{
		ID:        id,
		Active:    true,
		Email:     email,
		FirstName: "fn",
		LastName:  "ln",
		Role:      role,
	}
	err := u.ChangePassword(password)
	if err != nil {
		t.Fatalf("couldn't create a user: %v", err)
	}
	return u
}
