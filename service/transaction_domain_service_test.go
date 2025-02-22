package service_test

import (
    "testing"

    "dsql-simple-sample/domain"
    "dsql-simple-sample/service"
)

func TestTransactionDomainServiceImpl_CanTransfer(t *testing.T) {
    srv := service.NewTransactionDomainService()

    t.Run("Returns true if balance is greater than amount", func(t *testing.T) {
        fromAccount := &domain.Account{
            Balance: 1000,
        }
        amount := 500
        if !srv.CanTransfer(fromAccount, amount) {
            t.Errorf("Expected true when balance >= amount. Balance=%d, Amount=%d", fromAccount.Balance, amount)
        }
    })

    t.Run("Returns false if balance is less than amount", func(t *testing.T) {
        fromAccount := &domain.Account{
            Balance: 300,
        }
        amount := 500
        if srv.CanTransfer(fromAccount, amount) {
            t.Errorf("Expected false when balance < amount. Balance=%d, Amount=%d", fromAccount.Balance, amount)
        }
    })

    t.Run("Returns true if balance is exactly equal to amount", func(t *testing.T) {
        fromAccount := &domain.Account{
            Balance: 500,
        }
        amount := 500
        if !srv.CanTransfer(fromAccount, amount) {
            t.Errorf("Expected true when balance == amount. Balance=%d, Amount=%d", fromAccount.Balance, amount)
        }
    })
}