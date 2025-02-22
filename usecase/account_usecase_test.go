package usecase_test

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"dsql-simple-sample/domain"
	"dsql-simple-sample/infrastructure/db"
	"dsql-simple-sample/usecase"
)

//------------------------------
// Mock definitions
//------------------------------

// mockTxManager is a mock for TransactionManager
type mockTxManager struct {
	mock.Mock
}

// Do simulates a transactionless operation
func (m *mockTxManager) Do(ctx context.Context, fn func(ctx context.Context, q db.Queryer) error) error {
	args := m.Called(ctx, fn)
	// If needed, call the callback to simulate DB operation
	fnErr := fn(ctx, nil)
	returnValue := args.Error(0)
	if fnErr != nil {
		return fnErr
	}
	return returnValue
}

// DoTx simulates a transactional operation
func (m *mockTxManager) DoTx(ctx context.Context, fn func(ctx context.Context, q db.Queryer) error) error {
	args := m.Called(ctx, fn)
	fnErr := fn(ctx, nil)
	returnValue := args.Error(0)
	if fnErr != nil {
		return fnErr
	}
	return returnValue
}

// mockAccountRepo is a mock for AccountRepository
type mockAccountRepo struct {
	mock.Mock
}

// GetAccountByID simulates retrieving an account by user ID
func (m *mockAccountRepo) GetAccountByID(ctx context.Context, q db.Queryer, userId string) (*domain.Account, error) {
	args := m.Called(ctx, q, userId)
	acc, _ := args.Get(0).(*domain.Account)
	return acc, args.Error(1)
}

// UpdateAccount simulates updating an account's balance
func (m *mockAccountRepo) UpdateAccount(ctx context.Context, q db.Queryer, userId string, balance int) error {
	args := m.Called(ctx, q, userId, balance)
	return args.Error(0)
}

// mockTransactionRepo is a mock for TransactionRepository
type mockTransactionRepo struct {
	mock.Mock
}

// CreateTransaction simulates saving a transaction
func (m *mockTransactionRepo) CreateTransaction(ctx context.Context, q db.Queryer, fromId, toId string, amount int) error {
	args := m.Called(ctx, q, fromId, toId, amount)
	return args.Error(0)
}

// mockTxDomainService is a mock for TransactionDomainService
type mockTxDomainService struct {
	mock.Mock
}

// CanTransfer simulates checking if transfer is possible
func (m *mockTxDomainService) CanTransfer(account *domain.Account, amount int) bool {
	args := m.Called(account, amount)
	return args.Bool(0)
}

//------------------------------
// Test cases
//------------------------------

func Test_GetAccountByID_Success(t *testing.T) {
	ctx := context.Background()

	// Prepare mocks
	txManager := new(mockTxManager)
	accRepo := new(mockAccountRepo)
	txRepo := new(mockTransactionRepo)
	txDomainService := new(mockTxDomainService)

	// Create usecase instance
	uc := usecase.NewAccountUseCase(txManager, accRepo, txRepo, txDomainService)

	// Set mock behaviors
	txManager.
		On("Do", ctx, mock.Anything).
		Return(nil)

	expectedAccount := &domain.Account{
		UserId:  "user1",
		Balance: 1000,
	}
	accRepo.
		On("GetAccountByID", ctx, mock.Anything, "user1").
		Return(expectedAccount, nil)

	// Execute test
	acc, err := uc.GetAccountByID(ctx, "user1")

	// Verification
	assert.NoError(t, err)
	assert.Equal(t, expectedAccount, acc)
}

func Test_GetAccountByID_Error(t *testing.T) {
	ctx := context.Background()

	txManager := new(mockTxManager)
	accRepo := new(mockAccountRepo)
	txRepo := new(mockTransactionRepo)
	txDomainService := new(mockTxDomainService)

	uc := usecase.NewAccountUseCase(txManager, accRepo, txRepo, txDomainService)

	txManager.
		On("Do", ctx, mock.Anything).
		Return(nil)

	// Return an error simulating "user not found"
	accRepo.
		On("GetAccountByID", ctx, mock.Anything, "invalid_user").
		Return((*domain.Account)(nil), errors.New("user not found"))

	// Execute test
	acc, err := uc.GetAccountByID(ctx, "invalid_user")

	// Verification
	assert.Nil(t, acc)
	assert.Error(t, err)
}

func Test_Transfer_Success(t *testing.T) {
	ctx := context.Background()

	txManager := new(mockTxManager)
	accRepo := new(mockAccountRepo)
	txRepo := new(mockTransactionRepo)
	txDomainService := new(mockTxDomainService)

	uc := usecase.NewAccountUseCase(txManager, accRepo, txRepo, txDomainService)

	// Set mock behaviors
	txManager.
		On("DoTx", ctx, mock.Anything).
		Return(nil)

	fromAccount := &domain.Account{UserId: "fromUser", Balance: 1000}
	toAccount := &domain.Account{UserId: "toUser", Balance: 200}

	accRepo.
		On("GetAccountByID", ctx, mock.Anything, "fromUser").
		Return(fromAccount, nil)
	accRepo.
		On("GetAccountByID", ctx, mock.Anything, "toUser").
		Return(toAccount, nil)

	txDomainService.
		On("CanTransfer", fromAccount, 300).
		Return(true)

	accRepo.
		On("UpdateAccount", ctx, mock.Anything, "fromUser", 700).
		Return(nil)
	accRepo.
		On("UpdateAccount", ctx, mock.Anything, "toUser", 500).
		Return(nil)

	txRepo.
		On("CreateTransaction", ctx, mock.Anything, "fromUser", "toUser", 300).
		Return(nil)

	// Execute test
	err := uc.Transfer(ctx, "fromUser", "toUser", 300)

	// Verification
	assert.NoError(t, err)
	assert.Equal(t, 700, fromAccount.Balance)
	assert.Equal(t, 500, toAccount.Balance)
}

func Test_Transfer_InsufficientBalance(t *testing.T) {
	ctx := context.Background()

	txManager := new(mockTxManager)
	accRepo := new(mockAccountRepo)
	txRepo := new(mockTransactionRepo)
	txDomainService := new(mockTxDomainService)

	uc := usecase.NewAccountUseCase(txManager, accRepo, txRepo, txDomainService)

	txManager.
		On("DoTx", ctx, mock.Anything).
		Return(nil)

	fromAccount := &domain.Account{UserId: "fromUser", Balance: 100}
	toAccount := &domain.Account{UserId: "toUser", Balance: 200}

	accRepo.
		On("GetAccountByID", ctx, mock.Anything, "fromUser").
		Return(fromAccount, nil)
	accRepo.
		On("GetAccountByID", ctx, mock.Anything, "toUser").
		Return(toAccount, nil)

	// Simulate insufficient balance
	txDomainService.
		On("CanTransfer", fromAccount, 300).
		Return(false)

	// Execute test
	err := uc.Transfer(ctx, "fromUser", "toUser", 300)

	// Verification
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "Insufficient funds")
}
