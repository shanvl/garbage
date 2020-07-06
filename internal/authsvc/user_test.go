package authsvc

import (
	"reflect"
	"testing"
)

func TestNewUser(t *testing.T) {
	t.Parallel()
	type args struct {
		activationToken string
		id              string
		email           string
	}
	tests := []struct {
		name    string
		args    args
		want    *User
		wantErr bool
	}{
		{
			name: "no activation token",
			args: args{
				activationToken: "",
				id:              "id",
				email:           "email",
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "ok",
			args: args{
				activationToken: "token",
				id:              "id",
				email:           "email",
			},
			want: &User{
				ID:              "id",
				Active:          false,
				ActivationToken: "token",
				Email:           "email",
				FirstName:       "",
				LastName:        "",
				PasswordHash:    "",
				Role:            Member,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewUser(tt.args.activationToken, tt.args.id, tt.args.email)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewUser() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewUser() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestUser_Activate(t *testing.T) {
	t.Parallel()
	activationToken := "token"
	user := &User{
		ID:              "id",
		Active:          false,
		ActivationToken: activationToken,
		Email:           "email",
	}
	tests := []struct {
		name            string
		user            *User
		activationToken string
		wantUser        *User
		wantErr         bool
	}{
		{
			name:            "no token",
			user:            user,
			activationToken: "",
			wantErr:         true,
			wantUser:        user,
		},
		{
			name:            "ok",
			user:            user,
			activationToken: activationToken,
			wantErr:         false,
			wantUser: &User{
				ID:              user.ID,
				Active:          true,
				ActivationToken: "",
				Email:           user.Email,
				FirstName:       user.FirstName,
				LastName:        user.LastName,
				PasswordHash:    user.PasswordHash,
				Role:            user.Role,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.user.Activate(tt.activationToken)
			if (err != nil) != tt.wantErr {
				t.Errorf("Activate() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !reflect.DeepEqual(tt.user, tt.wantUser) {
				t.Errorf("Activate() user = %v, wantUser %v", tt.user, tt.wantUser)
			}
		})
	}
}

func TestUser_ChangePassword(t *testing.T) {
	t.Parallel()
	user := &User{}
	tests := []struct {
		name     string
		user     *User
		password string
		wantErr  bool
	}{
		{
			name:     "ok",
			user:     user,
			password: "password",
			wantErr:  false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.user.ChangePassword(tt.password)
			if (err != nil) != tt.wantErr {
				t.Errorf("ChangePassword() error = %v, wantErr %v", err, tt.wantErr)
			}
			if err == nil && tt.user.PasswordHash == "" {
				t.Errorf("ChangePassword() password hash wasn't created")
			}
		})
	}
}

func TestUser_Deactivate(t *testing.T) {
	t.Parallel()
	user := &User{
		Active: true,
	}
	tests := []struct {
		name            string
		user            *User
		activationToken string
		wantUser        *User
		wantErr         bool
	}{
		{
			name:            "no token",
			user:            user,
			activationToken: "",
			wantErr:         true,
			wantUser:        user,
		},
		{
			name:            "ok",
			user:            user,
			activationToken: "token",
			wantErr:         false,
			wantUser: &User{
				Active:          false,
				ActivationToken: "token",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.user.Deactivate(tt.activationToken)
			if (err != nil) != tt.wantErr {
				t.Errorf("Deactivate() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !reflect.DeepEqual(tt.user, tt.wantUser) {
				t.Errorf("Deactivate() user = %v, wantUser %v", tt.user, tt.wantUser)
			}
		})
	}
}

func TestUser_IsCorrectPassword(t *testing.T) {
	t.Parallel()
	password := "password"
	hash := testCreatePassword(t, password)
	user := &User{PasswordHash: hash}
	tests := []struct {
		name     string
		user     *User
		password string
		wantOk   bool
	}{
		{
			name:     "different password",
			user:     user,
			password: "123",
			wantOk:   false,
		},
		{
			name:     "ok",
			user:     user,
			password: password,
			wantOk:   true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ok := tt.user.IsCorrectPassword(tt.password)
			if ok != tt.wantOk {
				t.Errorf("IsCorrectPassword() ok == %v, wantOk == %v", ok, tt.wantOk)
			}
		})
	}
}

func testCreatePassword(t *testing.T, password string) string {
	hash, err := createPasswordHash(password)
	if err != nil {
		t.Fatal("error creating a password hash")
	}
	return hash
}
