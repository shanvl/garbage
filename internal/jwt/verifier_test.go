package jwt

import (
	"io/ioutil"
	"testing"

	"github.com/dgrijalva/jwt-go"
)

func TestGenerator_Verify(t *testing.T) {
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
	gen := newTestGenerator(t)
	ver := newTestVerifier(t)
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			token, err := gen.Generate(tt.args.tokenType, tt.args.clientID, tt.args.userID, tt.args.role)
			if err != nil {
				t.Fatalf("couldn't generate a token")
			}
			claims, err := ver.Verify(token)
			if (err != nil) != tt.wantErr {
				t.Errorf("Verify() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if err == nil && (claims.Subject != tt.args.userID || claims.Type != tt.args.tokenType || claims.
				ClientID != tt.args.clientID || claims.Role != tt.args.role) {
				t.Errorf("Verify() no error, claims don't match")
			}
		})
	}
}

func newTestVerifier(t *testing.T) *Verifier {
	b, err := ioutil.ReadFile("./keys_test/test.rsa.pub")
	if err != nil {
		t.Fatalf("couldn't read key file: %v", err)
	}
	pubKey, err := jwt.ParseRSAPublicKeyFromPEM(b)
	if err != nil {
		t.Fatalf("couldn't get public key: %v", err)
	}

	return NewVerifier(pubKey)
}
