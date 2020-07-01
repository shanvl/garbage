package jwt

import (
	"io/ioutil"
	"testing"
	"time"

	"github.com/dgrijalva/jwt-go"
)

func TestManagerRSA_Generate(t *testing.T) {
	type args struct {
		tokenType TokenType
		clientID  string
		userID    string
		role      string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "no client id",
			args: args{
				tokenType: Refresh,
				clientID:  "",
				userID:    "userid",
				role:      "admin",
			},
			wantErr: true,
		},
		{
			name: "no user id",
			args: args{
				tokenType: Refresh,
				clientID:  "clientid",
				userID:    "",
				role:      "admin",
			},
			wantErr: true,
		},
		{
			name: "no role",
			args: args{
				tokenType: Refresh,
				clientID:  "clientid",
				userID:    "userid",
				role:      "",
			},
			wantErr: true,
		},
		{
			name: "refresh",
			args: args{
				tokenType: Refresh,
				clientID:  "clientid",
				userID:    "userid",
				role:      "admin",
			},
			wantErr: false,
		},
		{
			name: "access",
			args: args{
				tokenType: Access,
				clientID:  "clientid",
				userID:    "userid",
				role:      "admin",
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
	type args struct {
		tokenType TokenType
		clientID  string
		userID    string
		role      string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "access",
			args: args{
				tokenType: Access,
				clientID:  "clientID",
				userID:    "userID",
				role:      "admin",
			},
			wantErr: false,
		},
		{
			name: "refresh",
			args: args{
				tokenType: Refresh,
				clientID:  "clientID",
				userID:    "userID",
				role:      "admin",
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
				ClientID != tt.args.clientID || claims.Role != tt.args.role) {
				t.Errorf("Verify() no error, claims don't match")
			}
		})
	}
}

func newTestManagerRSA(t *testing.T) Manager {
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
