package authent_test

import (
	"context"
	"errors"
	"testing"

	"github.com/dgrijalva/jwt-go"
	"github.com/shanvl/garbage/internal/authsvc"
	"github.com/shanvl/garbage/internal/authsvc/authent"
	"github.com/shanvl/garbage/internal/authsvc/mock"
)

func Test_service_Login(t *testing.T) {
	t.Parallel()
	const (
		repoGetError   = "repoerror"
		repoStoreError = "repostoreerror"
		inactiveUser   = "inactive"
		tmError        = "tmerror"
		validPassword  = "password"
	)
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
	r.StoreClientFn = func(ctx context.Context, client authent.Client) error {
		if client.RefreshToken == repoStoreError {
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
			user, creds, err := s.Login(ctx, tt.args.email, tt.args.password)
			if (err != nil) != tt.wantErr {
				t.Errorf("Login() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if err == nil && (user.ID == "" || user.Email == "" || user.FirstName == "" || user.LastName == "") {
				t.Errorf("Login() error == nil, invalid user == %+v", user)
			}
			if err == nil && (creds.ClientID == "" || creds.Access == "" || creds.Refresh == "") {
				t.Errorf("Login() error == nil, invalid creds == %+v", creds)
			}
		})
	}
}

func Test_service_Logout(t *testing.T) {
	t.Parallel()
	const (
		repoError = "repoerror"
	)
	ctx := context.Background()
	r := &mock.AuthRepo{}
	r.DeleteClientFn = func(ctx context.Context, clientID string) error {
		if clientID == repoError {
			return errors.New("error")
		}
		return nil
	}
	tm := &mock.TokenManager{}
	s := authent.NewService(r, tm)
	type args struct {
		clientID string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "no client id",
			args: args{
				clientID: "",
			},
			wantErr: true,
		},
		{
			name: "ok",
			args: args{
				clientID: "client id",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := s.Logout(ctx, tt.args.clientID)
			if (err != nil) != tt.wantErr {
				t.Errorf("Logout() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func Test_service_LogoutAllClients(t *testing.T) {
	t.Parallel()
	const (
		repoError = "repoerror"
	)
	ctx := context.Background()
	r := &mock.AuthRepo{}
	r.DeleteUserClientsFn = func(ctx context.Context, userID string) error {
		if userID == repoError {
			return errors.New("error")
		}
		return nil
	}
	tm := &mock.TokenManager{}
	s := authent.NewService(r, tm)
	type args struct {
		userID string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "no user id",
			args: args{
				userID: "",
			},
			wantErr: true,
		},
		{
			name: "ok",
			args: args{
				userID: "user id",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := s.LogoutAllClients(ctx, tt.args.userID)
			if (err != nil) != tt.wantErr {
				t.Errorf("LogoutAllClients() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func Test_service_RefreshTokens(t *testing.T) {
	t.Parallel()
	const (
		repoGetError      = "repogeterror"
		repoStoreError    = "repostoreerror"
		verifyError       = "verifyerror"
		generateError     = "generateerror"
		userID            = "userID"
		validRefreshToken = "token"
		clientID          = "clientID"
	)
	ctx := context.Background()
	r := &mock.AuthRepo{}
	r.ClientByIDFn = func(ctx context.Context, id string) (client authent.Client, err error) {
		if id == repoGetError {
			return authent.Client{}, errors.New("error")
		}
		return authent.Client{clientID, userID, validRefreshToken}, nil
	}
	r.StoreClientFn = func(ctx context.Context, client authent.Client) error {
		if client.ID == repoStoreError {
			return errors.New("error")
		}
		return nil
	}
	tm := &mock.TokenManager{}
	tm.GenerateFn = func(tokenType authsvc.TokenType, clientID, userID string, role authsvc.Role) (string, error) {
		if clientID == generateError {
			return "", errors.New("error")
		}
		return validRefreshToken, nil
	}
	tm.VerifyFn = func(token string) (authsvc.UserClaims, error) {
		if token == verifyError {
			return authsvc.UserClaims{}, errors.New("error")
		}
		return authsvc.UserClaims{ClientID: clientID, StandardClaims: jwt.StandardClaims{Subject: userID},
			Role: "member"}, nil
	}
	s := authent.NewService(r, tm)
	type args struct {
		refreshToken string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "no token",
			args: args{
				refreshToken: "",
			},
			wantErr: true,
		},
		{
			name: "verify token error",
			args: args{
				refreshToken: verifyError,
			},
			wantErr: true,
		},
		{
			name: "get client repo error",
			args: args{
				refreshToken: repoGetError,
			},
			wantErr: true,
		},
		{
			name: "generate tokens error",
			args: args{
				refreshToken: generateError,
			},
			wantErr: true,
		},
		{
			name: "store tokens error",
			args: args{
				refreshToken: repoStoreError,
			},
			wantErr: true,
		},
		{
			name: "ok",
			args: args{
				refreshToken: "token",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			creds, err := s.RefreshTokens(ctx, tt.args.refreshToken)
			if (err != nil) != tt.wantErr {
				t.Errorf("RefreshTokens() error = %v, wantErr %v", err, tt.wantErr)
			}
			if err == nil && (creds.ClientID == "" || creds.Refresh == "" || creds.Access == "") {
				t.Errorf("RefreshTokens() err == nil, invalid creds == %+v", creds)
			}
		})
	}
}
