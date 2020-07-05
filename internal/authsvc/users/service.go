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
	UserByActivationToken(ctx context.Context, activationToken string) (*authsvc.User, error)
	UserByID(ctx context.Context, id string) (*authsvc.User, error)
	Users(ctx context.Context, nameAndEmail string) ([]*authsvc.User, error)
}

// Service manages users
type Service interface {
	// ActivateUser changes the active state of the user to active and populates it with the provided additional info
	ActivateUser(ctx context.Context, activationToken, firstName, lastName, password string) (userID string, err error)
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

// ActivateUser changes the active state of the user to active and populates it with the provided additional info
func (s *service) ActivateUser(ctx context.Context, activationToken, firstName, lastName,
	password string) (userID string, err error) {
	// validate the arguments
	validErr := valid.EmptyError()
	if activationToken == "" {
		validErr.Add("activation token", "activation token is not provided")
	}
	if firstName == "" {
		validErr.Add("first name", "first name is not provided")
	}
	if lastName == "" {
		validErr.Add("last name", "last name is not provided")
	}
	if password == "" {
		validErr.Add("password", "password is not provided")
	}
	if !validErr.IsEmpty() {
		return "", validErr
	}

	// get the user
	user, err := s.repo.UserByActivationToken(ctx, activationToken)
	if err != nil {
		return "", fmt.Errorf("activate user: %w", err)
	}

	// activate the user
	err = user.Activate(activationToken)
	if err != nil {
		return "", fmt.Errorf("activate user: %w", err)
	}

	// set the user's password, first name and last name
	err = user.ChangePassword(password)
	if err != nil {
		return "", fmt.Errorf("activate user: %w", err)
	}
	user.FirstName = firstName
	user.LastName = lastName

	// store the user
	err = s.repo.StoreUser(ctx, user)
	if err != nil {
		return "", fmt.Errorf("activate user: %w", err)
	}
	return user.ID, nil
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