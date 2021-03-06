package jwt

import (
	"io/ioutil"
	"testing"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/shanvl/garbage/internal/authsvc"
)

func TestManagerRSA_Generate(t *testing.T) {
	t.Parallel()
	type args struct {
		tokenType authsvc.TokenType
		clientID  string
		userID    string
		role      authsvc.Role
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "no client id",
			args: args{
				tokenType: authsvc.Refresh,
				clientID:  "",
				userID:    "userid",
				role:      authsvc.Admin,
			},
			wantErr: true,
		},
		{
			name: "no user id",
			args: args{
				tokenType: authsvc.Refresh,
				clientID:  "clientid",
				userID:    "",
				role:      authsvc.Admin,
			},
			wantErr: true,
		},
		{
			name: "refresh",
			args: args{
				tokenType: authsvc.Refresh,
				clientID:  "clientid",
				userID:    "userid",
				role:      authsvc.Admin,
			},
			wantErr: false,
		},
		{
			name: "access",
			args: args{
				tokenType: authsvc.Access,
				clientID:  "clientid",
				userID:    "userid",
				role:      authsvc.Admin,
			},
			wantErr: false,
		},
	}

	manager := newTestManagerRSA(t)
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := manager.Generate(tt.args.tokenType, tt.args.clientID, tt.args.userID, tt.args.role)
			if (err != nil) != tt.wantErr {
				t.Errorf("Generate() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if err == nil && got == "" {
				t.Errorf("Generate() no error and no token")
			}
		})
	}
}

func TestManagerRSA_Verify(t *testing.T) {
	t.Parallel()
	type args struct {
		tokenType authsvc.TokenType
		clientID  string
		userID    string
		role      authsvc.Role
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "access",
			args: args{
				tokenType: authsvc.Access,
				clientID:  "clientID",
				userID:    "userID",
				role:      authsvc.Admin,
			},
			wantErr: false,
		},
		{
			name: "refresh",
			args: args{
				tokenType: authsvc.Refresh,
				clientID:  "clientID",
				userID:    "userID",
				role:      authsvc.Admin,
			},
			wantErr: false,
		},
	}
	manager := newTestManagerRSA(t)
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			token, err := manager.Generate(tt.args.tokenType, tt.args.clientID, tt.args.userID, tt.args.role)
			if err != nil {
				t.Fatalf("couldn't generate a token")
			}
			claims, err := manager.Verify(token)
			if (err != nil) != tt.wantErr {
				t.Errorf("Verify() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if err == nil && (claims.Subject != tt.args.userID || claims.Type != tt.args.tokenType.String() || claims.
				ClientID != tt.args.clientID || claims.Role != tt.args.role.String()) {
				t.Errorf("Verify() no error, claims don't match")
			}
		})
	}
}

func newTestManagerRSA(t *testing.T) authsvc.TokenManager {
	prB, err := ioutil.ReadFile("./keys_test/test.rsa")
	if err != nil {
		t.Fatalf("couldn't read key file: %v", err)
	}
	prKey, err := jwt.ParseRSAPrivateKeyFromPEM(prB)
	if err != nil {
		t.Fatalf("couldn't get private key: %v", err)
	}
	pubB, err := ioutil.ReadFile("./keys_test/test.rsa.pub")
	if err != nil {
		t.Fatalf("couldn't read key file: %v", err)
	}
	pubKey, err := jwt.ParseRSAPublicKeyFromPEM(pubB)
	if err != nil {
		t.Fatalf("couldn't get public key: %v", err)
	}

	return NewManagerRSA(30*time.Minute, 120*time.Hour, prKey, pubKey)
}
