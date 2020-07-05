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
	repoError := "repo error"
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
