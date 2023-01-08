package presenter_test

import (
	"database/sql"
	"testing"

	"github.com/kotlw/gentlemoney/internal/model"
	"github.com/kotlw/gentlemoney/internal/presenter"
	"github.com/kotlw/gentlemoney/internal/service"
	"github.com/kotlw/gentlemoney/internal/storage/inmemory"
	"github.com/kotlw/gentlemoney/internal/storage/sqlite"

	_ "github.com/mattn/go-sqlite3"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type AccountPresenterTestSuite struct {
	suite.Suite
	db           *sql.DB
	presenter    *presenter.Account
	initCurrency *model.Currency
}

func (s *AccountPresenterTestSuite) SetupSuite() {
	db, err := sql.Open("sqlite3", "file::memory:?cache=shared")
	require.NoError(s.T(), err, "occurred in SetupSuite")
	s.db = db

	currencyPersistentStorage, err := sqlite.NewCurrency(db)
	require.NoError(s.T(), err, "occurred in SetupSuite")

	currencyService, err := service.NewCurrency(currencyPersistentStorage, inmemory.NewCurrency())
	require.NoError(s.T(), err, "occurred in SetupSuite")
	s.presenter = presenter.NewAccount(currencyService)

	s.initCurrency = &model.Currency{Abbreviation: "USD", ExchangeRate: 100, IsMain: true}
	err = currencyService.Insert(s.initCurrency)
	require.NoError(s.T(), err, "occurred in SetupSuite")
}

func (s *AccountPresenterTestSuite) TestToMap() {
	account := &model.Account{Name: "Card1", Currency: s.initCurrency}
	expected := map[string]string{"ID": "0", "Name": "Card1", "Currency": s.initCurrency.Abbreviation}
	actual := s.presenter.ToMap(account)
	assert.Equal(s.T(), expected, actual)
}

func (s *AccountPresenterTestSuite) TestFromMapPositive() {
	for _, tc := range []struct {
		name     string
		give     map[string]string
		expected *model.Account
	}{
		{
			name:     "ExistingID",
			give:     map[string]string{"ID": "21", "Name": "Card1", "Currency": "USD"},
			expected: &model.Account{ID: int64(21), Name: "Card1", Currency: s.initCurrency},
		},
		{
			name:     "NotExistingID",
			give:     map[string]string{"Name": "Card1", "Currency": "USD"},
			expected: &model.Account{ID: int64(0), Name: "Card1", Currency: s.initCurrency},
		},
		{
			name:     "ExistingCurrency",
			give:     map[string]string{"Name": "Card1", "Currency": "USD"},
			expected: &model.Account{Name: "Card1", Currency: s.initCurrency},
		},
		{
			name:     "NotExistingCurrency",
			give:     map[string]string{"Name": "Card1", "Currency": "???"},
			expected: &model.Account{Name: "Card1", Currency: nil},
		},
	} {
		s.Run(tc.name, func() {
			actual, err := s.presenter.FromMap(tc.give)
			require.NoError(s.T(), err)
			assert.Equal(s.T(), tc.expected, actual)
		})
	}

}

func (s *AccountPresenterTestSuite) TestFromMapNegative() {
	for _, tc := range []struct {
		name     string
		give     map[string]string
		expected string
	}{
		{
			name:     "MissingName",
			give:     map[string]string{"Currency": "USD"},
			expected: `checkKeys: key "Name" is missing`,
		},
		{
			name:     "MissingCurrency",
			give:     map[string]string{"Name": "Card1"},
			expected: `checkKeys: key "Currency" is missing`,
		},
	} {
		s.Run(tc.name, func() {
			_, err := s.presenter.FromMap(tc.give)
			assert.EqualError(s.T(), err, tc.expected)
		})
	}
}

func (s *AccountPresenterTestSuite) TearDownSuite() {
	err := s.db.Close()
	require.NoError(s.T(), err, "occurred in TearDownSuite")
}

func TestAccountPresenterTestSuite(t *testing.T) {
	suite.Run(t, new(AccountPresenterTestSuite))
}
