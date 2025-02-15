package repository

import (
	"context"
)

type TransactionRepository interface {
	CreateTransaction(ctx context.Context, tx TxOrConn, fromId string, toId string, amount int) error
}
