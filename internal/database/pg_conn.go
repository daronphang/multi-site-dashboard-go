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
	lock sync.Mutex
)

// When using a connection pool, if the existing connection to db is broken, it will automatically
// perform a reconnection for every new connection a thread requests from the pool.
// Hence, it is safe to initialize the connection pool once and reusing it for all threads.
func ProvidePgConnPool(ctx context.Context, conf *config.Config) (*pgxpool.Pool, error) {
	lock.Lock()
	defer lock.Unlock()
	// Use established connection pool if exists.
	if pgConnPool != nil {
		if err := pgConnPool.Ping(ctx); err == nil {
			return pgConnPool, nil
		}
	}

	pool, err := createConnPool(ctx, conf)
	if err != nil {
		return nil, err
	}

	// Cache connection pool.
	pgConnPool = pool
	return pgConnPool, nil
}

func ProvidePgDriver(conf *config.Config) (database.Driver, error) {
	u := &url.URL{
		Scheme: "postgres",
		User: url.UserPassword(conf.Postgres.Username, conf.Postgres.Password),
		Host: fmt.Sprintf("%s:%d", conf.Postgres.Server, conf.Postgres.Port),
		Path: conf.Postgres.DBName,
	}
	conn, err := sql.Open("pgx", u.String())
	if err != nil {
		return nil, err
	}
	driver, err := postgres.WithInstance(conn, &postgres.Config{})
	if err != nil {
		return nil, err
	}

	return driver, nil
}

func createConnPool(ctx context.Context, conf *config.Config) (*pgxpool.Pool, error) {
	u := &url.URL{
		Scheme: "postgres",
		User: url.UserPassword(conf.Postgres.Username, conf.Postgres.Password),
		Host: fmt.Sprintf("%s:%d", conf.Postgres.Server, conf.Postgres.Port),
		Path: conf.Postgres.DBName,
	}
	pool, err :=  pgxpool.New(ctx, u.String())
	if err != nil {
		return nil, err
	}

	if err := pool.Ping(ctx); err != nil {
		return nil, err
	}
	
	return pool, nil
}