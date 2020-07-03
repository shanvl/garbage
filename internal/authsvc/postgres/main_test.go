package postgres_test

import (
	"log"
	"os"
	"testing"
	"time"

	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/shanvl/garbage/internal/eventsvc/postgres"
)

// instance to be used in the tests
var db *pgxpool.Pool

func TestMain(m *testing.M) {
	os.Exit(testMain(m))
}

// connects to the test db
func testMain(m *testing.M) int {
	// connect to the test db. Config values are hardcoded in order not to corrupt production db in case the wrong
	// compose file is used
	d, err := postgres.Connect(postgres.Config{
		Host:            "authsvc_db",
		Database:        "testdb",
		User:            "root",
		Password:        "root",
		Port:            5432,
		MaxConns:        20,
		MaxConnLifetime: 5 * time.Minute,
	})
	if err != nil {
		log.Printf("couldn't connect to testdb: %s\n", err)
		return 1
	}
	defer d.Close()
	// assign the db instance to the global variable so that it can be used later in the tests
	db = d

	return m.Run()
}
