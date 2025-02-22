package usecase

import (
	"context"
	"dsql-simple-sample/domain"
	"dsql-simple-sample/infrastructure/db"
	"dsql-simple-sample/repository"
	"dsql-simple-sample/service"
	"errors"
	"github.com/jackc/pgx/v5/pgxpool"
)

type AccountUseCase interface {
	Transfer(ctx context.Context, from, to domain.Account, amount int) error
	GetAccountByID(ctx context.Context, userId string) (*domain.Account, error)
}

type AccountUseCaseImpl struct {
	pool                  *pgxpool.Pool
	txManager             db.TransactionManager
	accountRepository     repository.AccountRepository
	transactionRepository repository.TransactionRepository
	txDomainService       service.TransactionDomainService
}

func NewAccountUseCase(
	txManager db.TransactionManager,
	accountRepository repository.AccountRepository,
	transactionRepository repository.TransactionRepository,
	txDomainService service.TransactionDomainService,
) *AccountUseCaseImpl {
	return &AccountUseCaseImpl{
		txManager:             txManager,
		accountRepository:     accountRepository,
		transactionRepository: transactionRepository,
		txDomainService:       txDomainService,
	}
}

func (u *AccountUseCaseImpl) GetAccountByID(ctx context.Context, userId string) (*domain.Account, error) {
	var account *domain.Account
	err := u.txManager.Do(ctx, func(ctx context.Context, q db.Queryer) error {
		acc, err := u.accountRepository.GetAccountByID(ctx, q, userId)
		if err != nil {
			return err
		}
		account = acc
		return nil
	})
	if err != nil {
		return nil, err
	}
	return account, nil
}

func (u *AccountUseCaseImpl) Transfer(ctx context.Context, from, to string, amount int) error {
	return u.txManager.DoTx(ctx, func(ctx context.Context, q db.Queryer) error {
		fromAccount, err := u.accountRepository.GetAccountByID(ctx, q, from)
		if err != nil {
			return err
		}

		toAccount, err := u.accountRepository.GetAccountByID(ctx, q, to)
		if err != nil {
			return err
		}

		if u.txDomainService.CanTransfer(fromAccount, amount) {
			fromAccount.Balance -= amount
			toAccount.Balance += amount

			if err := u.accountRepository.UpdateAccount(ctx, q, fromAccount.UserId, fromAccount.Balance); err != nil {
				return err
			}
			if err := u.accountRepository.UpdateAccount(ctx, q, toAccount.UserId, toAccount.Balance); err != nil {
				return err
			}

			transaction := domain.Transaction{
				FromId: fromAccount.UserId,
				ToId:   toAccount.UserId,
				Amount: amount,
			}
			if err := u.transactionRepository.CreateTransaction(ctx, q, transaction.FromId, transaction.ToId, transaction.Amount); err != nil {
				return err
			}
			return nil
		}
		return errors.New("Insufficient funds")
	})
}
