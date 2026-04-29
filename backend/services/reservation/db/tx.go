package db

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// TxFunc is a function that runs within a transaction, receiving a Querier and the underlying pgx.Tx.
type TxFunc func(q Querier, tx pgx.Tx) error

// WithTx runs fn inside a database transaction. The Querier and pgx.Tx passed to fn are scoped to the transaction.
// If fn returns an error, the transaction is rolled back. Otherwise, it is committed.
func WithTx(ctx context.Context, pool *pgxpool.Pool, fn TxFunc) error {
	tx, err := pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("begin tx: %w", err)
	}
	defer tx.Rollback(ctx) //nolint:errcheck

	if err := fn(New(tx), tx); err != nil {
		return err
	}
	return tx.Commit(ctx)
}
