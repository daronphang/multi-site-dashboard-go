package database

import (
	"context"
	"database/sql"
	"fmt"
	"multi-site-dashboard-go/internal/config"
	"net/url"
	"sync"

	"github.com/golang-migrate/migrate/v4/database"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	_ "github.com/jackc/pgx/v5/stdlib"
)


var (
	pgConnPool *pgxpool.Pool
	dbConn *sql.DB
	syncOnceConn sync.Once
	syncOncePool sync.Once
)

// When using a connection pool, if the existing connection to db is broken, it will automatically
// perform a reconnection for every new connection a thread requests from the pool.
// Hence, it is safe to initialize the connection pool once and reusing it for all threads.
func ProvidePgConnPool(ctx context.Context, cfg *config.Config) (*pgxpool.Pool, error) {
	var poolErr error
	syncOncePool.Do(func() {
		u := &url.URL{
			Scheme: "postgres",
			User: url.UserPassword(cfg.Postgres.Username, cfg.Postgres.Password),
			Host: fmt.Sprintf("%s:%d", cfg.Postgres.Server, cfg.Postgres.Port),
			Path: cfg.Postgres.DBName,
		}
		pool, err :=  pgxpool.New(ctx, u.String())
		if err != nil {
			poolErr = err
			return
		}
	
		if err := pool.Ping(ctx); err != nil {
			poolErr = err
			return
		}
		pgConnPool = pool
	})
	return pgConnPool, poolErr
}

func ProvidePgConn(cfg *config.Config) (*sql.DB, error) {
	var connErr error
	syncOnceConn.Do(func() {
		u := &url.URL{
			Scheme: "postgres",
			User: url.UserPassword(cfg.Postgres.Username, cfg.Postgres.Password),
			Host: fmt.Sprintf("%s:%d", cfg.Postgres.Server, cfg.Postgres.Port),
			// Path: cfg.Postgres.DBName,
		}
		conn, err := sql.Open("pgx", u.String())
		if err != nil {
			connErr = err
			return
		}
		dbConn = conn 
	})
	return dbConn, connErr
}

func ProvidePgDriver(conn *sql.DB) (database.Driver, error) {
	// For db migration.
	driver, err := postgres.WithInstance(conn, &postgres.Config{})
	if err != nil {
		return nil, err
	}
	return driver, nil
}

func SetupDatabase(cfg *config.Config) error {
	conn, err := ProvidePgConn(cfg)
	if err != nil {
		return err
	}
	_, err = conn.Exec(
		"SELECT 'CREATE DATABASE " + pgx.Identifier{cfg.Postgres.DBName}.Sanitize() + "' WHERE NOT EXISTS (SELECT FROM pg_database WHERE datname = $1)",
		cfg.Postgres.DBName,
	)
	if err != nil {
		return err
	}
	return nil
}