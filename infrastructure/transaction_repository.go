package infrastructure

import (
	"context"
	"dsql-simple-sample/infrastructure/db"
	"dsql-simple-sample/repository"
)

type TransactionRepositoryImpl struct {
}

func NewTransactionRepository() repository.TransactionRepository {
	return &TransactionRepositoryImpl{}
}

func (t TransactionRepositoryImpl) CreateTransaction(ctx context.Context, q db.Queryer, fromId string, toId string, amount int) error {
	_, err := q.Exec(ctx, `
		INSERT INTO simple_transaction (from_id,to_id,amount)
		VALUES ($1,$2,$3)
		`, fromId, toId, amount)
	return err
}
