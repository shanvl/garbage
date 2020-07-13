// package jwt is responsible for generating and verifying jwt.
package jwt

import (
	"crypto/rsa"
	"errors"
	"fmt"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/shanvl/garbage/internal/authsvc"
)

// managerRSA is an implementation of Manager which uses RSA method to sign and verify jwt
type managerRSA struct {
	accessTokenDuration  time.Duration
	refreshTokenDuration time.Duration
	privateKey           *rsa.PrivateKey
	publicKey            *rsa.PublicKey
}

func NewManagerRSA(accessTokenDuration, refreshTokenDuration time.Duration, privateKey *rsa.PrivateKey,
	publicKey *rsa.PublicKey) authsvc.TokenManager {

	return &managerRSA{
		accessTokenDuration:  accessTokenDuration,
		refreshTokenDuration: refreshTokenDuration,
		privateKey:           privateKey,
		publicKey:            publicKey,
	}
}

// Generate generates jwt
func (m *managerRSA) Generate(tokenType authsvc.TokenType, clientID, userID string, role authsvc.Role) (string, error) {
	if clientID == "" {
		return "", errors.New("clientID must be provided")
	}
	if userID == "" {
		return "", errors.New("userID must be provided")
	}
	var expAt int64
	if tokenType == authsvc.Access {
		expAt = time.Now().Add(m.accessTokenDuration).Unix()
	}
	if tokenType == authsvc.Refresh {
		expAt = time.Now().Add(m.refreshTokenDuration).Unix()
	}
	claims := authsvc.UserClaims{
		StandardClaims: jwt.StandardClaims{
			Subject:   userID,
			ExpiresAt: expAt,
		},
		ClientID: clientID,
		Role:     role.String(),
		Type:     tokenType.String(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
	return token.SignedString(m.privateKey)
}

// Verify verifies jwt
func (m *managerRSA) Verify(token string) (authsvc.UserClaims, error) {
	if token == "" {
		return authsvc.UserClaims{}, errors.New("no token has been provided")
	}
	t, err := jwt.ParseWithClaims(token, &authsvc.UserClaims{}, func(token *jwt.Token) (interface{}, error) {
		_, ok := token.Method.(*jwt.SigningMethodRSA)
		if !ok {
			return authsvc.UserClaims{}, errors.New("unexpected signing algorithm")
		}
		return m.publicKey, nil
	})
	if err != nil {
		return authsvc.UserClaims{}, fmt.Errorf("invalid token: %w", err)
	}

	claims, ok := t.Claims.(*authsvc.UserClaims)
	if !ok {
		return authsvc.UserClaims{}, errors.New("invalid claims")
	}
	return *claims, nil
}
