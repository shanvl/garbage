package postgres_test

import (
	"context"
	"reflect"
	"testing"

	"github.com/shanvl/garbage/internal/authsvc"
	"github.com/shanvl/garbage/internal/authsvc/postgres"
)

func TestRepository_ChangeUserRole(t *testing.T) {
	r := postgres.NewUsersRepo(db)
	ctx := context.Background()
	t.Run("member to admin", func(t *testing.T) {
		id := "someid"
		defer deleteUserByID(t, id)
		u := &authsvc.User{
			ID:   id,
			Role: authsvc.Member,
		}
		storeUser(t, u)
		if err := r.ChangeUserRole(ctx, u.ID, authsvc.Admin); err != nil {
			t.Errorf("ChangeUserRole() error == %v, wantErr == false", err)
		}
		u = userByID(t, u.ID)
		if u.Role != authsvc.Admin {
			t.Errorf("ChangeUserRole() user role == %s, want == admin", u.Role.String())
		}
	})
	t.Run("unknown user", func(t *testing.T) {
		id := "unknownuser"
		if err := r.ChangeUserRole(ctx, id, authsvc.Admin); err == nil {
			t.Errorf("ChangeUserRole() error == nil, wantErr == ErrUnknownUser")
		}
	})
}

func TestRepository_DeleteUser(t *testing.T) {
	r := postgres.NewUsersRepo(db)
	ctx := context.Background()
	u := &authsvc.User{
		ID: "id",
	}
	storeUser(t, u)
	defer deleteUserByID(t, u.ID)
	testCases := []struct {
		name    string
		id      string
		wantErr bool
	}{
		{
			name:    "known user",
			id:      u.ID,
			wantErr: false,
		},
		{
			name:    "unknown user",
			id:      "unknownid",
			wantErr: false,
		},
	}
	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			err := r.DeleteUser(ctx, tt.id)
			if (err != nil) != tt.wantErr {
				t.Errorf("DeleteUser() error == %v, wantErr == %v", err, tt.wantErr)
			}
		})
	}
}

func TestRepository_StoreUser(t *testing.T) {
	r := postgres.NewUsersRepo(db)
	ctx := context.Background()
	u := &authsvc.User{
		ID:              "someid",
		Active:          true,
		ActivationToken: "token",
		Email:           "email",
		FirstName:       "fn",
		LastName:        "ln",
		PasswordHash:    "ph",
		Role:            authsvc.Member,
	}
	defer deleteUserByID(t, u.ID)
	t.Run("new user", func(t *testing.T) {
		err := r.StoreUser(ctx, u)
		if err != nil {
			t.Errorf("StoreUser() error == %v, wantErr == false", err)
		}
		savedUser := userByID(t, u.ID)
		if !reflect.DeepEqual(savedUser, u) {
			t.Errorf("StoreUser() saved user == %+v, want == %+v", savedUser, u)
		}
	})
	t.Run("update user", func(t *testing.T) {
		u.FirstName = "changed name"
		err := r.StoreUser(ctx, u)
		if err != nil {
			t.Errorf("StoreUser() error == %v, wantErr == false", err)
		}
		savedUser := userByID(t, u.ID)
		if !reflect.DeepEqual(savedUser, u) {
			t.Errorf("StoreUser() saved user == %+v, want == %+v", savedUser, u)
		}
	})
}

const storeUserQ = `
	insert into users (id, active, activation_token, email, first_name, last_name, password_hash, role)
	values ($1, $2, $3, $4, $5, $6, $7, $8);
`

func storeUser(t *testing.T, u *authsvc.User) {
	t.Helper()
	_, err := db.Exec(context.Background(), storeUserQ, u.ID, u.Active, u.ActivationToken, u.Email, u.FirstName,
		u.LastName, u.PasswordHash, u.Role.String())
	if err != nil {
		t.Fatalf("test helper: couldn't store the user: %v", err)
	}
}

const userByIDQ = `
	select id, active, activation_token, email, first_name, last_name, password_hash, role
	from users
	where id = $1;
`

func userByID(t *testing.T, id string) *authsvc.User {
	t.Helper()
	u := &authsvc.User{}
	var roleStr string
	err := db.QueryRow(context.Background(), userByIDQ, id).Scan(&u.ID, &u.Active, &u.ActivationToken, &u.Email,
		&u.FirstName, &u.LastName, &u.PasswordHash, &roleStr)
	if err != nil {
		t.Fatalf("test helper: couldn't get a user: %v", err)
	}
	u.Role, err = authsvc.StringToRole(roleStr)
	if err != nil {
		t.Fatalf("test helper: couldn't get a user: %v", err)
	}
	return u
}

func deleteUserByID(t *testing.T, id string) {
	t.Helper()
	_, err := db.Exec(context.Background(), "delete from users where id = $1", id)
	if err != nil {
		t.Fatalf("test helper: couldn't delete a user: %v", err)
	}
}
