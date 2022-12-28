package service

import (
	"fmt"
	"sort"

	"gentlemoney/internal/model"
	"gentlemoney/internal/storage/inmemory"
	"gentlemoney/internal/storage/sqlite"
)

// currencyList is a wrapper to perform sort by model.Currency.Abbreviation.
type currencyList []*model.Currency

func (cc currencyList) Len() int           { return len(cc) }
func (cc currencyList) Less(i, j int) bool { return cc[i].Abbreviation < cc[j].Abbreviation }
func (cc currencyList) Swap(i, j int)      { cc[i], cc[j] = cc[j], cc[i] }

// Currency service contains business logic related to model.Currency.
type Currency struct {
	persistantStorage *sqlite.Currency
	inmemoryStorage   *inmemory.Currency
}

// NewCurrency returns Currency service.
func NewCurrency(persistantStorage *sqlite.Currency, inmemoryStorage *inmemory.Currency) (*Currency, error) {
    c := &Currency{
        persistantStorage: persistantStorage,
        inmemoryStorage: inmemoryStorage,
    }

    if err := c.Init(); err != nil {
        return nil, fmt.Errorf("c.Init: %w", err)
    }

    return c, nil
}

// Init initialize inmemory storage with data from persistant storage.
func (s *Currency) Init() error {
    cc, err := s.persistantStorage.GetAll()
    if err != nil {
        return fmt.Errorf("s.persistantStorage.GetAll: %w", err)
    }

    s.inmemoryStorage.Init(cc)

    return nil
}

// Insert appends currency to both persistant and inmemory storages.
func (s *Currency) Insert(c *model.Currency) error {
    id, err := s.persistantStorage.Insert(c)
    if err != nil {
        return fmt.Errorf("s.persistantStorage.Insert: %w", err)
    }
    
    c.ID = id
    s.inmemoryStorage.Insert(c)

    return nil
}

// Update updates currency in persistant storage. Since GetAll returns pointers to inmemory data
// after update the category we need to update it in persistant storage as well.
func (s *Currency) Update(c *model.Currency) error {
    if err := s.persistantStorage.Update(c); err != nil {
        return fmt.Errorf("s.persistantStorage.Update: %w", err)
    }
    return nil
}

// Delete deletes currency from inmemory and persistant storages.
func (s *Currency) Delete(c *model.Currency) error {
    if err := s.persistantStorage.Delete(c.ID); err != nil {
        return fmt.Errorf("s.persistantStorage.Delete: %w", err)
    }
    s.inmemoryStorage.Delete(c)
    return nil
}

// GetAll returns all currencies.
func (s *Currency) GetAll() []*model.Currency {
    return s.inmemoryStorage.GetAll()
}

// GetByID returns currency by given model.Currency.ID.
func (s *Currency) GetByID(id int64) *model.Currency {
    return s.inmemoryStorage.GetByID(id)
}

// GetByAbbreviation returns currency by given model.Currency.Abbreviation.
func (s *Currency) GetByAbbreviation(abbreviation string) *model.Currency {
    return s.inmemoryStorage.GetByAbbreviation(abbreviation)
}

// GetAllSorted returns all currencies sorted by model.Currency.Abbreviation.
func (s *Currency) GetAllSorted() []*model.Currency {
    cc := s.inmemoryStorage.GetAll()
    sort.Sort(currencyList(cc))
    return cc
}
