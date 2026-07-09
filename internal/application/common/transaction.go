package common

import (
	"context"

	db "github.com/go-api/internal/infrastructure/postgres/generated"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type txKeyType string

const txKey txKeyType = "tx"

type TransactionManager interface {
	Execute(ctx context.Context, fn func(ctx context.Context) error) error
}

type transactionManager struct {
	db *pgxpool.Pool
}

func NewTransactionManager(db *pgxpool.Pool) TransactionManager {
	return &transactionManager{db: db}
}

func (t *transactionManager) Execute(ctx context.Context, fn func(ctx context.Context) error) error {
	tx, err := t.db.Begin(ctx)
	if err != nil {
		return err
	}

	txCtx := context.WithValue(ctx, txKey, tx)

	err = fn(txCtx)
	if err != nil {
		_ = tx.Rollback(ctx)
		return err
	}

	return tx.Commit(ctx)
}

func GetTx(ctx context.Context) pgx.Tx {
	tx, ok := ctx.Value(txKey).(pgx.Tx)
	if ok {
		return tx
	}
	return nil
}

func GetQueries(ctx context.Context, pool *pgxpool.Pool) *db.Queries {
	tx := GetTx(ctx)
	if tx != nil {
		return db.New(tx)
	}
	return db.New(pool)
}
