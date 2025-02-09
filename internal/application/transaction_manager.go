package application

import (
	"errors"
	"math"
	"time"

	"github.com/google/uuid"
)

type Store interface {
	GetAccount(accountID uuid.UUID) (*Account, error)
	SaveAccount(account Account) error
	SaveTransaction(accountID uuid.UUID, transaction Transaction) error
	GetTransactions(accountID uuid.UUID, page, limit int) ([]Transaction, int, error)
}

type TransactionManager struct {
	store Store
}

func NewTransactionManager(store Store) *TransactionManager {
	return &TransactionManager{
		store: store,
	}
}

type Account struct {
	ID      uuid.UUID
	Balance float64
}

type Transaction struct {
	TransactionID   uuid.UUID
	Amount          float64
	TransactionTime time.Time
	TransactionType TransactionType
}

type TransactionType string

const (
	TransactionDeposit    TransactionType = "deposit"
	TransactionWithdrawal TransactionType = "withdrawal"
)

// GetAccount returns the current account nad its balance. If the account doesn't exist, the account is created.
func (t TransactionManager) GetAccount(accountID uuid.UUID) (*Account, error) {
	account, err := t.store.GetAccount(accountID)
	if err != nil {
		if !errors.Is(err, ErrAccountNotFound) {
			return nil, err
		}

		account = &Account{
			ID: accountID,
		}
		if err = t.store.SaveAccount(*account); err != nil {
			return nil, err
		}
	}

	return account, nil
}

// AddDeposit add a new deposit transaction to the account. If the account doesn't exist, the account is created.
func (t TransactionManager) AddDeposit(accountID uuid.UUID, amount float64) (*Transaction, error) {
	account, err := t.store.GetAccount(accountID)
	if err != nil {
		if !errors.Is(err, ErrAccountNotFound) {
			return nil, err
		}

		account = &Account{
			ID: accountID,
		}
		if err = t.store.SaveAccount(*account); err != nil {
			return nil, err
		}
	}

	deposit := Transaction{
		TransactionID:   uuid.New(),
		Amount:          math.Floor(amount*100) / 100,
		TransactionTime: time.Now().UTC(),
		TransactionType: TransactionDeposit,
	}

	if err = t.store.SaveTransaction(account.ID, deposit); err != nil {
		return nil, err
	}

	return &deposit, nil
}

// AddWithdrawal adds a new withdrawal transaction to an existing account. A withdrawal can only take place if account balance remains positive afterward.
func (t TransactionManager) AddWithdrawal(accountID uuid.UUID, amount float64) (*Transaction, error) {
	account, err := t.store.GetAccount(accountID)
	if err != nil {
		return nil, err
	}

	amount = math.Floor(amount*100) / 100
	if account.Balance < amount {
		return nil, ErrNotEnoughFunds
	}

	deposit := Transaction{
		TransactionID:   uuid.New(),
		Amount:          -amount,
		TransactionTime: time.Now().UTC(),
		TransactionType: TransactionWithdrawal,
	}

	if err = t.store.SaveTransaction(account.ID, deposit); err != nil {
		return nil, err
	}

	return &deposit, nil
}

// GetTransactions returns all transaction for an existing account.
func (t TransactionManager) GetTransactions(accountID uuid.UUID, page, limit int) ([]Transaction, int, error) {
	_, err := t.store.GetAccount(accountID)
	if err != nil {
		return nil, 0, err
	}

	transactions, nextPage, err := t.store.GetTransactions(accountID, page, limit)
	if err != nil {
		return nil, 0, err
	}

	return transactions, nextPage, nil
}
