package authent

import (
	"context"

	"github.com/shanvl/garbage/internal/authsvc"
)

type Repository interface {
	DeleteClient(ctx context.Context, clientID string) error
	DeleteUserClients(ctx context.Context, userID string) error
	StoreClient(ctx context.Context, clientID string, refreshToken string) error
	UserByID(ctx context.Context, userID string) (*authsvc.User, error)
}

type Service interface {
	Login(ctx context.Context, userID string) (*authsvc.User, AuthCreds, error)
	Logout(ctx context.Context, clientID string) error
	LogoutAllClients(ctx context.Context, userID string) error
}

type AuthCreds struct {
	AccessToken  string
	RefreshToken string
	ClientID     string
}
