package postgres

import (
	"context"

	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/shanvl/garbage/internal/authsvc"
	"github.com/shanvl/garbage/internal/authsvc/users"
)

type usersRepo struct {
	db *pgxpool.Pool
}

func NewUsersRepo(db *pgxpool.Pool) users.Repository {
	return &usersRepo{db}
}

const changeUserRoleQuery = `
	update users
	set role = $1
	where id = $2
	returning id;
`

// ChangeUserRole changes the user's role to the provided role
func (u *usersRepo) ChangeUserRole(ctx context.Context, id string, role authsvc.Role) error {
	var returnedID string
	err := u.db.QueryRow(ctx, changeUserRoleQuery, role.String(), id).Scan(&returnedID)
	if err != nil {
		return err
	}
	if returnedID == "" {
		return authsvc.ErrUnknownUser
	}
	return nil
}

const deleteUserQuery = `
	delete from users
	where id = $1;
`

// DeleteUser deletes the user
func (u *usersRepo) DeleteUser(ctx context.Context, id string) error {
	_, err := u.db.Exec(ctx, deleteUserQuery, id)
	return err
}

const storeUserQuery = `
	insert into users (id, active, activation_token, email, first_name, last_name, password_hash, role)
	values ($1, $2, $3, $4, $5, $6, $7, $8)
	on conflict (id) do update
		set (active, activation_token, email, first_name, last_name, password_hash, role)
				= ($2, $3, $4, $5, $6, $7, $8);
`

// StoreUser upserts the given to user to the db
func (u *usersRepo) StoreUser(ctx context.Context, user *authsvc.User) error {
	_, err := u.db.Exec(ctx, storeUserQuery, user.ID, user.Active, user.ActivationToken, user.Email, user.FirstName,
		user.LastName, user.PasswordHash, user.Role.String())
	return err
}

func (u *usersRepo) UserByActivationToken(ctx context.Context, activationToken string) (*authsvc.User, error) {
	panic("implement me")
}

func (u *usersRepo) UserByID(ctx context.Context, id string) (*authsvc.User, error) {
	panic("implement me")
}

func (u *usersRepo) Users(ctx context.Context, nameAndEmail string, sorting users.Sorting, amount,
	skip int) ([]*authsvc.User, int, error) {
	panic("implement me")
}
