package mock

import (
	"context"

	"github.com/shanvl/garbage/internal/authsvc"
)

// UsersRepo is a mock repository for users service
type UsersRepo struct {
	ChangeUserRoleFn      func(ctx context.Context, id, role authsvc.Role) error
	ChangeUserRoleInvoked bool

	DeleteUserFn      func(ctx context.Context, id string) error
	DeleteUserInvoked bool

	StoreUserFn      func(ctx context.Context, user *authsvc.User) error
	StoreUserInvoked bool

	UserByIDFn      func(ctx context.Context, id string) (*authsvc.User, error)
	UserByIDInvoked bool

	UserByActivationTokenFn      func(ctx context.Context, activationToken string) (*authsvc.User, error)
	UserByActivationTokenInvoked bool

	UsersFn      func(ctx context.Context, nameAndEmail string) ([]*authsvc.User, error)
	UsersInvoked bool
}

func (u *UsersRepo) ChangeUserRole(ctx context.Context, id, role authsvc.Role) error {
	u.ChangeUserRoleInvoked = true
	err := u.ChangeUserRoleFn(ctx, id, role)
	return err
}

func (u *UsersRepo) DeleteUser(ctx context.Context, id string) error {
	u.DeleteUserInvoked = true
	err := u.DeleteUserFn(ctx, id)
	return err
}

func (u *UsersRepo) StoreUser(ctx context.Context, user *authsvc.User) error {
	u.StoreUserInvoked = true
	err := u.StoreUserFn(ctx, user)
	return err
}

func (u *UsersRepo) UserByActivationToken(ctx context.Context, activationToken string) (*authsvc.User, error) {
	u.UserByActivationTokenInvoked = true
	user, err := u.UserByActivationTokenFn(ctx, activationToken)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (u *UsersRepo) UserByID(ctx context.Context, id string) (*authsvc.User, error) {
	u.UserByIDInvoked = true
	user, err := u.UserByIDFn(ctx, id)
	return user, err
}

func (u *UsersRepo) Users(ctx context.Context, nameAndEmail string) ([]*authsvc.User, error) {
	u.UsersInvoked = true
	users, err := u.UsersFn(ctx, nameAndEmail)
	return users, err
}
