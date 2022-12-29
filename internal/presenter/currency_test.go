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
	for _, tc := range []struct {
		name     string
		give     *model.Currency
		expected map[string]string
	}{
		{
			name:     "Zero",
			give:     &model.Currency{Abbreviation: "USD", ExchangeRate: 0, IsMain: true},
			expected: map[string]string{"Abbreviation": "USD", "ExchangeRate": "0.00", "IsMain": "*"},
		},
		{
			name:     "PositiveOneDigit",
			give:     &model.Currency{Abbreviation: "USD", ExchangeRate: 9, IsMain: false},
			expected: map[string]string{"Abbreviation": "USD", "ExchangeRate": "0.09", "IsMain": ""},
		},
		{
			name:     "PositiveTwoDigits",
			give:     &model.Currency{Abbreviation: "USD", ExchangeRate: 91, IsMain: true},
			expected: map[string]string{"Abbreviation": "USD", "ExchangeRate": "0.91", "IsMain": "*"},
		},
		{
			name:     "PositiveLongNumber",
			give:     &model.Currency{Abbreviation: "USD", ExchangeRate: 911321341243124141, IsMain: true},
			expected: map[string]string{"Abbreviation": "USD", "ExchangeRate": "9113213412431241.41", "IsMain": "*"},
		},
		{
			name:     "NegativeOneDigit",
			give:     &model.Currency{Abbreviation: "USD", ExchangeRate: -1, IsMain: true},
			expected: map[string]string{"Abbreviation": "USD", "ExchangeRate": "-0.01", "IsMain": "*"},
		},
		{
			name:     "NegativeTwoDigits",
			give:     &model.Currency{Abbreviation: "USD", ExchangeRate: -19, IsMain: true},
			expected: map[string]string{"Abbreviation": "USD", "ExchangeRate": "-0.19", "IsMain": "*"},
		},
	} {
		s.Run(tc.name, func() {
			actual := s.presenter.ToMap(tc.give)
			assert.Equal(s.T(), tc.expected, actual)
		})
	}
}

func (s *CurrencyPresenterTestSuite) TestFromMapPositive() {
	for _, tc := range []struct {
		name     string
		give     map[string]string
		expected *model.Currency
	}{
		{
			name:     "Zero",
			give:     map[string]string{"Abbreviation": "USD", "ExchangeRate": "0.00", "IsMain": "*"},
			expected: &model.Currency{Abbreviation: "USD", ExchangeRate: 0, IsMain: true},
		},
		{
			name:     "PositiveOneDigit",
			give:     map[string]string{"Abbreviation": "USD", "ExchangeRate": "0.09", "IsMain": ""},
			expected: &model.Currency{Abbreviation: "USD", ExchangeRate: 9, IsMain: false},
		},
		{
			name:     "PositiveTwoDigits",
			give:     map[string]string{"Abbreviation": "USD", "ExchangeRate": "0.91", "IsMain": "*"},
			expected: &model.Currency{Abbreviation: "USD", ExchangeRate: 91, IsMain: true},
		},
		{
			name:     "PositiveLongNumber",
			give:     map[string]string{"Abbreviation": "USD", "ExchangeRate": "9113213412431241.41", "IsMain": "*"},
			expected: &model.Currency{Abbreviation: "USD", ExchangeRate: 911321341243124141, IsMain: true},
		},
		{
			name:     "NegativeOneDigit",
			give:     map[string]string{"Abbreviation": "USD", "ExchangeRate": "-0.01", "IsMain": "*"},
			expected: &model.Currency{Abbreviation: "USD", ExchangeRate: -1, IsMain: true},
		},
		{
			name:     "NegativeTwoDigits",
			give:     map[string]string{"Abbreviation": "USD", "ExchangeRate": "-0.19", "IsMain": "*"},
			expected: &model.Currency{Abbreviation: "USD", ExchangeRate: -19, IsMain: true},
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
	for _, tc := range []struct {
		name     string
		give     map[string]string
		expected string
	}{
		{
			name:     "InvalidExchangeRate",
			give:     map[string]string{"Abbreviation": "USD", "ExchangeRate": "invalid", "IsMain": "*"},
			expected: `parseMoney: strconv.Atoi: strconv.Atoi: parsing "invalid": invalid syntax`,
		},
		{
			name:     "MissingAbbreviation",
			give:     map[string]string{"ExchangeRate": "0.99", "IsMain": "*"},
			expected: `checkKeys: key "Abbreviation" is missing`,
		},
		{
			name:     "MissingExchangeRate",
			give:     map[string]string{"Abbreviation": "USD", "IsMain": "*"},
			expected: `checkKeys: key "ExchangeRate" is missing`,
		},
		{
			name:     "MissingIsMain",
			give:     map[string]string{"Abbreviation": "USD", "ExchangeRate": "invalid"},
			expected: `checkKeys: key "IsMain" is missing`,
		},
	} {
		s.Run(tc.name, func() {
			_, err := s.presenter.FromMap(tc.give)
			assert.EqualError(s.T(), err, tc.expected)
		})
	}
}

func TestCurrencyPresenterTestSuite(t *testing.T) {
	suite.Run(t, new(CurrencyPresenterTestSuite))
}
