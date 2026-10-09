package db

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
)

// TxRunner runs fn inside a database transaction. The Querier passed to fn is scoped to the
// transaction, which is rolled back if fn returns an error and committed otherwise.
type TxRunner func(ctx context.Context, fn func(q Querier) error) error

// NewTxRunner returns a TxRunner that opens its transactions on pool, so callers can run
// transactions without holding the pool itself.
func NewTxRunner(pool *pgxpool.Pool) TxRunner {
	return func(ctx context.Context, fn func(q Querier) error) error {
		tx, err := pool.Begin(ctx)
		if err != nil {
			return fmt.Errorf("begin tx: %w", err)
		}
		defer tx.Rollback(ctx) //nolint:errcheck

		if err := fn(New(tx)); err != nil {
			return err
		}
		return tx.Commit(ctx)
	}
}
