package jwt

import (
	"crypto/rsa"
	"errors"
	"time"

	"github.com/dgrijalva/jwt-go"
)

// Generator generates jwt using RSA signing method
type Generator struct {
	accessTokenDuration  time.Duration
	refreshTokenDuration time.Duration
	privateKey           *rsa.PrivateKey
}

func NewGenerator(accessTokenDuration, refreshTokenDuration time.Duration, privateKey *rsa.PrivateKey) *Generator {
	return &Generator{
		accessTokenDuration:  accessTokenDuration,
		refreshTokenDuration: refreshTokenDuration,
		privateKey:           privateKey,
	}
}

// Generate generates jwt using RSA signing method
func (m *Generator) Generate(tokenType TokenType, clientID, userID, role string) (string, error) {
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
	claims := Claims{
		StandardClaims: jwt.StandardClaims{
			Subject:   userID,
			ExpiresAt: expAt,
		},
		ClientID: clientID,
		Role:     role,
		Type:     tokenType,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
	return token.SignedString(m.privateKey)
}
