package repository

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5/pgxpool"
)

type ExtQuerier interface {
	Querier
	ExecWithPgTx(ctx context.Context, cb func(Querier) (interface{}, error)) func() (interface{}, error)
}

// Enables multiple executions to be atomic by treating them as a single transaction.
func (q *Queries) ExecWithPgTx(ctx context.Context, cb func(Querier) (interface{}, error)) func() (interface{}, error) {
	return func() (interface{}, error) {
		db, ok := q.db.(*pgxpool.Pool)
		if !ok {
			return nil, errors.New("db is not an pgxpool.Pool")
		}
		tx, err := db.Begin(ctx)
		if err != nil {
			return nil, err
		}
		defer tx.Rollback(ctx)

		qtx := q.WithTx(tx)
		rv, err := cb(qtx)
		if err != nil {
			return nil, err
		}
		if err := tx.Commit(ctx); err != nil {
			return nil, err
		}
		return rv, nil
	}
}