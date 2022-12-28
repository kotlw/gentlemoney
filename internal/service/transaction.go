package service

import (
	"fmt"
	"sort"

	"github.com/kotlw/gentlemoney/internal/model"
	"github.com/kotlw/gentlemoney/internal/storage/inmemory"
	"github.com/kotlw/gentlemoney/internal/storage/sqlite"
)

// transactionList is a wrapper to perform sort by model.Transaction.Date.
type transactionList []*model.Transaction

func (tt transactionList) Len() int           { return len(tt) }
func (tt transactionList) Less(i, j int) bool { return tt[i].Date.After(tt[j].Date) }
func (tt transactionList) Swap(i, j int)      { tt[i], tt[j] = tt[j], tt[i] }

// Transaction service contains business logic related to model.Transaction.
type Transaction struct {
	persistantStorage *sqlite.Transaction
	inmemoryStorage   *inmemory.Transaction
}

// NewCurrency returns Transaction service.
func NewTransaction(
	persistantStorage *sqlite.Transaction,
	inmemoryStorage *inmemory.Transaction,
	categoryService *Category,
	accountService *Account) (*Transaction, error) {

	a := &Transaction{
		persistantStorage: persistantStorage,
		inmemoryStorage:   inmemoryStorage,
	}

	if err := a.Init(categoryService, accountService); err != nil {
		return nil, fmt.Errorf("c.Init: %w", err)
	}

	return a, nil
}

// Init initialize inmemory storage with data from persistant storage. It is also links existing
// categories and accounts to corresponding fields of model.Transaction.
func (s *Transaction) Init(categoryService *Category, accountService *Account) error {
	tt, err := s.persistantStorage.GetAll()
	if err != nil {
		return fmt.Errorf("s.persistantStorage.GetAll: %w", err)
	}

    for _, t := range tt {
        t.Category = categoryService.GetByID(t.Category.ID)
        t.Account = accountService.GetByID(t.Account.ID)
    }

	s.inmemoryStorage.Init(tt)

	return nil
}

// Insert appends transaction to both persistant and inmemory storages.
func (s *Transaction) Insert(c *model.Transaction) error {
	id, err := s.persistantStorage.Insert(c)
	if err != nil {
		return fmt.Errorf("s.persistantStorage.Insert: %w", err)
	}

	c.ID = id
	s.inmemoryStorage.Insert(c)

	return nil
}

// Update updates transaction in persistant storage. Since GetAll returns pointers to inmemory data
// after update the category we need to update it in persistant storage as well.
func (s *Transaction) Update(c *model.Transaction) error {
	if err := s.persistantStorage.Update(c); err != nil {
		return fmt.Errorf("s.persistantStorage.Update: %w", err)
	}
	return nil
}

// Delete deletes transaction from inmemory and persistant storages.
func (s *Transaction) Delete(c *model.Transaction) error {
	if err := s.persistantStorage.Delete(c.ID); err != nil {
		return fmt.Errorf("s.persistantStorage.Delete: %w", err)
	}
	s.inmemoryStorage.Delete(c)
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

// GetAllSorted returns all transactions sorted by model.Transaction.Date.
func (s *Transaction) GetAllSorted() []*model.Transaction {
	cc := s.inmemoryStorage.GetAll()
	sort.Sort(transactionList(cc))
	return cc
}
