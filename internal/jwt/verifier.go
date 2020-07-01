package jwt

import (
	"crypto/rsa"
	"errors"
	"fmt"

	"github.com/dgrijalva/jwt-go"
)

// Verifier verifies jwt with provided RSA public key
type Verifier struct {
	publicKey *rsa.PublicKey
}

func NewVerifier(publicKey *rsa.PublicKey) *Verifier {
	return &Verifier{publicKey: publicKey}
}

// Verify verifies jwt with provided RSA public key
func (v *Verifier) Verify(token string) (*Claims, error) {
	t, err := jwt.ParseWithClaims(token, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		_, ok := token.Method.(*jwt.SigningMethodRSA)
		if !ok {
			return nil, errors.New("unexpected signing algorithm")
		}
		return v.publicKey, nil
	})
	if err != nil {
		return nil, fmt.Errorf("invalid token: %w", err)
	}

	claims, ok := t.Claims.(*Claims)
	if !ok {
		return nil, errors.New("invalid claims")
	}
	return claims, nil
}
