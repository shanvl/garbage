package postgres_test

import (
	"context"
	"reflect"
	"testing"

	"github.com/shanvl/garbage/internal/authsvc"
	"github.com/shanvl/garbage/internal/authsvc/postgres"
	usersSvc "github.com/shanvl/garbage/internal/authsvc/users"
)

func TestRepository_ChangeUserRole(t *testing.T) {
	r := postgres.NewUsersRepo(db)
	ctx := context.Background()
	t.Run("member to admin", func(t *testing.T) {
		u := &authsvc.User{
			ID:   "someid",
			Role: authsvc.Member,
		}
		storeUser(t, u)
		defer deleteUserByID(t, u.ID)
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
		ID: "someid",
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

func TestRepository_UserByActivationToken(t *testing.T) {
	r := postgres.NewUsersRepo(db)
	ctx := context.Background()
	u := &authsvc.User{
		ID:              "someid",
		ActivationToken: "sometoken",
	}
	storeUser(t, u)
	defer deleteUserByID(t, u.ID)
	testCases := []struct {
		name            string
		activationToken string
		wantErr         bool
		wantUser        *authsvc.User
	}{
		{
			name:            "existing user",
			activationToken: u.ActivationToken,
			wantErr:         false,
			wantUser:        u,
		},
		{
			name:            "unknown token",
			activationToken: "unknowntoken",
			wantErr:         true,
			wantUser:        nil,
		},
	}
	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			u, err := r.UserByActivationToken(ctx, tt.activationToken)
			if (err != nil) != tt.wantErr {
				t.Errorf("UserByActivationToken() error == %v, wantErr == %v", err, tt.wantErr)
			}
			if !reflect.DeepEqual(u, tt.wantUser) {
				t.Errorf("UserByActivationToken() user == %+v, want == %+v", u, tt.wantUser)
			}
		})
	}
}

func TestRepository_UserByID(t *testing.T) {
	r := postgres.NewUsersRepo(db)
	ctx := context.Background()
	u := &authsvc.User{
		ID: "someid",
	}
	storeUser(t, u)
	defer deleteUserByID(t, u.ID)
	testCases := []struct {
		name     string
		id       string
		wantErr  bool
		wantUser *authsvc.User
	}{
		{
			name:     "existing user",
			id:       u.ID,
			wantErr:  false,
			wantUser: u,
		},
		{
			name:     "unknown id",
			id:       "unknownid",
			wantErr:  true,
			wantUser: nil,
		},
	}
	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			u, err := r.UserByID(ctx, tt.id)
			if (err != nil) != tt.wantErr {
				t.Errorf("UserByID() error == %v, wantErr == %v", err, tt.wantErr)
			}
			if !reflect.DeepEqual(u, tt.wantUser) {
				t.Errorf("UserByID() user == %+v, want == %+v", u, tt.wantUser)
			}
		})
	}
}

func TestRepository_Users(t *testing.T) {
	r := postgres.NewUsersRepo(db)
	ctx := context.Background()
	uu := []*authsvc.User{
		{ID: "someid1", FirstName: "FN1", LastName: "LN1", Email: "EMAIL1"},
		{ID: "someid2", FirstName: "FN2", LastName: "LN2", Email: "EMAIL2"},
		{ID: "someid3", FirstName: "FN3", LastName: "LN3", Email: "EMAIL3"},
	}
	for _, user := range uu {
		u := user
		storeUser(t, u)
		defer deleteUserByID(t, u.ID)
	}
	testCases := []struct {
		name         string
		textSearch   string
		sorting      usersSvc.Sorting
		amount, skip int
		wantErr      bool
		wantUsers    []*authsvc.User
	}{
		{
			name:       "name asc sorting and filled text search",
			textSearch: "fn ln email",
			sorting:    usersSvc.NameAsc,
			amount:     10,
			skip:       0,
			wantErr:    false,
			wantUsers:  []*authsvc.User{uu[0], uu[1], uu[2]},
		},
		{
			name:       "name desc sorting and filled text search",
			textSearch: "fn ln email",
			sorting:    usersSvc.NameDes,
			amount:     10,
			skip:       0,
			wantErr:    false,
			wantUsers:  []*authsvc.User{uu[2], uu[1], uu[0]},
		},
		{
			name:       "email asc sorting and filled text search",
			textSearch: "fn ln email",
			sorting:    usersSvc.EmailAsc,
			amount:     10,
			skip:       0,
			wantErr:    false,
			wantUsers:  []*authsvc.User{uu[0], uu[1], uu[2]},
		},
		{
			name:       "email des sorting and filled text search",
			textSearch: "fn ln email",
			sorting:    usersSvc.EmailDes,
			amount:     10,
			skip:       0,
			wantErr:    false,
			wantUsers:  []*authsvc.User{uu[2], uu[1], uu[0]},
		},
		{
			name:       "skip 2",
			textSearch: "fn ln email",
			sorting:    usersSvc.EmailAsc,
			amount:     10,
			skip:       2,
			wantErr:    false,
			wantUsers:  []*authsvc.User{uu[2]},
		},
		{
			name:       "no text search",
			textSearch: "",
			sorting:    usersSvc.EmailAsc,
			amount:     200,
			skip:       150,
			wantErr:    false,
			// will be ignored for this test
			wantUsers: []*authsvc.User{},
		},
	}
	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			users, _, err := r.Users(ctx, tt.textSearch, tt.sorting, tt.amount, tt.skip)
			if (err != nil) != tt.wantErr {
				t.Errorf("Users() error == %v, want == %v", err, tt.wantErr)
			}
			if !reflect.DeepEqual(users, tt.wantUsers) && tt.name != "no text search" {
				t.Errorf("Users() users == %+v, want == %+v", users, tt.wantUsers)
			}
		})
	}
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
