package grpc_test

import (
	"context"
	"strconv"
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

func TestServer_ActivateUser(t *testing.T) {
	ctx := context.Background()
	u := &authsvc.User{ID: "someid", ActivationToken: "sometoken", Active: false}
	storeUser(t, u)
	defer deleteUserByID(t, u.ID)
	tests := []struct {
		name string
		req  *authv1pb.ActivateUserRequest
		code codes.Code
	}{
		{
			name: "no token",
			req: &authv1pb.ActivateUserRequest{
				ActivationToken: "",
				FirstName:       "fn",
				LastName:        "ln",
				Password:        "psw",
			},
			code: codes.InvalidArgument,
		},
		{
			name: "no first name",
			req: &authv1pb.ActivateUserRequest{
				ActivationToken: u.ActivationToken,
				FirstName:       "",
				LastName:        "ln",
				Password:        "psw",
			},
			code: codes.InvalidArgument,
		},
		{
			name: "no last name",
			req: &authv1pb.ActivateUserRequest{
				ActivationToken: u.ActivationToken,
				FirstName:       "fn",
				LastName:        "",
				Password:        "password",
			},
			code: codes.InvalidArgument,
		},
		{
			name: "no password",
			req: &authv1pb.ActivateUserRequest{
				ActivationToken: u.ActivationToken,
				FirstName:       "fn",
				LastName:        "ln",
				Password:        "",
			},
			code: codes.InvalidArgument,
		},
		{
			name: "no invalid activation token",
			req: &authv1pb.ActivateUserRequest{
				ActivationToken: "someothertoken",
				FirstName:       "fn",
				LastName:        "ln",
				Password:        "password",
			},
			code: codes.NotFound,
		},
		{
			name: "ok",
			req: &authv1pb.ActivateUserRequest{
				ActivationToken: u.ActivationToken,
				FirstName:       "fn",
				LastName:        "ln",
				Password:        "password",
			},
			code: codes.OK,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			res, err := server.ActivateUser(ctx, tt.req)
			if tt.code == codes.OK {
				if err != nil {
					t.Errorf("ActivateUser() error == %v, wantErr == false", err)
				}
				if res == nil {
					t.Errorf("ActivateUser() res == nil, want != nil")
				}
				user := userByID(t, u.ID)
				if !user.Active {
					t.Errorf("ActivateUser() user wasn't activated")
				}
			} else {
				if err == nil {
					t.Errorf("ActivateUser() error == nil, wantErr == true")
				}
				if res != nil {
					t.Errorf("ActivateUser() res == %v, want == nil", res)
				}
				st, ok := status.FromError(err)
				if ok != true {
					t.Errorf("ActivateUser() couldn't get status from err %v", err)
				}
				if st.Code() != tt.code {
					t.Errorf("ActivateUser() err codes mismatch: code == %v, want == %v", st.Code(), tt.code)
				}
			}
		})
	}
}

func TestServer_ChangeUserRole(t *testing.T) {
	ctx := context.Background()
	u := &authsvc.User{ID: "someid", Role: authsvc.Member}
	storeUser(t, u)
	defer deleteUserByID(t, u.ID)
	tests := []struct {
		name string
		req  *authv1pb.ChangeUserRoleRequest
		code codes.Code
	}{
		{
			name: "no id",
			req: &authv1pb.ChangeUserRoleRequest{
				Id:   "",
				Role: authv1pb.Role_ROLE_ADMIN,
			},
			code: codes.InvalidArgument,
		},
		{
			name: "no such user",
			req: &authv1pb.ChangeUserRoleRequest{
				Id:   "somerandomid",
				Role: authv1pb.Role_ROLE_ADMIN,
			},
			code: codes.NotFound,
		},
		{
			name: "invalid role",
			req: &authv1pb.ChangeUserRoleRequest{
				Id:   u.ID,
				Role: authv1pb.Role_ROLE_UNKNOWN,
			},
			code: codes.InvalidArgument,
		},
		{
			name: "ok",
			req: &authv1pb.ChangeUserRoleRequest{
				Id:   u.ID,
				Role: authv1pb.Role_ROLE_ADMIN,
			},
			code: codes.OK,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			res, err := server.ChangeUserRole(ctx, tt.req)
			if tt.code == codes.OK {
				if err != nil {
					t.Errorf("ChangeUserRole() error == %v, wantErr == false", err)
				}
				if res == nil {
					t.Errorf("ChangeUserRole() res == nil, want != nil")
				}
				user := userByID(t, u.ID)
				if user.Role != authsvc.Admin {
					t.Errorf("ChangeUserRole() role wasn't changed: %v, want: admin", user.Role.String())
				}
			} else {
				if err == nil {
					t.Errorf("ChangeUserRole() error == nil, wantErr == true")
				}
				if res != nil {
					t.Errorf("ChangeUserRole() res == %v, want == nil", res)
				}
				st, ok := status.FromError(err)
				if ok != true {
					t.Errorf("ChangeUserRole() couldn't get status from err %v", err)
				}
				if st.Code() != tt.code {
					t.Errorf("ChangeUserRole() err codes mismatch: code == %v, want == %v", st.Code(), tt.code)
				}
			}
		})
	}
}

