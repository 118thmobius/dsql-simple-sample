package usecase

import (
	"context"
	"dsql-simple-sample/domain"
	"dsql-simple-sample/repository"
	"dsql-simple-sample/service"
	"errors"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type AccountUseCase interface {
	Transfer(ctx context.Context, from, to domain.Account, amount int) error
}

type AccountUseCaseImpl struct {
	accountRepository     repository.AccountRepository
	transactionRepository repository.TransactionRepository
	txDomainService       service.TransactionDomainService
	pool                  *pgxpool.Pool
}

func NewAccountUseCase(
	pool *pgxpool.Pool,
	accountRepository repository.AccountRepository,
	transactionRepository repository.TransactionRepository,
	txDomainService service.TransactionDomainService,
) *AccountUseCaseImpl {
	return &AccountUseCaseImpl{
		accountRepository:     accountRepository,
		transactionRepository: transactionRepository,
		txDomainService:       txDomainService,
		pool:                  pool,
	}
}

func (u *AccountUseCaseImpl) Transfer(ctx context.Context, from, to string, amount int) error {
	tx, err := u.pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer func() {
		if p := recover(); p != nil {
			_ = tx.Rollback(ctx)
			panic(p)
		} else if err != nil {
			_ = tx.Rollback(ctx)
		} else {
			err = tx.Commit(ctx)
		}
	}()

	fromAccount, err := u.accountRepository.GetAccountByID(ctx, txAdapter{tx}, from)
	if err != nil {
		return err
	}

	toAccount, err := u.accountRepository.GetAccountByID(ctx, txAdapter{tx}, to)
	if err != nil {
		return err
	}

	if u.txDomainService.CanTransfer(fromAccount, amount) {
		fromAccount.Balance -= amount
		toAccount.Balance += amount

		if err := u.accountRepository.UpdateAccount(ctx, txAdapter{tx}, fromAccount.UserId, fromAccount.Balance); err != nil {
			return err
		}
		if err := u.accountRepository.UpdateAccount(ctx, txAdapter{tx}, toAccount.UserId, toAccount.Balance); err != nil {
			return err
		}

		transaction := domain.Transaction{
			FromId: fromAccount.UserId,
			ToId:   toAccount.UserId,
			Amount: amount,
		}
		if err := u.transactionRepository.CreateTransaction(ctx, txAdapter{tx}, transaction.FromId, transaction.ToId, transaction.Amount); err != nil {
			return err
		}
		return nil
	}
	return errors.New("insufficient funds")
}

type txAdapter struct {
	tx pgx.Tx
}

func (t txAdapter) Exec(ctx context.Context, sql string, args ...interface{}) (repository.CommandTag, error) {
	cmdTag, err := t.tx.Exec(ctx, sql, args...)
	return cmdTag, err
}

func (t txAdapter) QueryRow(ctx context.Context, sql string, args ...interface{}) (repository.Row, error) {
	row := t.tx.QueryRow(ctx, sql, args...)
	return rowAdapter{row: row}, nil
}

type rowAdapter struct {
	row pgx.Row
}

func (r rowAdapter) Scan(dest ...interface{}) error {
	return r.row.Scan(dest...)
}
