package repository

import (
	"context"
	"dsql-simple-sample/domain"
	"dsql-simple-sample/infrastructure/db"
)

type AccountRepository interface {
	GetAccountByID(ctx context.Context, q db.Queryer, userId string) (*domain.Account, error)
	UpdateAccount(ctx context.Context, q db.Queryer, userId string, balance int) error
}
