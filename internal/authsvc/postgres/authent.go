package postgres

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/shanvl/garbage/internal/authsvc"
	"github.com/shanvl/garbage/internal/authsvc/authent"
)

type authentRepo struct {
	db *pgxpool.Pool
}

func NewAuthentRepo(db *pgxpool.Pool) authent.Repository {
	return &authentRepo{db}
}

const clientByIDQuery = `
	select id, user_id, refresh_token
	from clients
	where id = $1;
`

func (a *authentRepo) ClientByID(ctx context.Context, clientID string) (authent.Client, error) {
	c := authent.Client{}
	err := a.db.QueryRow(ctx, clientByIDQuery, clientID).Scan(&c.ID, &c.UserID, &c.RefreshToken)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return authent.Client{}, authsvc.ErrUnknownClient
		}
	}
	return c, nil
}

const deleteClientQuery = `
	delete from clients
	where id = $1;
`

func (a *authentRepo) DeleteClient(ctx context.Context, clientID string) error {
	_, err := a.db.Exec(ctx, deleteClientQuery, clientID)
	return err
}

const deleteUserClientsQuery = `
	delete from clients
	where user_id = $1;
`

func (a *authentRepo) DeleteUserClients(ctx context.Context, userID string) error {
	_, err := a.db.Exec(ctx, deleteUserClientsQuery, userID)
	return err
}

const storeClientQuery = `
	insert into clients (id, user_id, refresh_token)
	values ($1, $2, $3)
	on conflict (id) do update
		set (user_id, refresh_token) = ($2, $3);
`

func (a *authentRepo) StoreClient(ctx context.Context, client authent.Client) error {
	_, err := a.db.Exec(ctx, storeClientQuery, client.ID, client.UserID, client.RefreshToken)
	return err
}

const userByEmailQuery = `
	select id, active, activation_token, email, first_name, last_name, password_hash, role
	from users
	where lower(email) = $1;
`

func (a *authentRepo) UserByEmail(ctx context.Context, email string) (*authsvc.User, error) {
	u := &authsvc.User{}
	var roleStr string
	err := a.db.QueryRow(ctx, userByEmailQuery, email).Scan(&u.ID, &u.Active, &u.ActivationToken, &u.Email,
		&u.FirstName, &u.LastName, &u.PasswordHash, &roleStr)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, authsvc.ErrUnknownUser
		}
		return nil, err
	}
	u.Role, err = authsvc.StringToRole(roleStr)
	if err != nil {
		return nil, err
	}
	return u, nil
}
