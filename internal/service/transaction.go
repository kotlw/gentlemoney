package service

import (
	"fmt"

	"github.com/kotlw/gentlemoney/internal/model"
	"github.com/kotlw/gentlemoney/internal/storage/inmemory"
	"github.com/kotlw/gentlemoney/internal/storage/sqlite"
)

// Transaction service contains business logic related to model.Transaction.
type Transaction struct {
	persistentStorage *sqlite.Transaction
	inmemoryStorage   *inmemory.Transaction
}

// NewCurrency returns Transaction service.
func NewTransaction(
	persistentStorage *sqlite.Transaction,
	inmemoryStorage *inmemory.Transaction,
	categoryService *Category,
	accountService *Account) (*Transaction, error) {

	a := &Transaction{
		persistentStorage: persistentStorage,
		inmemoryStorage:   inmemoryStorage,
	}

	if err := a.Init(categoryService, accountService); err != nil {
		return nil, fmt.Errorf("c.Init: %w", err)
	}

	return a, nil
}

// Init initialize inmemory storage with data from persistent storage. It is also links existing
// categories and accounts to corresponding fields of model.Transaction.
func (s *Transaction) Init(categoryService *Category, accountService *Account) error {
	tt, err := s.persistentStorage.GetAll()
	if err != nil {
		return fmt.Errorf("s.persistentStorage.GetAll: %w", err)
	}

	for _, t := range tt {
		t.Category = categoryService.GetByID(t.Category.ID)
		t.Account = accountService.GetByID(t.Account.ID)
	}

	s.inmemoryStorage.Init(tt)

	return nil
}

// Insert appends transaction to both persistent and inmemory storages.
func (s *Transaction) Insert(t *model.Transaction) error {
	id, err := s.persistentStorage.Insert(t)
	if err != nil {
		return fmt.Errorf("s.persistentStorage.Insert: %w", err)
	}

	t.ID = id
	s.inmemoryStorage.Insert(t)

	return nil
}

// Update updates transaction in persistent and inmemory storages.
func (s *Transaction) Update(t *model.Transaction) error {
	if err := s.persistentStorage.Update(t); err != nil {
		return fmt.Errorf("s.persistentStorage.Update: %w", err)
	}

	s.inmemoryStorage.Update(t)

	return nil
}

// Delete deletes transaction from inmemory and persistent storages.
func (s *Transaction) Delete(t *model.Transaction) error {
	if err := s.persistentStorage.Delete(t.ID); err != nil {
		return fmt.Errorf("s.persistentStorage.Delete: %w", err)
	}
	s.inmemoryStorage.Delete(t)
	return nil
}

// GetAll returns all transactions.
func (s *Transaction) GetAll() []*model.Transaction {
	return s.inmemoryStorage.GetAll()
}

// GetByID returns transaction by given model.Transaction.ID.
func (s *Transaction) GetByID(id int64) *model.Transaction {
	return s.inmemoryStorage.GetByID(id)
}
