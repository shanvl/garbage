package grpc_test

import (
	"context"
	"testing"

	authv1pb "github.com/shanvl/garbage/api/auth/v1/pb"
	"github.com/shanvl/garbage/internal/authsvc"
	"github.com/shanvl/garbage/internal/authsvc/authent"
	"github.com/shanvl/garbage/internal/authsvc/grpc"
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

func TestServer_Logout(t *testing.T) {
	u := newUser(t, "someid", "someemail", "psw", authsvc.Member)
	storeUser(t, u)
	defer deleteUserByID(t, u.ID)
	c := authent.Client{
		ID:           "clientid",
		UserID:       u.ID,
		RefreshToken: "token",
	}
	storeClient(t, c)
	defer deleteClientByID(t, c.ID)
	type args struct {
		ctx context.Context
		req *authv1pb.LogoutRequest
	}
	tests := []struct {
		name string
		args args
		code codes.Code
	}{
		{
			name: "empty context",
			args: args{
				ctx: context.Background(),
				req: &authv1pb.LogoutRequest{},
			},
			code: codes.Internal,
		},
		{
			name: "invalid context type",
			args: args{
				ctx: context.WithValue(context.Background(), grpc.AuthCtxKey, struct{}{}),
				req: &authv1pb.LogoutRequest{},
			},
			code: codes.Internal,
		},
		{
			name: "ok",
			args: args{
				ctx: context.WithValue(context.Background(), grpc.AuthCtxKey, &authsvc.UserClaims{ClientID: c.ID}),
				req: &authv1pb.LogoutRequest{},
			},
			code: codes.OK,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			res, err := server.Logout(tt.args.ctx, tt.args.req)
			if tt.code == codes.OK {
				if err != nil {
					t.Errorf("Logout() error == %v, wantErr == false", err)
				}
				if res == nil {
					t.Errorf("Logout() res == nil, want != nil")
				}
			} else {
				if err == nil {
					t.Errorf("Logout() error == nil, wantErr == true")
				}
				if res != nil {
					t.Errorf("Logout() res == %v, want == nil", res)
				}
				st, ok := status.FromError(err)
				if ok != true {
					t.Errorf("Logout() couldn't get status from err %v", err)
				}
				if st.Code() != tt.code {
					t.Errorf("Logout() err codes mismatch: code == %v, want == %v", st.Code(), tt.code)
				}
				if tt.code == codes.OK {
					c := getClient(t, c.ID)
					emptyC := authent.Client{}
					if c != emptyC {
						t.Errorf("Logout() code == OK, client wasn't deleted %+v", c)
					}
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

func storeClient(t *testing.T, client authent.Client) {
	t.Helper()
	err := authentRepo.StoreClient(context.Background(), client)
	if err != nil {
		t.Fatalf("couldn't store client: %v", err)
	}
}

func deleteClientByID(t *testing.T, id string) {
	t.Helper()
	err := authentRepo.DeleteClient(context.Background(), id)
	if err != nil {
		t.Fatalf("couldn't delete client: %v", err)
	}
}

func deleteUserByID(t *testing.T, id string) {
	t.Helper()
	err := usersRepo.DeleteUser(context.Background(), id)
	if err != nil {
		t.Fatalf("couldn't delete user: %v", err)
	}
}

func getClient(t *testing.T, id string) authent.Client {
	t.Helper()
	c, err := authentRepo.ClientByID(context.Background(), id)
	if err != nil {
		t.Fatalf("couldn't get a client: %v", err)
	}
	return c
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
