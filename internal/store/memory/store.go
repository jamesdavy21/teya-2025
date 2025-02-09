package memory

import (
	"slices"

	"github.com/google/uuid"
	"github.com/jamesdavy21/teya-2025/internal/application"
)

type Store struct {
	accounts     map[uuid.UUID]application.Account
	transactions map[uuid.UUID][]application.Transaction
}

func NewInMemoryStore() *Store {
	return &Store{
		accounts:     make(map[uuid.UUID]application.Account),
		transactions: make(map[uuid.UUID][]application.Transaction),
	}
}

// GetAccount get account from in-memory store.
func (s *Store) GetAccount(accountID uuid.UUID) (*application.Account, error) {
	if account, ok := s.accounts[accountID]; ok {
		return &account, nil
	}

	return nil, application.ErrAccountNotFound
}

// SaveAccount save account into in-memory store.
func (s *Store) SaveAccount(account application.Account) error {
	s.accounts[account.ID] = account
	return nil
}

// SaveTransaction save a new transaction against an account and update account balance.
func (s *Store) SaveTransaction(accountID uuid.UUID, transaction application.Transaction) error {
	transactions := s.transactions[accountID]
	transactions = append(transactions, transaction)
	s.transactions[accountID] = transactions

	account := s.accounts[accountID]
	account.Balance += transaction.Amount
	s.accounts[accountID] = account

	return nil
}

// GetTransactions returns a list of transactions from the in-memory store. If page or limits exceeds limits returns nothing.
// Transactions get returned in descending order based on transaction.TransactionTime.
func (s *Store) GetTransactions(accountID uuid.UUID, page, limit int) ([]application.Transaction, int, error) {
	transactions := s.transactions[accountID]

	slices.SortFunc(transactions, func(a, b application.Transaction) int {
		return b.TransactionTime.Compare(a.TransactionTime)
	})

	nextPage := page + 1
	page = page * limit
	if page >= (len(transactions)) {
		page = len(transactions)
	}

	limit = limit + page
	if len(transactions) < limit {
		limit = len(transactions)
	}

	transactions = transactions[page:limit]
	if len(transactions) <= limit {
		nextPage = 0
	}

	return transactions, nextPage, nil
}
