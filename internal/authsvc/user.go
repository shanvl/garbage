package authsvc

import (
	"errors"
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

var ErrInvalidActivationToken = errors.New("invalid activation token")
var ErrUnknownUser = errors.New("unknown user")
var ErrUnknownRole = errors.New("unknown role")

// User is a user of the app
type User struct {
	ID              string
	Active          bool
	ActivationToken string
	Email           string
	FirstName       string
	LastName        string
	PasswordHash    string
	Role            Role
}

// Role is a role of the user. Different roles grant different permissions to use services
type Role int

const (
	Admin Role = iota
	Member
	Root
)

var roleStringValues = []string{"admin", "member", "root"}

// String returns the string value of a role
func (r Role) String() string {
	if r < 0 || int(r) >= len(roleStringValues) {
		return "unknown"
	}
	return roleStringValues[r]
}

var stringToRoleMap = map[string]Role{
	"admin":  Admin,
	"member": Member,
	"root":   Root,
}

// StringToRole converts a string to a role
func StringToRole(s string) (Role, error) {
	role, ok := stringToRoleMap[s]
	if !ok {
		return 0, fmt.Errorf("%w: %s", ErrUnknownRole, s)
	}
	return role, nil
}

// NewUser creates a new user
func NewUser(activationToken, id, email string) (*User, error) {
	if activationToken == "" {
		return nil, fmt.Errorf("new user: %w", ErrInvalidActivationToken)
	}
	return &User{
		ID:              id,
		ActivationToken: activationToken,
		Email:           email,
		Role:            Member,
	}, nil
}

// Activate changes the user's active state to active if the provided activation token equals to the user's
// activation token
func (u *User) Activate(activationToken string) error {
	if activationToken == "" || activationToken != u.ActivationToken {
		return ErrInvalidActivationToken
	}
	u.Active = true
	u.ActivationToken = ""
	return nil
}

// ChangePassword changes the user's password
func (u *User) ChangePassword(password string) error {
	passwordHash, err := createPasswordHash(password)
	if err != nil {
		return fmt.Errorf("change user password: %w", err)
	}
	u.PasswordHash = passwordHash
	return nil
}

// ChangeRole changes the user's role
func (u *User) ChangeRole(role Role) {
	u.Role = role
}

// Deactivate changes the user's active state to false.
// It requires an activation token which can be used later to activate the user
func (u *User) Deactivate(activationToken string) error {
	if activationToken == "" {
		return fmt.Errorf("deactivate user: %w", ErrInvalidActivationToken)
	}
	u.Active = false
	u.ActivationToken = activationToken
	return nil
}

// IsCorrectPassword compares the proved password with user's password hash
func (u *User) IsCorrectPassword(password string) bool {
	return comparePasswordHash(password, u.PasswordHash)
}

func createPasswordHash(password string) (string, error) {
	passwordHash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", fmt.Errorf("create password hash: %w", err)
	}
	return string(passwordHash), nil
}

func comparePasswordHash(password string, passwordHash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(passwordHash), []byte(password))
	return err == nil
}
