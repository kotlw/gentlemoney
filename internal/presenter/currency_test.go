package presenter_test

import (
	"testing"

	"github.com/kotlw/gentlemoney/internal/model"
	"github.com/kotlw/gentlemoney/internal/presenter"

	_ "github.com/mattn/go-sqlite3"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type CurrencyPresenterTestSuite struct {
	suite.Suite
	presenter *presenter.Currency
}

func (s *CurrencyPresenterTestSuite) SetupSuite() {
	s.presenter = presenter.NewCurrency()
}

func (s *CurrencyPresenterTestSuite) TestToMap() {
	currency := &model.Currency{Abbreviation: "USD"}
	expected := map[string]string{"ID": "0", "Abbreviation": "USD"}
	actual := s.presenter.ToMap(currency)
	assert.Equal(s.T(), expected, actual)
}

func (s *CurrencyPresenterTestSuite) TestFromMapPositive() {
	for _, tc := range []struct {
		name     string
		give     map[string]string
		expected *model.Currency
	}{
		{
			name:     "ExistingID",
			give:     map[string]string{"ID": "99", "Abbreviation": "USD"},
			expected: &model.Currency{ID: int64(99), Abbreviation: "USD"},
		},
		{
			name:     "NotExistingID",
			give:     map[string]string{"Abbreviation": "USD"},
			expected: &model.Currency{ID: int64(0), Abbreviation: "USD"},
		},
	} {
		s.Run(tc.name, func() {
			actual, err := s.presenter.FromMap(tc.give)
			require.NoError(s.T(), err)
			assert.Equal(s.T(), tc.expected, actual)
		})
	}
}

func (s *CurrencyPresenterTestSuite) TestFromMapNegative() {
	_, err := s.presenter.FromMap(map[string]string{"Name": "None"})
	assert.EqualError(s.T(), err, `checkKeys: key "Abbreviation" is missing`)
}

func TestCurrencyPresenterTestSuite(t *testing.T) {
	suite.Run(t, new(CurrencyPresenterTestSuite))
}
