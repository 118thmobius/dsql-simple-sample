package db

import (
	"context"
	"github.com/jackc/pgx/v5/pgxpool"
)

type TransactionManager interface {
	Do(ctx context.Context, fn func(ctx context.Context, q Queryer) error) error
	DoTx(ctx context.Context, fn func(ctx context.Context, q Queryer) error) error
}

type txManager struct {
	pool *pgxpool.Pool
}

func NewTransactionManager(pool *pgxpool.Pool) TransactionManager {
	return &txManager{pool: pool}
}

func (t txManager) Do(ctx context.Context, fn func(ctx context.Context, q Queryer) error) error {
	conn, err := t.pool.Acquire(ctx)
	if err != nil {
		return err
	}
	defer conn.Release()

	return fn(ctx, conn)
}

func (t txManager) DoTx(ctx context.Context, fn func(ctx context.Context, q Queryer) error) error {
	conn, err := t.pool.Acquire(ctx)
	if err != nil {
		return err
	}
	defer conn.Release()

	tx, err := conn.Begin(ctx)
	if err != nil {
		return err
	}

	defer func() {
		if p := recover(); p != nil {
			_ = tx.Rollback(ctx)
			panic(p)
		}
	}()

	if err := fn(ctx, tx); err != nil {
		_ = tx.Rollback(ctx)
		return err
	}

	return tx.Commit(ctx)
}
