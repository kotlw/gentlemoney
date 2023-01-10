package service

import (
	"fmt"

	"github.com/kotlw/gentlemoney/internal/model"
	"github.com/kotlw/gentlemoney/internal/storage/inmemory"
	"github.com/kotlw/gentlemoney/internal/storage/sqlite"
)

// Currency service contains business logic related to model.Currency.
type Currency struct {
	persistentStorage *sqlite.Currency
	inmemoryStorage   *inmemory.Currency
}

// NewCurrency returns Currency service.
func NewCurrency(persistentStorage *sqlite.Currency, inmemoryStorage *inmemory.Currency) (*Currency, error) {
	c := &Currency{
		persistentStorage: persistentStorage,
		inmemoryStorage:   inmemoryStorage,
	}

	if err := c.Init(); err != nil {
		return nil, fmt.Errorf("c.Init: %w", err)
	}

	return c, nil
}

// Init initialize inmemory storage with data from persistent storage.
func (s *Currency) Init() error {
	cc, err := s.persistentStorage.GetAll()
	if err != nil {
		return fmt.Errorf("s.persistentStorage.GetAll: %w", err)
	}

	s.inmemoryStorage.Init(cc)

	return nil
}

// Insert appends currency to both persistent and inmemory storages.
func (s *Currency) Insert(c *model.Currency) error {
	id, err := s.persistentStorage.Insert(c)
	if err != nil {
		return fmt.Errorf("s.persistentStorage.Insert: %w", err)
	}

	c.ID = id
	s.inmemoryStorage.Insert(c)

	return nil
}

// Update updates currency in persistent storage. Since GetAll returns pointers to inmemory data
// after update the category we need to update it in persistent storage as well.
func (s *Currency) Update(c *model.Currency) error {
	if err := s.persistentStorage.Update(c); err != nil {
		return fmt.Errorf("s.persistentStorage.Update: %w", err)
	}

	s.inmemoryStorage.Update(c)

	return nil
}

// Delete deletes currency from inmemory and persistent storages.
func (s *Currency) Delete(c *model.Currency) error {
	if err := s.persistentStorage.Delete(c.ID); err != nil {
		return fmt.Errorf("s.persistentStorage.Delete: %w", err)
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
