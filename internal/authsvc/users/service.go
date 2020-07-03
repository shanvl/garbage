package users

import (
	"context"

	"github.com/shanvl/garbage/internal/authsvc"
)

type Repository interface {
	ChangeUserRole(ctx context.Context, id, role authsvc.Role) error
	DeleteUser(ctx context.Context, id string) error
	// Upsert
	StoreUser(ctx context.Context, user *authsvc.User) error
	UserByID(ctx context.Context, id string) (*authsvc.User, error)
	Users(ctx context.Context, nameAndEmail string) ([]*authsvc.User, error)
}

type Service interface {
	ActivateUser(ctx context.Context, activateToken, firstName, lastName, password string) error
	ChangeUserRole(ctx context.Context, id string, role authsvc.Role) error
	CreateUser(ctx context.Context, email string) (id string, activationToken string, err error)
	DeleteUser(ctx context.Context, id string) error
	UserByID(ctx context.Context, id string) (*authsvc.User, error)
	Users(ctx context.Context, nameAndEmail string) ([]*authsvc.User, error)
}
