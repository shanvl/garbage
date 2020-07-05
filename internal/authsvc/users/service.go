package users

import (
	"context"
	"fmt"

	gonanoid "github.com/matoous/go-nanoid"
	"github.com/shanvl/garbage/internal/authsvc"
	"github.com/shanvl/garbage/pkg/valid"
)

// Repository is a repo required by Service
type Repository interface {
	ChangeUserRole(ctx context.Context, id, role authsvc.Role) error
	DeleteUser(ctx context.Context, id string) error
	// Upsert
	StoreUser(ctx context.Context, user *authsvc.User) error
	UserByID(ctx context.Context, id string) (*authsvc.User, error)
	Users(ctx context.Context, nameAndEmail string) ([]*authsvc.User, error)
}

// Service manages users
type Service interface {
	ActivateUser(ctx context.Context, activateToken, firstName, lastName, password string) error
	ChangeUserRole(ctx context.Context, id string, role authsvc.Role) error
	// CreateUser creates and stores a user, which must then be activated with the returned activation token
	// Note, that the user's password is not needed here, it is required on the activation step
	CreateUser(ctx context.Context, email string) (id string, activationToken string, err error)
	DeleteUser(ctx context.Context, id string) error
	UserByID(ctx context.Context, id string) (*authsvc.User, error)
	Users(ctx context.Context, nameAndEmail string) ([]*authsvc.User, error)
}

type service struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &service{repo}
}

func (service) ActivateUser(ctx context.Context, activateToken, firstName, lastName, password string) error {
	panic("implement me")
}

func (service) ChangeUserRole(ctx context.Context, id string, role authsvc.Role) error {
	panic("implement me")
}

// CreateUser creates and stores a user, which must then be activated with the returned activation token
// Note, that the user's password is not needed here, it is required on the activation step
func (s *service) CreateUser(ctx context.Context, email string) (string, string, error) {
	if email == "" {
		return "", "", valid.NewError("email", "email is required")
	}
	// create activation token
	activationToken, err := gonanoid.Nanoid(14)
	if err != nil {
		return "", "", fmt.Errorf("create user: create activation token: %w", err)
	}
	// create id
	userID, err := gonanoid.Nanoid(14)
	if err != nil {
		return "", "", fmt.Errorf("create user: creat id: %w", err)
	}
	// create a new user
	user, err := authsvc.NewUser(activationToken, userID, email)
	if err != nil {
		return "", "", fmt.Errorf("create user: %w", err)
	}
	// store the user
	err = s.repo.StoreUser(ctx, user)
	if err != nil {
		return "", "", fmt.Errorf("create user: %w", err)
	}
	return userID, activationToken, nil
}

func (service) DeleteUser(ctx context.Context, id string) error {
	panic("implement me")
}

func (service) UserByID(ctx context.Context, id string) (*authsvc.User, error) {
	panic("implement me")
}

func (service) Users(ctx context.Context, nameAndEmail string) ([]*authsvc.User, error) {
	panic("implement me")
}
