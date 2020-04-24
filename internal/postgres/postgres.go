// Package postgres manages postgres db
package postgres

import (
	"time"

	"github.com/jackc/pgx"
	"github.com/jackc/pgx/stdlib"
	"github.com/jmoiron/sqlx"
)

// Config allows to configure the db
type Config struct {
	Host, Database, User, Password  string
	Port                            int
	MaxOpenConns, MaxIdleConns      int
	AcquireTimeout, ConnMaxLifetime time.Duration
	// SimpleProtocol becomes needed when using PgBouncer
	PreferSimpleProtocol bool
	// Logger allows to log the driver's events
	Logger pgx.Logger
}

// Connect establishes a connection to the db server using a provided config
func Connect(c Config) (*sqlx.DB, error) {
	// if we are using SimpleProtocol, the driver must escape params because statements won't be prepared
	escapeParams := "off"
	if c.PreferSimpleProtocol {
		escapeParams = "on"
	}

	// set up the pgx connection pool
	pool, err := pgx.NewConnPool(pgx.ConnPoolConfig{
		ConnConfig: pgx.ConnConfig{
			Host:                 c.Host,
			Port:                 uint16(c.Port),
			Database:             c.Database,
			User:                 c.User,
			Password:             c.Password,
			Logger:               c.Logger,
			PreferSimpleProtocol: c.PreferSimpleProtocol,
			RuntimeParams: map[string]string{
				"standard_conforming_strings": escapeParams,
			},
		},
		MaxConnections: c.MaxOpenConns,
		AfterConnect:   nil,
		AcquireTimeout: c.AcquireTimeout,
	})
	if err != nil {
		return nil, err
	}
	pgxDB := stdlib.OpenDBFromPool(pool, stdlib.OptionPreferSimpleProtocol(c.PreferSimpleProtocol))

	// set up sqlx
	db := sqlx.NewDb(pgxDB, "pgx")
	db.SetMaxIdleConns(c.MaxOpenConns)
	db.SetMaxIdleConns(c.MaxIdleConns)
	db.SetConnMaxLifetime(c.ConnMaxLifetime)

	// ping the db
	err = db.Ping()
	if err != nil {
		return nil, err
	}

	return db, nil
}
