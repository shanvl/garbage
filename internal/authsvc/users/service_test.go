package users

import (
	"context"
	"errors"
	"testing"

	"github.com/shanvl/garbage/internal/authsvc"
	"github.com/shanvl/garbage/internal/authsvc/mock"
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
	s := NewService(repo)
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
	s := NewService(repo)
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
