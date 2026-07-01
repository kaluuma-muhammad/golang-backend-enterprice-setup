package postgres

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

type TransactionManager struct {
	pool *pgxpool.Pool
}

func NewTransactionManager(pool *pgxpool.Pool) *TransactionManager {

	return &TransactionManager{pool: pool}
}

func (t *TransactionManager) Execute(
	ctx context.Context,
	fn func(ctx context.Context) error,
) error {

	tx, err := t.pool.Begin(ctx)
	if err != nil {
		return err
	}

	defer tx.Rollback(ctx)

	txCtx := context.WithValue(
		ctx,
		"tx",
		tx,
	)

	if err := fn(txCtx); err != nil {
		return err
	}

	return tx.Commit(ctx)
}