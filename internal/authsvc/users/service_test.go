package users_test

import (
	"context"
	"errors"
	"testing"

	"github.com/shanvl/garbage/internal/authsvc"
	"github.com/shanvl/garbage/internal/authsvc/mock"
	"github.com/shanvl/garbage/internal/authsvc/users"
)

func Test_service_CreateUser(t *testing.T) {
	ctx := context.Background()
	const repoError = "repo error"
	repo := &mock.UsersRepo{}
	repo.StoreUserFn = func(ctx context.Context, user *authsvc.User) error {
		if user.Email == repoError {
			return errors.New("error")
		}
		return nil
	}
	s := users.NewService(repo)
	type args struct {
		email string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "no email",
			args: args{
				email: "",
			},
			wantErr: true,
		},
		{
			name: "long email",
			args: args{
				email: "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa@aaa.com",
			},
			wantErr: true,
		},
		{
			name: "repo's error",
			args: args{
				email: repoError,
			},
			wantErr: true,
		},
		{
			name: "ok",
			args: args{
				email: "email",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			id, activationToken, err := s.CreateUser(ctx, tt.args.email)
			if (err != nil) != tt.wantErr {
				t.Errorf("CreateUser() error = %v, wantErr %v", err, tt.wantErr)
			}
			if err == nil && id == "" {
				t.Errorf("CreateUser() err == nil, len(id) == 0")
			}
			if err == nil && activationToken == "" {
				t.Errorf("CreateUser() err == nil, len(activationToken) == 0")
			}
		})
	}
}

func Test_service_ActivateUser(t *testing.T) {
	ctx := context.Background()
	const repoError = "repo error"
	const activationToken = "token"
	repo := &mock.UsersRepo{}
	repo.UserByActivationTokenFn = func(ctx context.Context, activationToken string) (*authsvc.User, error) {
		if activationToken == repoError {
			return nil, errors.New("error")
		}
		return &authsvc.User{
			ID:              "id",
			Active:          false,
			ActivationToken: activationToken,
			Email:           "email",
		}, nil
	}
	repo.StoreUserFn = func(ctx context.Context, user *authsvc.User) error {
		if user.Email == repoError {
			return errors.New("error")
		}
		return nil
	}
	s := users.NewService(repo)
	type args struct {
		activationToken string
		firstName       string
		lastName        string
		password        string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "no activation token",
			args: args{
				activationToken: "",
				firstName:       "fn",
				lastName:        "ln",
				password:        "psw",
			},
			wantErr: true,
		},
		{
			name: "no first name",
			args: args{
				activationToken: activationToken,
				firstName:       "",
				lastName:        "ln",
				password:        "psw",
			},
			wantErr: true,
		},
		{
			name: "too long first name",
			args: args{
				activationToken: activationToken,
				firstName:       "zzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzz",
				lastName:        "ln",
				password:        "psw",
			},
			wantErr: true,
		},
		{
			name: "no last name",
			args: args{
				activationToken: activationToken,
				firstName:       "fn",
				lastName:        "",
				password:        "psw",
			},
			wantErr: true,
		},
		{
			name: "too long last name",
			args: args{
				activationToken: activationToken,
				firstName:       "fn",
				lastName:        "zzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzz",
				password:        "psw",
			},
			wantErr: true,
		},
		{
			name: "no password",
			args: args{
				activationToken: activationToken,
				firstName:       "fn",
				lastName:        "ln",
				password:        "",
			},
			wantErr: true,
		},
		{
			name: "no password",
			args: args{
				activationToken: activationToken,
				firstName:       "fn",
				lastName:        "ln",
				password:        "zzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzz",
			},
			wantErr: true,
		},
		{
			name: "repo error",
			args: args{
				activationToken: repoError,
				firstName:       "fn",
				lastName:        "ln",
				password:        "pws",
			},
			wantErr: true,
		},
		{
			name: "ok",
			args: args{
				activationToken: activationToken,
				firstName:       "fn",
				lastName:        "ln",
				password:        "psw",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			id, err := s.ActivateUser(ctx, tt.args.activationToken, tt.args.firstName,
				tt.args.lastName, tt.args.password)
			if (err != nil) != tt.wantErr {
				t.Errorf("ActivateUser() error = %v, wantErr %v", err, tt.wantErr)
			}
			if err == nil && id == "" {
				t.Errorf("ActivateUser() error == nil, len(userID) == 0")
			}
		})
	}
}

