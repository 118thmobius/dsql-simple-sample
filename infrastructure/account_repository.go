package infrastructure

import (
	"context"
	"dsql-simple-sample/domain"
	"dsql-simple-sample/repository"
	"errors"
)

type AccountRepositoryImpl struct {
}

func NewAccountRepository() repository.AccountRepository {
	return &AccountRepositoryImpl{}
}

func (r AccountRepositoryImpl) GetAccountByID(ctx context.Context, tx repository.TxOrConn, userId string) (*domain.Account, error) {
	row, err := tx.QueryRow(ctx, `
		SELECT user_id, city, balance
		FROM simple_account
		WHERE user_id = $1
		`, userId)
	if err != nil {
		return nil, err
	}
	var account domain.Account
	err = row.Scan(
		&account.UserId,
		&account.City,
		&account.Balance,
	)
	if err != nil {
		return nil, err
	}
	return &account, nil
}

func (r AccountRepositoryImpl) UpdateAccount(ctx context.Context, tx repository.TxOrConn, userId string, balance int) error {
	cmdTag, err := tx.Exec(ctx, `
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
