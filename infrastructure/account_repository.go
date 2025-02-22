package infrastructure

import (
	"context"
	"dsql-simple-sample/domain"
	"dsql-simple-sample/infrastructure/db"
	"dsql-simple-sample/repository"
	"errors"
)

type AccountRepositoryImpl struct {
}

func NewAccountRepository() repository.AccountRepository {
	return &AccountRepositoryImpl{}
}

func (r AccountRepositoryImpl) GetAccountByID(ctx context.Context, q db.Queryer, userId string) (*domain.Account, error) {
	row := q.QueryRow(ctx, `
		SELECT user_id, city, balance
		FROM simple_account
		WHERE user_id = $1
		`, userId)
	var account domain.Account
	err := row.Scan(
		&account.UserId,
		&account.City,
		&account.Balance,
	)
	if err != nil {
		return nil, err
	}
	return &account, nil
}

func (r AccountRepositoryImpl) UpdateAccount(ctx context.Context, q db.Queryer, userId string, balance int) error {
	cmdTag, err := q.Exec(ctx, `
		UPDATE simple_account 
		SET balance = $1 
		WHERE user_id = $2
		`, balance, userId)
	if err != nil {
		return err
	}
	if affected := cmdTag.RowsAffected(); affected == 0 {
		return errors.New("no rows affected")
	}
	return nil
}
