package domain

import (
	"reflect"
	"testing"
)

func TestNewTransaction(t *testing.T) {
	fromAccount := Account{UserId: "user_123"}
	toAccount := Account{UserId: "user_456"}
	amount := 100
	transaction := NewTransaction(fromAccount, toAccount, amount)

	expectedTransaction := &Transaction{
		FromId: "user_123",
		ToId:   "user_456",
		Amount: 100,
	}

	if !reflect.DeepEqual(transaction, expectedTransaction) {
		t.Errorf("NewTransaction() failed, expected %v, got %v", expectedTransaction, transaction)
	}
}
