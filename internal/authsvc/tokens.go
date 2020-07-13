package authsvc

import (
	"errors"

	"github.com/dgrijalva/jwt-go"
)

var (
	ErrInvalidAccessToken  = errors.New("invalid access token")
	ErrInvalidRefreshToken = errors.New("invalid refresh token")
	ErrUnknownClient       = errors.New("unknown client")
)

// TokenManager is used for generation and verification of tokens
type TokenManager interface {
	Generate(tokenType TokenType, clientID, userID string, role Role) (string, error)
	Verify(token string) (UserClaims, error)
}

// TokenType is a type of token. There are two token types: refresh tokens and access tokens.
type TokenType int

const (
	// Access tokens have a very short life time and not stored in db
	Access TokenType = iota
	// Refresh tokens have a rather long life time and stored in db
	Refresh
)

var tokenStringValues = [...]string{"access", "refresh"}

func (t TokenType) String() string {
	return tokenStringValues[t]
}

// UserClaims are a token payload
type UserClaims struct {
	jwt.StandardClaims
	// ClientID is used to distinguish between different user's clients (browsers, apps etc)
	// in order to have an option to revoke the corresponding refresh token and thus sign the user out of that client
	ClientID string
	// Role is a string representation of the user's role
	Role string
	// Type is a string representation of the token's type
	Type string
}
