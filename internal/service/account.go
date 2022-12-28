package service

import (
	"fmt"
	"sort"

	"github.com/kotlw/gentlemoney/internal/model"
	"github.com/kotlw/gentlemoney/internal/storage/inmemory"
	"github.com/kotlw/gentlemoney/internal/storage/sqlite"
)

// accountList is a wrapper to perform sort by model.Account.Name.
type accountList []*model.Account

func (aa accountList) Len() int           { return len(aa) }
func (aa accountList) Less(i, j int) bool { return aa[i].Name < aa[j].Name }
func (aa accountList) Swap(i, j int)      { aa[i], aa[j] = aa[j], aa[i] }

// Account service contains business logic related to model.Account.
type Account struct {
	persistantStorage *sqlite.Account
	inmemoryStorage   *inmemory.Account
}

// NewAccount returns Account service.
func NewAccount(
	persistantStorage *sqlite.Account,
	inmemoryStorage *inmemory.Account,
	currencyService *Currency) (*Account, error) {

	a := &Account{
		persistantStorage: persistantStorage,
		inmemoryStorage:   inmemoryStorage,
	}

	if err := a.Init(currencyService); err != nil {
		return nil, fmt.Errorf("c.Init: %w", err)
	}

	return a, nil
}

// Init initialize inmemory storage with data from persistant storage. It is also links existing
// currencies to corresponding fields of model.Account.
func (s *Account) Init(currencyService *Currency) error {
	aa, err := s.persistantStorage.GetAll()
	if err != nil {
		return fmt.Errorf("s.persistantStorage.GetAll: %w", err)
	}

    for _, a := range aa {
        a.Currency = currencyService.GetByID(a.Currency.ID)
    }

	s.inmemoryStorage.Init(aa)

	return nil
}

// Insert appends account to both persistant and inmemory storages.
func (s *Account) Insert(c *model.Account) error {
	id, err := s.persistantStorage.Insert(c)
	if err != nil {
		return fmt.Errorf("s.persistantStorage.Insert: %w", err)
	}

	c.ID = id
	s.inmemoryStorage.Insert(c)

	return nil
}

// Update updates account in persistant storage. Since GetAll returns pointers to inmemory data
// after update the category we need to update it in persistant storage as well.
func (s *Account) Update(c *model.Account) error {
	if err := s.persistantStorage.Update(c); err != nil {
		return fmt.Errorf("s.persistantStorage.Update: %w", err)
	}
	return nil
}

// Delete deletes account from inmemory and persistant storages.
func (s *Account) Delete(c *model.Account) error {
	if err := s.persistantStorage.Delete(c.ID); err != nil {
		return fmt.Errorf("s.persistantStorage.Delete: %w", err)
	}
	s.inmemoryStorage.Delete(c)
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

// GetAllSorted returns all accounts sorted by model.Account.Name.
func (s *Account) GetAllSorted() []*model.Account {
	cc := s.inmemoryStorage.GetAll()
	sort.Sort(accountList(cc))
	return cc
}
