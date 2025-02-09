package application_test

import (
	"testing"

	"github.com/google/uuid"
	"github.com/jamesdavy21/teya-2025/internal/application"
	"github.com/jamesdavy21/teya-2025/internal/store/memory"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetAccountCreatesAccountIfNotFound(t *testing.T) {
	store := memory.NewInMemoryStore()
	m := application.NewTransactionManager(store)

	accountID := uuid.New()
	account, err := m.GetAccount(accountID)
	require.NoError(t, err)

	assert.Equal(t, accountID, account.ID)
	assert.Equal(t, float64(0), account.Balance)
}

func TestGetAccountReturnsExistingAccount(t *testing.T) {
	store := memory.NewInMemoryStore()

	existingAccount := application.Account{
		ID:      uuid.New(),
		Balance: 10,
	}
	require.NoError(t, store.SaveAccount(existingAccount))

	m := application.NewTransactionManager(store)
	account, err := m.GetAccount(existingAccount.ID)
	require.NoError(t, err)

	assert.Equal(t, existingAccount.ID, account.ID)
	assert.Equal(t, existingAccount.Balance, account.Balance)
}

func TestTransactionManager_AddDeposit(t *testing.T) {
	store := memory.NewInMemoryStore()
	m := application.NewTransactionManager(store)

	accountID := uuid.New()
	transaction, err := m.AddDeposit(accountID, 10.555)
	require.NoError(t, err)
	assert.Equal(t, 10.55, transaction.Amount)
	assert.Equal(t, application.TransactionDeposit, transaction.TransactionType)

	account, err := m.GetAccount(accountID)
	require.NoError(t, err)

	assert.Equal(t, 10.55, account.Balance)
}

func TestTransactionManager_AddWithdrawSuccess(t *testing.T) {
	store := memory.NewInMemoryStore()
	m := application.NewTransactionManager(store)

	accountID := uuid.New()
	_, err := m.AddDeposit(accountID, 20.6)
	require.NoError(t, err)

	transaction, err := m.AddWithdrawal(accountID, 10.71)
	require.NoError(t, err)
	assert.Equal(t, -10.71, transaction.Amount)
	assert.Equal(t, application.TransactionWithdrawal, transaction.TransactionType)

	account, err := m.GetAccount(accountID)
	require.NoError(t, err)

	assert.Equal(t, 9.89, account.Balance)
}

func TestTransactionManager_AddWithdrawAccountNotFound(t *testing.T) {
	store := memory.NewInMemoryStore()
	m := application.NewTransactionManager(store)

	accountID := uuid.New()
	_, err := m.AddWithdrawal(accountID, 10)
	require.ErrorIs(t, err, application.ErrAccountNotFound)
}
