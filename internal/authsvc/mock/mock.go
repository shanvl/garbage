package mock

import (
	"context"

	"github.com/shanvl/garbage/internal/authsvc"
	"github.com/shanvl/garbage/internal/authsvc/users"
)

// UsersRepo mocks users service's repository
type UsersRepo struct {
	ChangeUserRoleFn      func(ctx context.Context, id string, role authsvc.Role) error
	ChangeUserRoleInvoked bool

	DeleteUserFn      func(ctx context.Context, id string) error
	DeleteUserInvoked bool

	StoreUserFn      func(ctx context.Context, user *authsvc.User) error
	StoreUserInvoked bool

	UserByIDFn      func(ctx context.Context, id string) (*authsvc.User, error)
	UserByIDInvoked bool

	UserByActivationTokenFn      func(ctx context.Context, activationToken string) (*authsvc.User, error)
	UserByActivationTokenInvoked bool

	UsersFn func(ctx context.Context, nameAndEmail string, sorting users.Sorting, amount,
		skip int) ([]*authsvc.User, int, error)
	UsersInvoked bool
}

func (u *UsersRepo) ChangeUserRole(ctx context.Context, id string, role authsvc.Role) error {
	u.ChangeUserRoleInvoked = true
	return u.ChangeUserRoleFn(ctx, id, role)
}

func (u *UsersRepo) DeleteUser(ctx context.Context, id string) error {
	u.DeleteUserInvoked = true
	return u.DeleteUserFn(ctx, id)
}

func (u *UsersRepo) StoreUser(ctx context.Context, user *authsvc.User) error {
	u.StoreUserInvoked = true
	return u.StoreUserFn(ctx, user)
}

func (u *UsersRepo) UserByActivationToken(ctx context.Context, activationToken string) (*authsvc.User, error) {
	u.UserByActivationTokenInvoked = true
	return u.UserByActivationTokenFn(ctx, activationToken)
}

func (u *UsersRepo) UserByID(ctx context.Context, id string) (*authsvc.User, error) {
	u.UserByIDInvoked = true
	return u.UserByIDFn(ctx, id)
}

func (u *UsersRepo) Users(ctx context.Context, nameAndEmail string, sorting users.Sorting, amount,
	skip int) ([]*authsvc.User, int, error) {
	u.UsersInvoked = true
	return u.UsersFn(ctx, nameAndEmail, sorting, amount, skip)
}

// AuthRepo mocks auth service's repository
type AuthRepo struct {
	DeleteClientFn      func(ctx context.Context, clientID string) error
	DeleteClientInvoked bool

	DeleteUserClientsFn      func(ctx context.Context, userID string) error
	DeleteUserClientsInvoked bool

	StoreClientFn      func(ctx context.Context, clientID string, refreshToken string) error
	StoreClientInvoked bool

	UserByIDFn      func(ctx context.Context, userID string) (*authsvc.User, error)
	UserByIDInvoked bool
}

func (a *AuthRepo) DeleteClient(ctx context.Context, clientID string) error {
	a.DeleteClientInvoked = true
	return a.DeleteClientFn(ctx, clientID)
}

func (a *AuthRepo) DeleteUserClients(ctx context.Context, userID string) error {
	a.DeleteUserClientsInvoked = true
	return a.DeleteUserClientsFn(ctx, userID)
}

func (a *AuthRepo) StoreClient(ctx context.Context, clientID string, refreshToken string) error {
	a.StoreClientInvoked = true
	return a.StoreClientFn(ctx, clientID, refreshToken)
}

func (a *AuthRepo) UserByID(ctx context.Context, userID string) (*authsvc.User, error) {
	a.UserByIDInvoked = true
	return a.UserByID(ctx, userID)
}

// TokenManager mocks authsvc.TokenManager
type TokenManager struct {
	GenerateFn      func(tokenType authsvc.TokenType, clientID, userID string, role authsvc.Role) (string, error)
	GenerateInvoked bool

	VerifyFn      func(token string) (*authsvc.UserClaims, error)
	VerifyInvoked bool
}

func (t *TokenManager) Generate(tokenType authsvc.TokenType, clientID, userID string, role authsvc.Role) (string, error) {
	t.GenerateInvoked = true
	return t.GenerateFn(tokenType, clientID, userID, role)
}

func (t *TokenManager) Verify(token string) (*authsvc.UserClaims, error) {
	t.VerifyInvoked = true
	return t.VerifyFn(token)
}
