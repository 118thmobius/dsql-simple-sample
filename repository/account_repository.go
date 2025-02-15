package repository

import (
	"context"
	"dsql-simple-sample/domain"
)

type AccountRepository interface {
	GetAccountByID(ctx context.Context, tx TxOrConn, userId string) (*domain.Account, error)
	UpdateAccount(ctx context.Context, tx TxOrConn, userId string, balance int) error
}
