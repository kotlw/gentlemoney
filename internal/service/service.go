package service

import (
	"fmt"

	"github.com/kotlw/gentlemoney/internal/storage/inmemory"
	"github.com/kotlw/gentlemoney/internal/storage/sqlite"
)

// Service is a facade structure which aggregates all Services. It is used for convenience.
type Service struct {
	category    *Category
	currency    *Currency
	account     *Account
	transaction *Transaction
}

// New returns new Service. 
func New(ps *sqlite.SqliteStorage, is *inmemory.InmemoryStorage) (s *Service, err error) {
	s = &Service{}

	if s.category, err = NewCategory(ps.Category(), is.Category()); err != nil {
		return nil, fmt.Errorf("NewCategory: %w", err)
	}
	if s.currency, err = NewCurrency(ps.Currency(), is.Currency()); err != nil {
		return nil, fmt.Errorf("NewCurrency: %w", err)
	}
	if s.account, err = NewAccount(ps.Account(), is.Account(), s.currency); err != nil {
		return nil, fmt.Errorf("NewAccount: %w", err)
	}
	if s.transaction, err = NewTransaction(ps.Transaction(), is.Transaction(), s.category, s.account); err != nil {
		return nil, fmt.Errorf("NewTransaction: %w", err)
	}

	return s, nil
}

// Category returns category service.
func (s *Service) Category() *Category {
	return s.category
}

// Currency returns currency service.
func (s *Service) Currency() *Currency {
	return s.currency
}

// Account returns account service.
func (s *Service) Account() *Account {
	return s.account
}

// Transaction returns transaction service.
func (s *Service) Transaction() *Transaction {
	return s.transaction
}
