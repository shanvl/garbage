package authent_test

import (
	"context"
	"errors"
	"testing"

	"github.com/shanvl/garbage/internal/authsvc"
	"github.com/shanvl/garbage/internal/authsvc/authent"
	"github.com/shanvl/garbage/internal/authsvc/mock"
)

func Test_service_Login(t *testing.T) {
	t.Parallel()
	const repoGetError = "repoerror"
	const repoStoreError = "repostoreerror"
	const inactiveUser = "inactive"
	const tmError = "tmerror"
	const validPassword = "password"
	ctx := context.Background()
	r := &mock.AuthRepo{}
	r.UserByEmailFn = func(ctx context.Context, email string) (*authsvc.User, error) {
		if email == repoGetError {
			return nil, errors.New("error")
		}
		u := &authsvc.User{
			ID:        "id",
			Active:    true,
			Email:     email,
			FirstName: "fn",
			LastName:  "ln",
			Role:      authsvc.Member,
		}
		if email == inactiveUser {
			u.Active = false
		}
		if email == tmError {
			u.ID = tmError
		}
		if email == repoStoreError {
			u.ID = repoStoreError
		}
		err := u.ChangePassword(validPassword)
		if err != nil {
			t.Fatalf("couldn't get mock user")
		}
		return u, nil
	}
	r.StoreClientFn = func(ctx context.Context, clientID string, refreshToken string) error {
		if refreshToken == repoStoreError {
			return errors.New("error")
		}
		return nil
	}
	tm := &mock.TokenManager{}
	tm.GenerateFn = func(tokenType authsvc.TokenType, clientID, userID string, role authsvc.Role) (string, error) {
		if userID == tmError {
			return "", errors.New("tm error")
		}
		if userID == repoStoreError {
			return repoStoreError, nil
		}
		return "token", nil
	}
	s := authent.NewService(r, tm)
	type args struct {
		email    string
		password string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "no email",
			args: args{
				email:    "",
				password: "psw",
			},
			wantErr: true,
		},
		{
			name: "no password",
			args: args{
				email:    "email",
				password: "",
			},
			wantErr: true,
		},
		{
			name: "inactive user",
			args: args{
				email:    inactiveUser,
				password: "psw",
			},
			wantErr: true,
		},
		{
			name: "invalid password",
			args: args{
				email:    "email",
				password: "psw",
			},
			wantErr: true,
		},
		{
			name: "repo error",
			args: args{
				email:    repoGetError,
				password: validPassword,
			},
			wantErr: true,
		},
		{
			name: "token manager error",
			args: args{
				email:    tmError,
				password: validPassword,
			},
			wantErr: true,
		},
		{
			name: "repo store error",
			args: args{
				email:    repoStoreError,
				password: validPassword,
			},
			wantErr: true,
		},
		{
			name: "ok",
			args: args{
				email:    "email",
				password: validPassword,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, _, err := s.Login(ctx, tt.args.email, tt.args.password)
			if (err != nil) != tt.wantErr {
				t.Errorf("Login() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}
