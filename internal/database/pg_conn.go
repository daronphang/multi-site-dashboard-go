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
	"github.com/jackc/pgx/v5/pgxpool"
	_ "github.com/jackc/pgx/v5/stdlib"
)


var (
	pgConnPool *pgxpool.Pool
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

func ProvidePgConn(cfg *config.Config, withPath bool) (*sql.DB, error) {
	var path = cfg.Postgres.DBName
	if !withPath {
		path = ""
	}
	u := &url.URL{
		Scheme: "postgres",
		User: url.UserPassword(cfg.Postgres.Username, cfg.Postgres.Password),
		Host: fmt.Sprintf("%s:%d", cfg.Postgres.Server, cfg.Postgres.Port),
		Path: path,
	}
	conn, err := sql.Open("pgx", u.String())
	if err != nil {
		return nil, err
	}
	return conn, nil
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
	conn, err := ProvidePgConn(cfg, false)
	if err != nil {
		return err
	}

	rv := conn.QueryRow("SELECT 1 FROM pg_database WHERE datname = $1", cfg.Postgres.DBName)
	if err := rv.Scan(); err == sql.ErrNoRows {
		_, err = conn.Exec("CREATE DATABASE " + cfg.Postgres.DBName)
		if err != nil {
			return err
		}
	}
	if err := conn.Close(); err != nil {
		return err
	}
	return nil
}