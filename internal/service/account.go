package service

import (
	"fmt"

	"github.com/kotlw/gentlemoney/internal/model"
	"github.com/kotlw/gentlemoney/internal/storage/inmemory"
	"github.com/kotlw/gentlemoney/internal/storage/sqlite"
)

// Account service contains business logic related to model.Account.
type Account struct {
	persistentStorage *sqlite.Account
	inmemoryStorage   *inmemory.Account
}

// NewAccount returns Account service.
func NewAccount(
	persistentStorage *sqlite.Account,
	inmemoryStorage *inmemory.Account,
	currencyService *Currency) (*Account, error) {

	a := &Account{
		persistentStorage: persistentStorage,
		inmemoryStorage:   inmemoryStorage,
	}

	if err := a.Init(currencyService); err != nil {
		return nil, fmt.Errorf("c.Init: %w", err)
	}

	return a, nil
}

// Init initialize inmemory storage with data from persistent storage. It is also links existing
// currencies to corresponding fields of model.Account.
func (s *Account) Init(currencyService *Currency) error {
	aa, err := s.persistentStorage.GetAll()
	if err != nil {
		return fmt.Errorf("s.persistentStorage.GetAll: %w", err)
	}

	for _, a := range aa {
		a.Currency = currencyService.GetByID(a.Currency.ID)
	}

	s.inmemoryStorage.Init(aa)

	return nil
}

// Insert appends account to both persistent and inmemory storages.
func (s *Account) Insert(a *model.Account) error {
	id, err := s.persistentStorage.Insert(a)
	if err != nil {
		return fmt.Errorf("s.persistentStorage.Insert: %w", err)
	}

	a.ID = id
	s.inmemoryStorage.Insert(a)

	return nil
}

// Update updates account in persistent storage. Since GetAll returns pointers to inmemory data
// after update the category we need to update it in persistent storage as well.
func (s *Account) Update(a *model.Account) error {
	if err := s.persistentStorage.Update(a); err != nil {
		return fmt.Errorf("s.persistentStorage.Update: %w", err)
	}

	s.inmemoryStorage.Update(a)

	return nil
}

// Delete deletes account from inmemory and persistent storages.
func (s *Account) Delete(a *model.Account) error {
	if err := s.persistentStorage.Delete(a.ID); err != nil {
		return fmt.Errorf("s.persistentStorage.Delete: %w", err)
	}

	s.inmemoryStorage.Delete(a)

	return nil
}

// GetAll returns all accounts.
func (s *Account) GetAll() []*model.Account {
	return s.inmemoryStorage.GetAll()
}

// GetByID returns account by given model.Account.ID.
func (s *Account) GetByID(id int64) *model.Account {
	return s.inmemoryStorage.GetByID(id)
}

// GetByName returns currency by given model.Account.Name.
func (s *Account) GetByName(name string) *model.Account {
	return s.inmemoryStorage.GetByName(name)
}
