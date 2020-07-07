package authent

import (
	"context"
	"fmt"

	gonanoid "github.com/matoous/go-nanoid"
	"github.com/shanvl/garbage/internal/authsvc"
	"github.com/shanvl/garbage/pkg/valid"
)

type Repository interface {
	DeleteClient(ctx context.Context, clientID string) error
	DeleteUserClients(ctx context.Context, userID string) error
	StoreClient(ctx context.Context, clientID string, refreshToken string) error
	UserByEmail(ctx context.Context, email string) (*authsvc.User, error)
}

type Service interface {
	// Login generates, saves and returns auth credentials for the user if the given password and the email are correct
	Login(ctx context.Context, email, password string) (*authsvc.User, *AuthCreds, error)
	Logout(ctx context.Context, clientID string) error
	LogoutAllClients(ctx context.Context, userID string) error
	RefreshTokens(ctx context.Context, clientID, refreshToken string) (*AuthCreds, error)
}

type service struct {
	repo         Repository
	tokenManager authsvc.TokenManager
}

func NewService(repository Repository, tokenManager authsvc.TokenManager) Service {
	return &service{repository, tokenManager}
}

// Login generates, saves and returns auth credentials for the user if the given password and the email are correct
func (s *service) Login(ctx context.Context, email, password string) (*authsvc.User, *AuthCreds, error) {
	// validate the arguments
	errValid := valid.EmptyError()
	if email == "" {
		errValid.Add("email", "email is required")
	}
	if password == "" {
		errValid.Add("password", "password is required")
	}
	if !errValid.IsEmpty() {
		return nil, nil, errValid
	}
	// get the user by its email
	user, err := s.repo.UserByEmail(ctx, email)
	if err != nil {
		return nil, nil, err
	}
	// check if the user is in active state
	if user.Active == false {
		return nil, nil, authsvc.ErrInactiveUser
	}
	// check the password
	if !user.IsCorrectPassword(password) {
		return nil, nil, authsvc.ErrInvalidPassword
	}
	// generate auth credentials
	creds, err := s.generateAuthCreds(user.ID, user.Role)
	if err != nil {
		return nil, nil, err
	}
	// store the credentials
	err = s.repo.StoreClient(ctx, creds.ClientID, creds.RefreshToken)
	if err != nil {
		return nil, nil, err
	}
	return user, creds, nil
}

func (s *service) Logout(ctx context.Context, clientID string) error {
	panic("implement me")
}

func (s *service) LogoutAllClients(ctx context.Context, userID string) error {
	panic("implement me")
}

func (s *service) RefreshTokens(ctx context.Context, clientID, refreshToken string) (*AuthCreds, error) {
	// validate the arguments
	errValid := valid.EmptyError()
	if clientID == "" {
		errValid.Add("clientID", "clientID is required")
	}
	if refreshToken == "" {
		errValid.Add("refreshToken", "refreshToken is required")
	}
	if !errValid.IsEmpty() {
		return nil, errValid
	}
	panic("implement me")
}

// generateAuthCreds creates client id, access token and refresh token
func (s *service) generateAuthCreds(userID string, role authsvc.Role) (*AuthCreds, error) {
	clientID, err := gonanoid.Nanoid(15)
	if err != nil {
		return nil, fmt.Errorf("client id generation error: %w", err)
	}
	accessToken, err := s.tokenManager.Generate(authsvc.Access, clientID, userID, role)
	if err != nil {
		return nil, fmt.Errorf("access token generation error: %w", err)
	}
	refreshToken, err := s.tokenManager.Generate(authsvc.Refresh, clientID, userID, role)
	if err != nil {
		return nil, fmt.Errorf("refresh token generation error: %w", err)
	}
	return &AuthCreds{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ClientID:     clientID,
	}, nil
}

type AuthCreds struct {
	AccessToken  string
	RefreshToken string
	ClientID     string
}