func Test_service_ChangeUserRole(t *testing.T) {
	ctx := context.Background()
	const repoError = "repo error"
	repo := &mock.UsersRepo{}
	repo.ChangeUserRoleFn = func(ctx context.Context, id string, role authsvc.Role) error {
		if id == repoError {
			return errors.New("error")
		}
		return nil
	}
	s := users.NewService(repo)
	type args struct {
		id   string
		role authsvc.Role
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "no id",
			args: args{
				id:   "",
				role: authsvc.Member,
			},
			wantErr: true,
		},
		{
			name: "repo error",
			args: args{
				id:   repoError,
				role: authsvc.Member,
			},
			wantErr: true,
		},
		{
			name: "no id",
			args: args{
				id:   "id",
				role: authsvc.Member,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := s.ChangeUserRole(ctx, tt.args.id, tt.args.role)
			if (err != nil) != tt.wantErr {
				t.Errorf("ChangeUserRole() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_service_DeleteUser(t *testing.T) {
	ctx := context.Background()
	const repoError = "repo error"
	repo := &mock.UsersRepo{}
	repo.DeleteUserFn = func(ctx context.Context, id string) error {
		if id == repoError {
			return errors.New("error")
		}
		return nil
	}
	s := users.NewService(repo)
	type args struct {
		id string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "no id",
			args: args{
				id: "",
			},
			wantErr: true,
		},
		{
			name: "repo error",
			args: args{
				id: repoError,
			},
			wantErr: true,
		},
		{
			name: "ok",
			args: args{
				id: "id",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := s.DeleteUser(ctx, tt.args.id)
			if (err != nil) != tt.wantErr {
				t.Errorf("DeleteUser() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_service_UserByID(t *testing.T) {
	ctx := context.Background()
	const repoError = "repo error"
	repo := &mock.UsersRepo{}
	repo.UserByIDFn = func(ctx context.Context, id string) (*authsvc.User, error) {
		if id == repoError {
			return nil, errors.New("error")
		}
		return &authsvc.User{
			ID:     "id",
			Active: false,
			Email:  "email",
			Role:   authsvc.Member,
		}, nil
	}
	s := users.NewService(repo)
	type args struct {
		id string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "no id",
			args: args{
				id: "",
			},
			wantErr: true,
		},
		{
			name: "repo error",
			args: args{
				id: repoError,
			},
			wantErr: true,
		},
		{
			name: "ok",
			args: args{
				id: "id",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			user, err := s.UserByID(ctx, tt.args.id)
			if (err != nil) != tt.wantErr {
				t.Errorf("UserByID() error = %v, wantErr %v", err, tt.wantErr)
			}
			if err == nil && user == nil {
				t.Errorf("UserByID() error == nil, user == nil")
			}
		})
	}
}

func Test_service_Users(t *testing.T) {
	ctx := context.Background()
	const repoError = "repo error"
	repo := &mock.UsersRepo{}
	uu := []*authsvc.User{{}, {}, {}}
	repo.UsersFn = func(ctx context.Context, nameAndEmail string, sorting users.Sorting, amount, skip int) ([]*authsvc.User,
		int, error) {
		if nameAndEmail == repoError {
			return nil, 0, errors.New("error")
		}
		return uu, 3, nil
	}
	s := users.NewService(repo)
	type args struct {
		nameAndEmail string
		sorting      users.Sorting
		amount, skip int
	}
	tests := []struct {
		name      string
		args      args
		wantErr   bool
		wantTotal int
	}{
		{
			name: "negative amount",
			args: args{
				nameAndEmail: "name",
				sorting:      users.NameDes,
				amount:       -50,
				skip:         0,
			},
			wantErr:   false,
			wantTotal: len(uu),
		},
		{
			name: "negative skip",
			args: args{
				nameAndEmail: "name",
				sorting:      users.NameDes,
				amount:       10,
				skip:         -50,
			},
			wantErr:   false,
			wantTotal: len(uu),
		},
		{
			name: "unspecified sorting",
			args: args{
				nameAndEmail: "name",
				sorting:      users.Unspecified,
				amount:       10,
				skip:         0,
			},
			wantErr:   false,
			wantTotal: len(uu),
		},
		{
			name: "ok",
			args: args{
				nameAndEmail: "name",
				sorting:      users.NameDes,
				amount:       10,
				skip:         0,
			},
			wantErr:   false,
			wantTotal: len(uu),
		},
		{
			name: "repo error",
			args: args{
				nameAndEmail: repoError,
				amount:       10,
				skip:         0,
			},
			wantErr:   true,
			wantTotal: 0,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, total, err := s.Users(ctx, tt.args.nameAndEmail, tt.args.sorting, tt.args.amount, tt.args.skip)
			if (err != nil) != tt.wantErr {
				t.Errorf("Users() error == %v, wantErr == %v", err, tt.wantErr)
			}
			if total != tt.wantTotal {
				t.Errorf("Users() total == %d, wantTotal == %d", total, tt.wantTotal)
			}
		})
	}
}
