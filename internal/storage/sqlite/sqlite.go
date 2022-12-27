package sqlite

import (
	"database/sql"
	"fmt"
)

// SqliteStorage is a facade structure which aggregates all sqlite storages. It is used for convenience.
type SqliteStorage struct {
	category    *Category
	currency    *Currency
	account     *Account
	transaction *Transaction
}

// New creates object which aggregates all storages.
func New(db *sql.DB) (s *SqliteStorage, err error) {
	s = &SqliteStorage{}

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

// Category returns category sqlite storage.
func (s *SqliteStorage) Category() *Category {
	return s.category
}

// Currency returns currency sqlite storage.
func (s *SqliteStorage) Currency() *Currency {
	return s.currency
}

// Account returns account sqlite storage.
func (s *SqliteStorage) Account() *Account {
	return s.account
}

// Transaction returns transaction sqlite storage.
func (s *SqliteStorage) Transaction() *Transaction {
	return s.transaction
}
