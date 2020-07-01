package jwt

// TokenType is a type of jwt. There are two jwt types: refresh token and access token.
// Refresh tokens have a longer living time and stored in db,
// while access tokens are not stored in db and have a short living time
type TokenType int

const (
	Access TokenType = iota
	Refresh
)

var stringValues = [...]string{"access", "refresh"}

func (t TokenType) String() string {
	return stringValues[t]
}
