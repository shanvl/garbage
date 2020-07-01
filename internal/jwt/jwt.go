// package jwt is responsible for generating and verifying jwt
package jwt

import (
	"github.com/dgrijalva/jwt-go"
)

// TokenType is a type of jwt. There are two jwt types: refresh token and access token.
// Refresh tokens have a longer living time and stored in db,
// while access tokens are not stored in db and have shorter living time
type TokenType string

const (
	Access  TokenType = "access"
	Refresh TokenType = "refresh"
)

// Claims are a jwt payload
type Claims struct {
	jwt.StandardClaims
	// ClientID is used to distinguish between different user's clients (browsers,
	// apps) in order to have an option to revoke the refresh token and thus sign the user out of that client
	ClientID string
	// Role is a user's role
	Role string
	Type TokenType
}