func TestServer_FindUser(t *testing.T) {
	ctx := context.Background()
	u := &authsvc.User{ID: "someid"}
	storeUser(t, u)
	defer deleteUserByID(t, u.ID)
	tests := []struct {
		name string
		req  *authv1pb.FindUserRequest
		code codes.Code
	}{
		{
			name: "no id",
			req: &authv1pb.FindUserRequest{
				Id: "",
			},
			code: codes.InvalidArgument,
		},
		{
			name: "no such user",
			req: &authv1pb.FindUserRequest{
				Id: "somerandomid",
			},
			code: codes.NotFound,
		},
		{
			name: "ok",
			req: &authv1pb.FindUserRequest{
				Id: u.ID,
			},
			code: codes.OK,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			res, err := server.FindUser(ctx, tt.req)
			if tt.code == codes.OK {
				if err != nil {
					t.Errorf("FindUser() error == %v, wantErr == false", err)
				}
				if res == nil {
					t.Errorf("FindUser() res == nil, want != nil")
				}
				if res.User.Id != u.ID {
					t.Errorf("FindUser() id mismatch, want: %s, got: %s", res.User.Id, u.ID)
				}
			} else {
				if err == nil {
					t.Errorf("FindUser() error == nil, wantErr == true")
				}
				if res != nil {
					t.Errorf("FindUser() res == %v, want == nil", res)
				}
				st, ok := status.FromError(err)
				if ok != true {
					t.Errorf("FindUser() couldn't get status from err %v", err)
				}
				if st.Code() != tt.code {
					t.Errorf("FindUser() err codes mismatch: code == %v, want == %v", st.Code(), tt.code)
				}
			}
		})
	}
}

func TestServer_FindUsers(t *testing.T) {
	ctx := context.Background()
	length := 3
	for i := 0; i < length; i++ {
		iStr := strconv.Itoa(i)
		u := &authsvc.User{ID: "someid" + iStr, FirstName: "fnnn" + iStr, LastName: "lnnn" + iStr, Email: "emaillll" + iStr}
		storeUser(t, u)
		defer deleteUserByID(t, u.ID)
	}
	tests := []struct {
		name string
		req  *authv1pb.FindUsersRequest
		code codes.Code
	}{
		{
			name: "no text search",
			req: &authv1pb.FindUsersRequest{
				NameAndEmail: "",
				Sorting:      authv1pb.UserSorting_USER_SORTING_NAME_DESC,
				Amount:       10,
				Skip:         0,
			},
			code: codes.OK,
		},
		{
			name: "with text search",
			req: &authv1pb.FindUsersRequest{
				NameAndEmail: "fnnn lnnn",
				Sorting:      authv1pb.UserSorting_USER_SORTING_NAME_DESC,
				Amount:       10,
				Skip:         0,
			},
			code: codes.OK,
		},
		{
			name: "with unknown sorting",
			req: &authv1pb.FindUsersRequest{
				NameAndEmail: "fn ln",
				Sorting:      authv1pb.UserSorting_USER_SORTING_UNKNOWN,
				Amount:       10,
				Skip:         0,
			},
			code: codes.OK,
		},
		{
			name: "skip overflow",
			req: &authv1pb.FindUsersRequest{
				NameAndEmail: "fn ln",
				Sorting:      authv1pb.UserSorting_USER_SORTING_NAME_DESC,
				Amount:       10,
				Skip:         999,
			},
			code: codes.OK,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			res, err := server.FindUsers(ctx, tt.req)
			if tt.code == codes.OK {
				if err != nil {
					t.Errorf("FindUsers() error == %v, wantErr == false", err)
				}
				if res == nil {
					t.Errorf("FindUsers() res == nil, want != nil")
				}
				if tt.name == "with text search" && res.Total != uint32(length) && len(res.Users) != length {
					t.Errorf("FindUsers() expected length == %d, got %d", length, len(res.Users))
				}
			} else {
				if err == nil {
					t.Errorf("FindUsers() error == nil, wantErr == true")
				}
				if res != nil {
					t.Errorf("FindUsers() res == %v, want == nil", res)
				}
				st, ok := status.FromError(err)
				if ok != true {
					t.Errorf("FindUsers() couldn't get status from err %v", err)
				}
				if st.Code() != tt.code {
					t.Errorf("FindUsers() err codes mismatch: code == %v, want == %v", st.Code(), tt.code)
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
