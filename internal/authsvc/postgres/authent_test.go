package postgres_test

import (
	"context"
	"reflect"
	"testing"

	"github.com/shanvl/garbage/internal/authsvc"
	"github.com/shanvl/garbage/internal/authsvc/authent"
	"github.com/shanvl/garbage/internal/authsvc/postgres"
)

func TestRepository_ClientByID(t *testing.T) {
	r := postgres.NewAuthentRepo(db)
	ctx := context.Background()
	u := &authsvc.User{ID: "someid"}
	storeUser(t, u)
	defer deleteUserByID(t, u.ID)
	c := authent.Client{ID: "someid", UserID: u.ID}
	storeClient(t, c)
	defer deleteClientByID(t, c.ID)
	testCases := []struct {
		name       string
		id         string
		wantErr    bool
		wantClient authent.Client
	}{
		{
			name:       "existing client",
			id:         c.ID,
			wantErr:    false,
			wantClient: c,
		},
		{
			name:       "unknown id",
			id:         "unknownid",
			wantErr:    true,
			wantClient: authent.Client{},
		},
	}
	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			u, err := r.ClientByID(ctx, tt.id)
			if (err != nil) != tt.wantErr {
				t.Errorf("ClientByID() error == %v, wantErr == %v", err, tt.wantErr)
			}
			if !reflect.DeepEqual(u, tt.wantClient) {
				t.Errorf("ClientByID() user == %+v, want == %+v", u, tt.wantClient)
			}
		})
	}
}

func TestRepository_DeleteClient(t *testing.T) {
	r := postgres.NewAuthentRepo(db)
	ctx := context.Background()
	u := &authsvc.User{
		ID: "someid",
	}
	storeUser(t, u)
	defer deleteUserByID(t, u.ID)
	c := authent.Client{ID: "someid", UserID: u.ID}
	storeClient(t, c)
	defer deleteClientByID(t, c.ID)
	testCases := []struct {
		name    string
		id      string
		wantErr bool
	}{
		{
			name:    "known client",
			id:      c.ID,
			wantErr: false,
		},
		{
			name:    "unknown client",
			id:      "unknownid",
			wantErr: false,
		},
	}
	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			err := r.DeleteClient(ctx, tt.id)
			if (err != nil) != tt.wantErr {
				t.Errorf("DeleteClient() error == %v, wantErr == %v", err, tt.wantErr)
			}
			if _, err = clientByIDWithErr(t, c.ID); err == nil {
				t.Errorf("DeleteClient() client wasn't deleted")
			}
		})
	}
}

func TestRepository_DeleteUserClients(t *testing.T) {
	r := postgres.NewAuthentRepo(db)
	ctx := context.Background()
	u := &authsvc.User{
		ID: "someid",
	}
	storeUser(t, u)
	defer deleteUserByID(t, u.ID)
	c1 := authent.Client{ID: "someid", UserID: u.ID}
	storeClient(t, c1)
	defer deleteClientByID(t, c1.ID)
	c2 := authent.Client{ID: "someid1", UserID: u.ID}
	storeClient(t, c2)
	defer deleteClientByID(t, c2.ID)
	testCases := []struct {
		name    string
		id      string
		wantErr bool
	}{
		{
			name:    "user with clients",
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
			err := r.DeleteUserClients(ctx, tt.id)
			if (err != nil) != tt.wantErr {
				t.Errorf("DeleteUserClients() error == %v, wantErr == %v", err, tt.wantErr)
			}
			if _, err = clientByIDWithErr(t, c1.ID); err == nil {
				t.Errorf("DeleteUserClients() client wasn't deleted")
			}
			if _, err = clientByIDWithErr(t, c2.ID); err == nil {
				t.Errorf("DeleteUserClients() client wasn't deleted")
			}
		})
	}
}

func TestRepository_StoreClient(t *testing.T) {
	r := postgres.NewAuthentRepo(db)
	ctx := context.Background()
	u := &authsvc.User{ID: "userid"}
	storeUser(t, u)
	defer deleteUserByID(t, u.ID)
	c := authent.Client{
		ID:           "someid",
		UserID:       u.ID,
		RefreshToken: "token",
	}
	defer deleteClientByID(t, c.ID)
	t.Run("new client", func(t *testing.T) {
		err := r.StoreClient(ctx, c)
		if err != nil {
			t.Errorf("StoreClient() error == %v, wantErr == false", err)
		}
		savedClient := clientByID(t, c.ID)
		if !reflect.DeepEqual(savedClient, c) {
			t.Errorf("StoreClient() saved user == %+v, want == %+v", savedClient, c)
		}
	})
	t.Run("update client", func(t *testing.T) {
		c.RefreshToken = "new token"
		err := r.StoreClient(ctx, c)
		if err != nil {
			t.Errorf("StoreClient() error == %v, wantErr == false", err)
		}
		savedClient := clientByID(t, c.ID)
		if !reflect.DeepEqual(savedClient, c) {
			t.Errorf("StoreClient() saved user == %+v, want == %+v", savedClient, c)
		}
	})
}

func TestRepository_UserByEmail(t *testing.T) {
	r := postgres.NewAuthentRepo(db)
	ctx := context.Background()
	u := &authsvc.User{
		ID:    "someid",
		Email: "someemail",
	}
	storeUser(t, u)
	defer deleteUserByID(t, u.ID)
	testCases := []struct {
		name     string
		email    string
		wantErr  bool
		wantUser *authsvc.User
	}{
		{
			name:     "existing user",
			email:    u.Email,
			wantErr:  false,
			wantUser: u,
		},
		{
			name:     "unknown email",
			email:    "unknownid",
			wantErr:  true,
			wantUser: nil,
		},
	}
	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			u, err := r.UserByEmail(ctx, tt.email)
			if (err != nil) != tt.wantErr {
				t.Errorf("UserByEmail() error == %v, wantErr == %v", err, tt.wantErr)
			}
			if !reflect.DeepEqual(u, tt.wantUser) {
				t.Errorf("UserByEmail() user == %+v, want == %+v", u, tt.wantUser)
			}
		})
	}
}

const storeClientQ = `
	insert into clients (id, user_id, refresh_token)
	values ($1, $2, $3);
`

func storeClient(t *testing.T, c authent.Client) {
	t.Helper()
	_, err := db.Exec(context.Background(), storeClientQ, c.ID, c.UserID, c.RefreshToken)
	if err != nil {
		t.Fatalf("test helper: couldn't store the client: %v", err)
	}
}

func deleteClientByID(t *testing.T, id string) {
	t.Helper()
	_, err := db.Exec(context.Background(), "delete from clients where id = $1", id)
	if err != nil {
		t.Fatalf("test helper: couldn't delete a client: %v", err)
	}
}

const clientByIDQ = `
	select id, user_id, refresh_token
	from clients
	where id = $1;
`

func clientByID(t *testing.T, id string) authent.Client {
	t.Helper()
	c := authent.Client{}
	err := db.QueryRow(context.Background(), clientByIDQ, id).Scan(&c.ID, &c.UserID, &c.RefreshToken)
	if err != nil {
		t.Fatalf("couldn't get the client: %v", err)
	}
	return c
}

func clientByIDWithErr(t *testing.T, id string) (authent.Client, error) {
	t.Helper()
	c := authent.Client{}
	err := db.QueryRow(context.Background(), clientByIDQ, id).Scan(&c.ID, &c.UserID, &c.RefreshToken)
	if err != nil {
		return authent.Client{}, err
	}
	return c, nil
}
