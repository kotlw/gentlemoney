package storage

import (
	"database/sql"
	"fmt"
)

// Storage is a facade structure which aggregates all storages. It is used for convenience.
type Storage struct {
	category    *Category
	currency    *Currency
	account     *Account
	transaction *Transaction
}

// New creates object which aggregates all storages.
func New(db *sql.DB) (s *Storage, err error) {
	s = &Storage{}

	if s.category, err = NewCategory(db); err != nil {
		return nil, fmt.Errorf("NewCategory: %w", err)
	}
	if s.currency, err = NewCurrency(db); err != nil {
		return nil, fmt.Errorf("NewCurrency: %w", err)
	}
	if s.account, err = NewAccount(db); err != nil {
		return nil, fmt.Errorf("NewAccount: %w", err)
	}
	if s.transaction, err = NewTransaction(db); err != nil {
		return nil, fmt.Errorf("NewTransaction: %w", err)
	}

	return s, nil
}

// Category returns category storage.
func (s *Storage) Category() *Category {
	return s.category
}

// Currency returns currency storage.
func (s *Storage) Currency() *Currency {
	return s.currency
}

// Account returns account storage.
func (s *Storage) Account() *Account {
	return s.account
}

// Transaction returns transaction storage.
func (s *Storage) Transaction() *Transaction {
	return s.transaction
}
