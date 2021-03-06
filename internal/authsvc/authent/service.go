// package authent is responsible for authentication
package authent

import (
	"context"
	"fmt"

	gonanoid "github.com/matoous/go-nanoid"
	"github.com/shanvl/garbage/internal/authsvc"
	"github.com/shanvl/garbage/pkg/valid"
)

type Repository interface {
	ClientByID(ctx context.Context, clientID string) (client Client, err error)
	DeleteClient(ctx context.Context, clientID string) error
	DeleteUserClients(ctx context.Context, userID string) error
	StoreClient(ctx context.Context, client Client) error
	UserByEmail(ctx context.Context, email string) (*authsvc.User, error)
}

// Service is responsible for authentication
type Service interface {
	// Login generates, saves and returns auth credentials for the user if the given password and the email are correct
	Login(ctx context.Context, email, password string) (*authsvc.User, AuthCreds, error)
	// Logout deletes the user's client and refresh token from the db, thus, logging the user out
	Logout(ctx context.Context, clientID string) error
	// LogoutAllClients deletes all the user's clients and refresh tokens, thus, logging the user out of every device
	LogoutAllClients(ctx context.Context, userID string) error
	// RefreshTokens verifies the given refresh token and then creates, saves and returns new auth credentials
	RefreshTokens(ctx context.Context, refreshToken string) (AuthCreds, error)
}

type service struct {
	repo         Repository
	tokenManager authsvc.TokenManager
}

func NewService(repository Repository, tokenManager authsvc.TokenManager) Service {
	return &service{repository, tokenManager}
}

// Login generates, saves and returns auth credentials for the user if the given password and the email are correct
func (s *service) Login(ctx context.Context, email, password string) (*authsvc.User, AuthCreds, error) {
	// validate the arguments
	errValid := valid.EmptyError()
	if email == "" {
		errValid.Add("email", "email is required")
	}
	if password == "" {
		errValid.Add("password", "password is required")
	}
	if !errValid.IsEmpty() {
		return nil, AuthCreds{}, errValid
	}
	// get the user by its email
	user, err := s.repo.UserByEmail(ctx, email)
	if err != nil {
		return nil, AuthCreds{}, err
	}
	// check if the user is in active state
	if user.Active == false {
		return nil, AuthCreds{}, authsvc.ErrInactiveUser
	}
	// check the password
	if !user.IsCorrectPassword(password) {
		return nil, AuthCreds{}, authsvc.ErrInvalidPassword
	}
	// generate auth credentials
	creds, err := s.generateAuthCreds(user.ID, user.Role)
	if err != nil {
		return nil, AuthCreds{}, err
	}
	// store the credentials
	c := Client{
		ID:           creds.ClientID,
		UserID:       user.ID,
		RefreshToken: creds.Refresh,
	}
	err = s.repo.StoreClient(ctx, c)
	if err != nil {
		return nil, AuthCreds{}, err
	}
	return user, creds, nil
}

// Logout deletes the user's client and refresh token from the db, thus, logging the user out
func (s *service) Logout(ctx context.Context, clientID string) error {
	if clientID == "" {
		return valid.NewError("clientID", "clientID is required")
	}
	return s.repo.DeleteClient(ctx, clientID)
}

// LogoutAllClients deletes all the user's clients and refresh tokens, thus, logging the user out of every device
func (s *service) LogoutAllClients(ctx context.Context, userID string) error {
	if userID == "" {
		return valid.NewError("userID", "userID is required")
	}
	return s.repo.DeleteUserClients(ctx, userID)
}

// RefreshTokens verifies the given refresh token and then creates, saves and returns new auth credentials
func (s *service) RefreshTokens(ctx context.Context, refreshToken string) (AuthCreds, error) {
	// validate the arguments
	if refreshToken == "" {
		return AuthCreds{}, valid.NewError("refreshToken", "refreshToken is required")
	}
	// verify the token and extract its claims
	claims, err := s.tokenManager.Verify(refreshToken)
	if err != nil {
		return AuthCreds{}, fmt.Errorf("%w: %v", authsvc.ErrInvalidRefreshToken, err)
	}
	// use the claims to get and compare clientID, userID and refreshToken saved in the db
	client, err := s.repo.ClientByID(ctx, claims.ClientID)
	if err != nil {
		return AuthCreds{}, err
	}
	if client.ID != claims.ClientID || client.RefreshToken != refreshToken || client.UserID != claims.Subject {
		return AuthCreds{}, authsvc.ErrInvalidRefreshToken
	}
	// convert string role from claims to authsvc.Role
	role, err := authsvc.StringToRole(claims.Role)
	if err != nil {
		return AuthCreds{}, err
	}
	// generate new tokens
	tokens, err := s.generateTokenPair(client.ID, client.UserID, role)
	if err != nil {
		return AuthCreds{}, err
	}
	// store a newly created refresh token along with the clientID and the userID
	c := Client{
		ID:           client.ID,
		UserID:       client.UserID,
		RefreshToken: tokens.Refresh,
	}
	err = s.repo.StoreClient(ctx, c)
	if err != nil {
		return AuthCreds{}, err
	}
	return AuthCreds{
		Tokens:   tokens,
		ClientID: c.ID,
	}, nil
}

// generateAuthCreds creates client id, access token and refresh token
func (s *service) generateAuthCreds(userID string, role authsvc.Role) (AuthCreds, error) {
	// create clientID
	clientID, err := gonanoid.Nanoid(15)
	if err != nil {
		return AuthCreds{}, fmt.Errorf("client id generation error: %w", err)
	}
	// create access and refresh tokens
	tokens, err := s.generateTokenPair(clientID, userID, role)
	if err != nil {
		return AuthCreds{}, err
	}
	return AuthCreds{
		Tokens:   tokens,
		ClientID: clientID,
	}, nil
}

// generateTokenPair generates access and refresh tokens
func (s *service) generateTokenPair(clientID, userID string, role authsvc.Role) (Tokens, error) {
	accessToken, err := s.tokenManager.Generate(authsvc.Access, clientID, userID, role)
	if err != nil {
		return Tokens{}, fmt.Errorf("access token generation error: %w", err)
	}
	refreshToken, err := s.tokenManager.Generate(authsvc.Refresh, clientID, userID, role)
	if err != nil {
		return Tokens{}, fmt.Errorf("refresh token generation error: %w", err)
	}
	return Tokens{
		Access:  accessToken,
		Refresh: refreshToken,
	}, nil
}

type AuthCreds struct {
	Tokens
	ClientID string
}

type Client struct {
	ID           string
	UserID       string
	RefreshToken string
}

type Tokens struct {
	Access, Refresh string
}
