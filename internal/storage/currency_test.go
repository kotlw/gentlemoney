package storage_test

import (
	"database/sql"
	"testing"

	"gentlemoney/internal/model"
	"gentlemoney/internal/storage"

	_ "github.com/mattn/go-sqlite3"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type CurrencyStorageTestSuite struct {
	suite.Suite
	db             *sql.DB
	storage        *storage.Currency
	InitCurrencies []*model.Currency
}

func (s *CurrencyStorageTestSuite) SetupSuite() {
	db, err := sql.Open("sqlite3", "file::memory:?cache=shared")
	require.NoError(s.T(), err, "occured in SetupSuite")
	s.db = db

	s.storage, err = storage.NewCurrency(db)
	require.NoError(s.T(), err, "occured in SetupSuite")

	// id's settled by sqlite on insert incrementally starting from 1,
	// so here they are initialized for match purpose
	s.InitCurrencies = []*model.Currency{
		{ID: 1, Abbreviation: "USD", ExchangeRate: 100, IsMain: true},
		{ID: 2, Abbreviation: "EUR", ExchangeRate: 124, IsMain: false},
	}
}

func (s *CurrencyStorageTestSuite) SetupTest() {
	stmt, err := s.db.Prepare(`INSERT INTO currency(abbreviation, exchangeRate, isMain) VALUES (?, ?, ?);`)
	require.NoError(s.T(), err, "occured in SetupTest")

	for _, currency := range s.InitCurrencies {
		_, err := stmt.Exec(currency.Abbreviation, currency.ExchangeRate, currency.IsMain)
		require.NoError(s.T(), err, "occured in SetupTest")
	}
}

func (s *CurrencyStorageTestSuite) TestInsertPositive() {
	currency := &model.Currency{ID: 3, Abbreviation: "PLN", ExchangeRate: 200, IsMain: false}
	expectedCurrencies := append(s.InitCurrencies, currency)

	_, err := s.storage.Insert(currency)
	require.NoError(s.T(), err)
	assert.ElementsMatch(s.T(), s.fetchActualData(), expectedCurrencies)
}

func (s *CurrencyStorageTestSuite) TestInsertNegative() {
	_, err := s.storage.Insert(s.InitCurrencies[1])
	assert.EqualError(s.T(), err, "e.db.Exec: UNIQUE constraint failed: currency.abbreviation")
}

func (s *CurrencyStorageTestSuite) TestUpdatePositive() {
	expectedCurrencies := make([]*model.Currency, len(s.InitCurrencies))
	copy(expectedCurrencies, s.InitCurrencies)
	expectedCurrencies[1].Abbreviation = "CZN"

	err := s.storage.Update(expectedCurrencies[1])
	require.NoError(s.T(), err)
	assert.ElementsMatch(s.T(), s.fetchActualData(), expectedCurrencies)
}

func (s *CurrencyStorageTestSuite) TestUpdateNegative() {
	err := s.storage.Update(&model.Currency{ID: 10})
	assert.EqualError(s.T(), err, "total affected rows 0 while expected 1")
}

func (s *CurrencyStorageTestSuite) TestDeletePositive() {
	expectedCurrencies := []*model.Currency{s.InitCurrencies[0]}

	err := s.storage.Delete(2)
	require.NoError(s.T(), err)
	assert.ElementsMatch(s.T(), s.fetchActualData(), expectedCurrencies)
}

func (s *CurrencyStorageTestSuite) TestDeleteNegative() {
	err := s.storage.Delete(10)
	assert.EqualError(s.T(), err, "total affected rows 0 while expected 1")
}

func (s *CurrencyStorageTestSuite) TestGetAll() {
	allCurrencies, err := s.storage.GetAll()
	require.NoError(s.T(), err)
	assert.Equal(s.T(), allCurrencies, s.InitCurrencies)
}

func (s *CurrencyStorageTestSuite) fetchActualData() []*model.Currency {
	rows, err := s.db.Query(`SELECT id, abbreviation, exchangeRate, isMain FROM currency;`)
	require.NoError(s.T(), err)
	defer func() {
		err = rows.Close()
		require.NoError(s.T(), err)
	}()

	res := make([]*model.Currency, 0, 3)
	for rows.Next() {
		t := model.NewEmptyCurrency()
		err = rows.Scan(&t.ID, &t.Abbreviation, &t.ExchangeRate, &t.IsMain)
		require.NoError(s.T(), err)
		res = append(res, t)
	}

	return res
}

func (s *CurrencyStorageTestSuite) TearDownTest() {
	stmt, err := s.db.Prepare(`DELETE FROM currency;`)
	require.NoError(s.T(), err, "occured in TearDownTest")

	_, err = stmt.Exec()
	require.NoError(s.T(), err, "occured in TearDownTest")
}

func (s *CurrencyStorageTestSuite) TearDownSuite() {
	err := s.db.Close()
	require.NoError(s.T(), err, "occured in TearDownSuite")
}

func TestCurrencyStorageTestSuite(t *testing.T) {
	suite.Run(t, new(CurrencyStorageTestSuite))
}
