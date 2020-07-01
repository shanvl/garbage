package jwt

import (
	"io/ioutil"
	"testing"
	"time"

	"github.com/dgrijalva/jwt-go"
)

func TestGenerator_Generate(t *testing.T) {
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
			name: "no user id",
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

	gen := newTestGenerator(t)
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := gen.Generate(tt.args.tokenType, tt.args.clientID, tt.args.userID, tt.args.role)
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

func newTestGenerator(t *testing.T) *Generator {
	b, err := ioutil.ReadFile("./keys_test/test.rsa")
	if err != nil {
		t.Fatalf("couldn't read key file: %v", err)
	}
	prKey, err := jwt.ParseRSAPrivateKeyFromPEM(b)
	if err != nil {
		t.Fatalf("couldn't get private key: %v", err)
	}

	return NewGenerator(30*time.Minute, 120*time.Hour, prKey)
}
