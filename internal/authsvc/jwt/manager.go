// package jwt is responsible for generating and verifying jwt.
package jwt

import (
	"crypto/rsa"
	"errors"
	"fmt"
	"time"

	"github.com/dgrijalva/jwt-go"
)

// Manager is used to generate and verify jwt
type Manager interface {
	Generate(tokenType TokenType, clientID, userID, role string) (string, error)
	Verify(token string) (*UserClaims, error)
}

// managerRSA is an implementation of Manager which uses RSA method to sign and verify jwt
type managerRSA struct {
	accessTokenDuration  time.Duration
	refreshTokenDuration time.Duration
	privateKey           *rsa.PrivateKey
	publicKey            *rsa.PublicKey
}

func NewManagerRSA(accessTokenDuration, refreshTokenDuration time.Duration, privateKey *rsa.PrivateKey,
	publicKey *rsa.PublicKey) Manager {

	return &managerRSA{
		accessTokenDuration:  accessTokenDuration,
		refreshTokenDuration: refreshTokenDuration,
		privateKey:           privateKey,
		publicKey:            publicKey,
	}
}

// Generate generates jwt
func (m *managerRSA) Generate(tokenType TokenType, clientID, userID, role string) (string, error) {
	if clientID == "" {
		return "", errors.New("clientID must be provided")
	}
	if userID == "" {
		return "", errors.New("userID must be provided")
	}
	if role == "" {
		return "", errors.New("role must be provided")
	}
	var expAt int64
	if tokenType == Access {
		expAt = time.Now().Add(m.accessTokenDuration).Unix()
	}
	if tokenType == Refresh {
		expAt = time.Now().Add(m.refreshTokenDuration).Unix()
	}
	claims := UserClaims{
		StandardClaims: jwt.StandardClaims{
			Subject:   userID,
			ExpiresAt: expAt,
		},
		ClientID: clientID,
		Role:     role,
		Type:     tokenType.String(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
	return token.SignedString(m.privateKey)
}

// Verify verifies jwt
func (m *managerRSA) Verify(token string) (*UserClaims, error) {
	t, err := jwt.ParseWithClaims(token, &UserClaims{}, func(token *jwt.Token) (interface{}, error) {
		_, ok := token.Method.(*jwt.SigningMethodRSA)
		if !ok {
			return nil, errors.New("unexpected signing algorithm")
		}
		return m.publicKey, nil
	})
	if err != nil {
		return nil, fmt.Errorf("invalid token: %w", err)
	}

	claims, ok := t.Claims.(*UserClaims)
	if !ok {
		return nil, errors.New("invalid claims")
	}
	return claims, nil
}

// UserClaims are a jwt payload
type UserClaims struct {
	jwt.StandardClaims
	// ClientID is used to distinguish between different user's clients (browsers,
	// apps) in order to have an option to revoke the refresh token and thus sign the user out of that client
	ClientID string
	// Role is a string representation of a user's role
	Role string
	// Type is a string representation of the token's type
	Type string
}
