package inmemory

import (
	"gentlemoney/internal/model"
)

// Currency is used to acces inmemory storage.
type Currency struct {
	currencies     []*model.Currency
	currencyByID   map[int64]*model.Currency
	currencyByAbbr map[string]*model.Currency
}

// NewCurrency returns new currency inmemory storage.
func NewCurrency() *Currency {
	return &Currency{
		currencies:     make([]*model.Currency, 0, 20),
		currencyByID:   make(map[int64]*model.Currency),
		currencyByAbbr: make(map[string]*model.Currency),
	}
}

// Init initialize inmemory storage with given slice of data.
func (s *Currency) Init(cc []*model.Currency) {
	for _, c := range cc {
		s.currencyByID[c.ID] = c
		s.currencyByAbbr[c.Abbreviation] = c
	}
	s.currencies = cc
}

// Insert appends currency to inmemory storage.
func (s *Currency) Insert(c *model.Currency) {
	s.currencyByID[c.ID] = c
	s.currencyByAbbr[c.Abbreviation] = c
	s.currencies = append(s.currencies, c)
}

// Delete removes currency from current inmemory storage.
func (s *Currency) Delete(c *model.Currency) {
	delete(s.currencyByID, c.ID)
	delete(s.currencyByAbbr, c.Abbreviation)

	if len(s.currencies) == 1 {
		s.currencies = make([]*model.Currency, 0, 20)
	} else {
        i := findIndex(c, s.currencies)
		s.currencies = append(s.currencies[:i], s.currencies[i+1:]...)
	}
}

// GetAll returns slice of currencies.
func (s *Currency) GetAll() []*model.Currency {
	return s.currencies
}

// GetByID returns currency by its id.
func (s *Currency) GetByID(id int64) *model.Currency {
	return s.currencyByID[id]
}

// GetByAbbreviation returns currency by its abbreviation.
func (s *Currency) GetByAbbreviation(abbreviation string) *model.Currency {
	return s.currencyByAbbr[abbreviation]
}
