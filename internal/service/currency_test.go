package service_test

import (
	"database/sql"
	"sort"
	"testing"

	"github.com/kotlw/gentlemoney/internal/model"
	"github.com/kotlw/gentlemoney/internal/service"
	"github.com/kotlw/gentlemoney/internal/storage/inmemory"
	"github.com/kotlw/gentlemoney/internal/storage/sqlite"

	_ "github.com/mattn/go-sqlite3"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type currencyList []*model.Currency

func (cc currencyList) Len() int           { return len(cc) }
func (cc currencyList) Less(i, j int) bool { return cc[i].Abbreviation < cc[j].Abbreviation }
func (cc currencyList) Swap(i, j int)      { cc[i], cc[j] = cc[j], cc[i] }

type CurrencyServiceTestSuite struct {
	suite.Suite
	db                *sql.DB
	persistentStorage *sqlite.Currency
	inmemoryStorage   *inmemory.Currency
	service           *service.Currency
	InitCurrencies    []*model.Currency
}

func (s *CurrencyServiceTestSuite) SetupSuite() {
	db, err := sql.Open("sqlite3", "file::memory:?cache=shared")
	require.NoError(s.T(), err, "occurred in SetupSuite")
	s.db = db

	s.persistentStorage, err = sqlite.NewCurrency(db)
	require.NoError(s.T(), err, "occurred in SetupSuite")

	s.inmemoryStorage = inmemory.NewCurrency()

	s.service, err = service.NewCurrency(s.persistentStorage, s.inmemoryStorage)
	require.NoError(s.T(), err, "occurred in SetupSuite")

	// id's settled by sqlite on insert incrementally starting from 1,
	// so here they are initialized for match purpose
	s.InitCurrencies = []*model.Currency{
		{ID: 1, Abbreviation: "USD", ExchangeRate: 100, IsMain: true},
		{ID: 2, Abbreviation: "EUR", ExchangeRate: 124, IsMain: false},
	}
}

func (s *CurrencyServiceTestSuite) SetupTest() {
	for _, c := range s.InitCurrencies {
		_, err := s.persistentStorage.Insert(c)
		require.NoError(s.T(), err, "occurred in SetupTest")
	}

	err := s.service.Init()
	require.NoError(s.T(), err, "occurred in SetupTest")
}

func (s *CurrencyServiceTestSuite) TestInsertPositive() {
	currency := &model.Currency{ID: 3, Abbreviation: "PLN", ExchangeRate: 200, IsMain: false}
	expectedCurrencies := append(s.InitCurrencies, currency)

	err := s.service.Insert(currency)
	require.NoError(s.T(), err)

	persistentCurrencies, err := s.persistentStorage.GetAll()
	require.NoError(s.T(), err)
	inmemoryCurrencies := s.inmemoryStorage.GetAll()
	assert.ElementsMatch(s.T(), persistentCurrencies, expectedCurrencies)
	assert.ElementsMatch(s.T(), inmemoryCurrencies, expectedCurrencies)
}

func (s *CurrencyServiceTestSuite) TestInsertNegative() {
	err := s.service.Insert(s.InitCurrencies[0])
	assert.ErrorContains(s.T(), err, "s.persistentStorage.Insert: e.db.Exec: UNIQUE constraint failed: currency.abbreviation")
}

func (s *CurrencyServiceTestSuite) TestUpdatePositive() {
	cc := s.service.GetAll()
	cc[0].Abbreviation = "CZN"

	err := s.service.Update(cc[0])
	require.NoError(s.T(), err)

	persistentCurrencies, err := s.persistentStorage.GetAll()
	inmemoryCurrencies := s.inmemoryStorage.GetAll()
	require.NoError(s.T(), err)
	assert.ElementsMatch(s.T(), persistentCurrencies, inmemoryCurrencies)
}

func (s *CurrencyServiceTestSuite) TestUpdateNegative() {
	cc := s.service.GetAll()
	cc[0].Abbreviation = "CZN"
	cc[0].ID = 10

	err := s.service.Update(cc[0])
	assert.ErrorContains(s.T(), err, "s.persistentStorage.Update: total affected rows 0 while expected 1")
	cc[0].ID = 1 // return real id to proper teardown
}

func (s *CurrencyServiceTestSuite) TestDeletePositive() {
	cc := s.service.GetAll()
	expectedCurrencies := []*model.Currency{cc[0]}

	err := s.service.Delete(cc[1])
	require.NoError(s.T(), err)

	persistentCurrencies, err := s.persistentStorage.GetAll()
	require.NoError(s.T(), err)
	inmemoryCurrencies := s.inmemoryStorage.GetAll()
	assert.ElementsMatch(s.T(), persistentCurrencies, expectedCurrencies)
	assert.ElementsMatch(s.T(), inmemoryCurrencies, expectedCurrencies)
}

func (s *CurrencyServiceTestSuite) TestDeleteNegative() {
	cc := s.service.GetAll()
	cc[0].ID = 10

	err := s.service.Delete(cc[0])
	assert.ErrorContains(s.T(), err, "s.persistentStorage.Delete: total affected rows 0 while expected 1")
	cc[0].ID = 1 // return real id to proper teardown
}

func (s *CurrencyServiceTestSuite) TestGetByID() {
	c := s.service.GetByID(2)
	assert.Equal(s.T(), s.InitCurrencies[1].Abbreviation, c.Abbreviation)
}

func (s *CurrencyServiceTestSuite) TestGetByAbbreviation() {
	c := s.service.GetByAbbreviation(s.InitCurrencies[1].Abbreviation)
	assert.Equal(s.T(), s.InitCurrencies[1].ID, c.ID)
}

func (s *CurrencyServiceTestSuite) TestGetAllSorted() {
	cc := s.service.GetAllSorted()
	expectedCurrencies := make([]*model.Currency, len(s.InitCurrencies))
	copy(expectedCurrencies, s.InitCurrencies)
	sort.Sort(currencyList(expectedCurrencies))
	assert.ElementsMatch(s.T(), cc, expectedCurrencies)
}

func (s *CurrencyServiceTestSuite) TearDownTest() {
	for {
		cc := s.service.GetAll()
		if len(cc) == 0 {
			break
		}

		err := s.persistentStorage.Delete(cc[0].ID)
		require.NoError(s.T(), err, "occurred in TearDownTest")
		s.inmemoryStorage.Delete(cc[0])
	}
}

func (s *CurrencyServiceTestSuite) TearDownSuite() {
	err := s.db.Close()
	require.NoError(s.T(), err, "occurred in TearDownSuite")
}

func TestCurrencyServiceTestSuite(t *testing.T) {
	suite.Run(t, new(CurrencyServiceTestSuite))
}
