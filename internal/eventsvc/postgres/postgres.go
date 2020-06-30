// Package postgres manages postgres db
package postgres

import (
	"context"
	"time"

	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
)

// code which postgres returns on the violation of a foreign key
const foreignKeyViolationCode = "23503"

// Config allows to configure the db
type Config struct {
	Host, Database, User, Password string
	Port                           int
	MaxConns                       int
	MaxConnLifetime                time.Duration
	// SimpleProtocol becomes needed when using PgBouncer
	PreferSimpleProtocol bool
	// Logger allows to log the driver's events
	Logger pgx.Logger
}

// Connect establishes a connection to the db server using a provided config
func Connect(c Config) (*pgxpool.Pool, error) {
	// if we decide to use SimpleProtocol, the driver will have to escape params because statements won't be prepared
	escapeParams := "off"
	if c.PreferSimpleProtocol {
		escapeParams = "on"
	}

	// create a config for the pool
	conf, err := pgxpool.ParseConfig("")
	if err != nil {
		return nil, err
	}
	// max conns and max conn time
	conf.MaxConns, conf.MaxConnLifetime = int32(c.MaxConns), c.MaxConnLifetime
	// db credentials
	conf.ConnConfig.Host, conf.ConnConfig.Port, conf.ConnConfig.Database, conf.ConnConfig.User,
		conf.ConnConfig.Password = c.Host, uint16(c.Port), c.Database, c.User, c.Password
	// logger
	conf.ConnConfig.Logger = c.Logger
	// simple protocol
	conf.ConnConfig.PreferSimpleProtocol = c.PreferSimpleProtocol
	conf.ConnConfig.RuntimeParams = map[string]string{"standard_conforming_strings": escapeParams}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	// create the pool and ping the db
	return pgxpool.ConnectConfig(ctx, conf)
}
