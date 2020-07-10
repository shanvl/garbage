// package users is responsible for user management
package users

import (
	"context"

	gonanoid "github.com/matoous/go-nanoid"
	"github.com/shanvl/garbage/internal/authsvc"
	"github.com/shanvl/garbage/pkg/valid"
)

// Repository is a repo required by Service
type Repository interface {
	ChangeUserRole(ctx context.Context, id string, role authsvc.Role) error
	DeleteUser(ctx context.Context, id string) error
	StoreUser(ctx context.Context, user *authsvc.User) error
	UserByActivationToken(ctx context.Context, activationToken string) (*authsvc.User, error)
	UserByID(ctx context.Context, id string) (*authsvc.User, error)
	Users(ctx context.Context, nameAndEmail string, sorting Sorting, amount, skip int) ([]*authsvc.User, int, error)
}

// Service manages users
type Service interface {
	// ActivateUser changes the active state of the user to active and populates it with the provided additional info
	ActivateUser(ctx context.Context, activationToken, firstName, lastName, password string) (userID string, err error)
	// ChangeUserRole changes the user's role to the provided role
	ChangeUserRole(ctx context.Context, id string, role authsvc.Role) error
	// CreateUser creates and stores a user, which must then be activated with the returned activation token
	// Note, that the user's password is not needed here, it is required on the activation step
	CreateUser(ctx context.Context, email string) (id string, activationToken string, err error)
	// DeleteUser deletes the user
	DeleteUser(ctx context.Context, id string) error
	// UserByID returns the user with the specified id
	UserByID(ctx context.Context, id string) (*authsvc.User, error)
	// Users returns a sorted list of users
	// "nameAndEmail" may consist of any combination of the email, first name and last name parts
	Users(ctx context.Context, nameAndEmail string, sorting Sorting, amount, skip int) (users []*authsvc.User,
		total int, err error)
}

const (
	DefaultAmount = 50
	MaxAmount     = 1000
)

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
		validErr.Add("activation token", "activation token is required")
	}
	if firstName == "" {
		validErr.Add("firstName", "first name is required")
	}
	if len(firstName) > 35 {
		validErr.Add("lastName", "length of the first name can't be more than 35")
	}
	if lastName == "" {
		validErr.Add("lastName", "last name is required")
	}
	if len(lastName) > 35 {
		validErr.Add("lastName", "length of the last name can't be more than 35")
	}
	if password == "" {
		validErr.Add("password", "password is required")
	}
	if len(password) > 50 {
		validErr.Add("password", "length of the password can't be more than 50")
	}
	if !validErr.IsEmpty() {
		return "", validErr
	}

	// get the user
	user, err := s.repo.UserByActivationToken(ctx, activationToken)
	if err != nil {
		return "", err
	}

	// activate the user
	err = user.Activate(activationToken)
	if err != nil {
		return "", err
	}

	// set the user's password, first name and last name
	err = user.ChangePassword(password)
	if err != nil {
		return "", err
	}
	user.FirstName = firstName
	user.LastName = lastName

	// store the user
	err = s.repo.StoreUser(ctx, user)
	if err != nil {
		return "", err
	}
	return user.ID, nil
}

// ChangeUserRole changes the user's role to the provided one
func (s *service) ChangeUserRole(ctx context.Context, id string, role authsvc.Role) error {
	// validate the arguments
	if id == "" {
		return valid.NewError("id", "id is required")
	}

	return s.repo.ChangeUserRole(ctx, id, role)
}

// CreateUser creates and stores a user, which must then be activated with the returned activation token
// Note, that the user's password is not needed here, it is required on the activation step
func (s *service) CreateUser(ctx context.Context, email string) (string, string, error) {
	if email == "" {
		return "", "", valid.NewError("email", "email is required")
	}
	if len(email) > 50 {
		return "", "", valid.NewError("email", "length of the email can't be more than 50")
	}
	// create activation token
	activationToken, err := gonanoid.Nanoid(14)
	if err != nil {
		return "", "", err
	}
	// create id
	userID, err := gonanoid.Nanoid(14)
	if err != nil {
		return "", "", err
	}
	// create a new user
	user, err := authsvc.NewUser(activationToken, userID, email)
	if err != nil {
		return "", "", err
	}
	// store the user
	err = s.repo.StoreUser(ctx, user)
	if err != nil {
		return "", "", err
	}
	return userID, activationToken, nil
}

// DeleteUser deletes the user
func (s *service) DeleteUser(ctx context.Context, id string) error {
	if id == "" {
		return valid.NewError("id", "id is required")
	}
	return s.repo.DeleteUser(ctx, id)
}

// UserByID returns the user with the specified id
func (s *service) UserByID(ctx context.Context, id string) (*authsvc.User, error) {
	if id == "" {
		return nil, valid.NewError("id", "id is required")
	}
	return s.repo.UserByID(ctx, id)
}

// Users returns a sorted list of users
// "nameAndEmail" may consist of any combination of the email, first name and last name parts
func (s *service) Users(ctx context.Context, nameAndEmail string, sorting Sorting, amount, skip int) ([]*authsvc.User,
	int, error) {

	// validate sorting
	if sorting == Unspecified {
		sorting = NameDes
	}
	// validate amount and skip
	if amount <= 0 || amount > MaxAmount {
		amount = DefaultAmount
	}
	if skip < 0 {
		skip = 0
	}
	return s.repo.Users(ctx, nameAndEmail, sorting, amount, skip)
}

type User struct {
	ID        string
	Email     string
	FirstName string
	LastName  string
	Role      authsvc.Role
}
