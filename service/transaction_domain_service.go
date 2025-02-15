package service

import "dsql-simple-sample/domain"

type TransactionDomainService interface {
	CanTransfer(fromAccount *domain.Account, amount int) bool
}

type TransactionDomainServiceImpl struct {
}

func NewTransactionDomainService() TransactionDomainService {
	return &TransactionDomainServiceImpl{}
}

func (s *TransactionDomainServiceImpl) CanTransfer(fromAccount *domain.Account, amount int) bool {
	return fromAccount.Balance >= amount
}
