package repository

import (
	"context"
	"dsql-simple-sample/infrastructure/db"
)

type TransactionRepository interface {
	CreateTransaction(ctx context.Context, q db.Queryer, fromId string, toId string, amount int) error
}
