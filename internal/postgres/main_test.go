package postgres_test

import (
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/shanvl/garbage-events-service/internal/postgres"
)

// instance to be used in the tests
var db *sqlx.DB

func TestMain(m *testing.M) {
	os.Exit(testMain(m))
}

// TODO: add logger
// connects to the test db
func testMain(m *testing.M) int {
	// connect to the test db. Config values are hardcoded in order not to corrupt production db in case the wrong
	// compose file is used
	d, err := postgres.Connect(postgres.Config{
		Host:            "db",
		Database:        "testdb",
		User:            "root",
		Password:        "root",
		Port:            5432,
		MaxOpenConns:    20,
		MaxIdleConns:    20,
		ConnMaxLifetime: 5 * time.Minute,
	})
	if err != nil {
		fmt.Printf("couldn't connect to testdb: %s\n", err)
		return 1
	}
	defer d.Close()
	// assign the db instance to the global variable so that it can be used later in the tests
	db = d

	return m.Run()
}
