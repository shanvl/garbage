package postgres

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/jackc/pgconn"
	"github.com/jackc/pgtype"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/jmoiron/sqlx"
	"github.com/shanvl/garbage/internal/authsvc"
	usersSvc "github.com/shanvl/garbage/internal/authsvc/users"
	pgtextsearch "github.com/shanvl/garbage/pkg/pg-text-search"
)

type usersRepo struct {
	db *pgxpool.Pool
}

func NewUsersRepo(db *pgxpool.Pool) usersSvc.Repository {
	return &usersRepo{db}
}

var sortingToOrderMap = map[usersSvc.Sorting]string{
	usersSvc.NameAsc:  "last_name asc, first_name asc",
	usersSvc.NameDes:  "last_name desc, first_name desc",
	usersSvc.EmailAsc: "email asc",
	usersSvc.EmailDes: "email desc",
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
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		if pgErr.ConstraintName == "users_unique_lower_email_idx" {
			err = authsvc.ErrDuplicateEmail
		}
	}
	return err
}

const userByActivationTokenQuery = `
	select id, active, activation_token, email, first_name, last_name, password_hash, role
	from users
	where activation_token = $1;
`

// UserByActivationToken gets a user with provided activation token
func (u *usersRepo) UserByActivationToken(ctx context.Context, activationToken string) (*authsvc.User, error) {
	user := &authsvc.User{}
	var roleStr string
	err := u.db.QueryRow(ctx, userByActivationTokenQuery, activationToken).Scan(&user.ID, &user.Active, &user.ActivationToken,
		&user.Email, &user.FirstName, &user.LastName, &user.PasswordHash, &roleStr)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, authsvc.ErrUnknownUser
		}
		return nil, err
	}
	user.Role, err = authsvc.StringToRole(roleStr)
	if err != nil {
		return nil, err
	}
	return user, nil
}

const userByIDQuery = `
	select id, active, activation_token, email, first_name, last_name, password_hash, role
	from users
	where id = $1;
`

// UserByID gets a user with the given id
func (u *usersRepo) UserByID(ctx context.Context, id string) (*authsvc.User, error) {
	user := &authsvc.User{}
	var roleStr string
	err := u.db.QueryRow(ctx, userByIDQuery, id).Scan(&user.ID, &user.Active, &user.ActivationToken,
		&user.Email, &user.FirstName, &user.LastName, &user.PasswordHash, &roleStr)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, authsvc.ErrUnknownUser
		}
		return nil, err
	}
	user.Role, err = authsvc.StringToRole(roleStr)
	if err != nil {
		return nil, err
	}
	return user, nil
}

const usersQuery = `
	with query as (
    select id, active, activation_token, email, first_name, last_name, password_hash, role
    from users
	where 1=1 %s
	),  pagination as (
		select *
		from query
		order by %s
		limit ? offset ?
	)
	select *
	from pagination
			 right join (select count(*) FROM query) as c(total) on true;
`

// Users returns a list of sorted users
// "nameAndEmail" may consist of any combination of the email, first name and last name parts
func (u *usersRepo) Users(ctx context.Context, nameAndEmail string, sorting usersSvc.Sorting, amount,
	skip int) ([]*authsvc.User, int, error) {

	// get the "order by" part of the query
	orderBy := sortingToOrderMap[sorting]
	// create the "where" part of the query
	where := strings.Builder{}
	var args []interface{}
	if nameAndEmail != "" {
		textSearch := pgtextsearch.PrepareQuery(nameAndEmail)
		where.WriteString("and text_search @@ to_tsquery('simple', ?) ")
		args = append(args, textSearch)
	}
	// add limit and offset
	args = append(args, amount, skip)
	// embed the "where" and the "order by" parts to the query
	q := fmt.Sprintf(usersQuery, where.String(), orderBy)
	// change "?" to "$" in the query
	q = sqlx.Rebind(sqlx.BindType("pgx"), q)

	rows, err := u.db.Query(ctx, q, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()
	// "total" column will always be returned, so other columns might be null
	var (
		id       pgtype.Varchar
		active   pgtype.Bool
		actToken pgtype.Text
		email    pgtype.Varchar
		fN       pgtype.Varchar
		lN       pgtype.Varchar
		pasH     pgtype.Text
		role     pgtype.Text
		total    int
	)
	users := []*authsvc.User{}
	for rows.Next() {
		err := rows.Scan(&id, &active, &actToken, &email, &fN, &lN, &pasH, &role, &total)
		if err != nil {
			return nil, 0, err
		}
		// next will happen if the offset >= total rows found or no users matching the provided criteria have been
		// found. In that case we simply return total w/o additional work
		if id.Status != pgtype.Present {
			return nil, 0, nil
		}
		// create user
		user := &authsvc.User{
			ID:              id.String,
			Active:          active.Bool,
			ActivationToken: actToken.String,
			Email:           email.String,
			FirstName:       fN.String,
			LastName:        lN.String,
			PasswordHash:    pasH.String,
		}
		// convert role string to Role type
		r, err := authsvc.StringToRole(role.String)
		if err != nil {
			return nil, 0, err
		}
		user.Role = r
		// push the user to the users
		users = append(users, user)
	}
	if rows.Err() != nil {
		return nil, 0, nil
	}
	return users, total, nil
}
