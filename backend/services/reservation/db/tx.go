package db

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
)

// TxFunc is a function that runs within a transaction, receiving a Querier scoped to that transaction.
type TxFunc func(q Querier) error

// WithTx runs fn inside a database transaction. The Querier passed to fn is scoped to the transaction.
// If fn returns an error, the transaction is rolled back. Otherwise, it is committed.
func WithTx(ctx context.Context, pool *pgxpool.Pool, fn TxFunc) error {
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
