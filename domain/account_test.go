package domain

import "testing"

func TestAccountFields(t *testing.T) {
	account := Account{
		UserId:  "user123",
		City:    "Tokyo",
		Balance: 1000,
	}

	if account.UserId != "user123" {
		t.Errorf("expected UserId to be 'user123', but got '%s'", account.UserId)
	}

	if account.City != "Tokyo" {
		t.Errorf("expected City to be 'Tokyo', but got '%s'", account.City)
	}

	if account.Balance != 1000 {
		t.Errorf("expected Balance to be 1000, but got %d", account.Balance)
	}
}
