package inmemory_test

import (
	"testing"

	"github.com/kotlw/gentlemoney/internal/model"
	"github.com/kotlw/gentlemoney/internal/storage/inmemory"

	_ "github.com/mattn/go-sqlite3"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type CurrencyInmemoryStorageTestSuite struct {
	suite.Suite
	storage        *inmemory.Currency
	InitCurrencies []*model.Currency
}

func (s *CurrencyInmemoryStorageTestSuite) SetupSuite() {
	s.storage = inmemory.NewCurrency()
	s.InitCurrencies = []*model.Currency{
		{ID: 1, Abbreviation: "USD", ExchangeRate: 100, IsMain: true},
		{ID: 2, Abbreviation: "EUR", ExchangeRate: 124, IsMain: false},
	}
}

func (s *CurrencyInmemoryStorageTestSuite) SetupTest() {
	s.storage.Init(s.InitCurrencies)
}

func (s *CurrencyInmemoryStorageTestSuite) TestInsertPositive() {
	currency := &model.Currency{ID: 3, Abbreviation: "PLN", ExchangeRate: 200, IsMain: false}
	expectedCurrencies := append(s.InitCurrencies, currency)

	s.storage.Insert(currency)

	assert.ElementsMatch(s.T(), s.storage.GetAll(), expectedCurrencies)
	assert.Equal(s.T(), currency, s.storage.GetByID(currency.ID))
	assert.Equal(s.T(), currency, s.storage.GetByAbbreviation(currency.Abbreviation))
}

func (s *CurrencyInmemoryStorageTestSuite) TestDeletePositive() {
	s.storage.Delete(s.InitCurrencies[1])

	assert.ElementsMatch(s.T(), s.storage.GetAll(), s.InitCurrencies[:1])
}

func (s *CurrencyInmemoryStorageTestSuite) TearDownTest() {
	for {
		cc := s.storage.GetAll()
		if len(cc) == 0 {
			break
		}
		s.storage.Delete(cc[0])
	}
}

func TestCurrencyInmemoryStorageTestSuite(t *testing.T) {
	suite.Run(t, new(CurrencyInmemoryStorageTestSuite))
}
