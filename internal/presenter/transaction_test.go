package presenter_test

import (
	"database/sql"
	"testing"
	"time"

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

type TransactionPresenterTestSuite struct {
	suite.Suite
	db           *sql.DB
	presenter    *presenter.Transaction
	initCategory *model.Category
	initCurrency *model.Currency
	initAccount  *model.Account
}

func (s *TransactionPresenterTestSuite) SetupSuite() {
	db, err := sql.Open("sqlite3", "file::memory:?cache=shared")
	require.NoError(s.T(), err, "occurred in SetupSuite")
	s.db = db

	persistentStorage, err := sqlite.New(db)
	require.NoError(s.T(), err, "occurred in SetupSuite")

	inmemoryStorage := inmemory.New()

	service, err := service.New(persistentStorage, inmemoryStorage)
	require.NoError(s.T(), err, "occurred in SetupSuite")

	s.presenter = presenter.NewTransaction(service.Account(), service.Category())

	s.initCategory = &model.Category{Title: "Health"}
	err = service.Category().Insert(s.initCategory)
	require.NoError(s.T(), err, "occurred in SetupSuite")

	s.initCurrency = &model.Currency{Abbreviation: "USD", ExchangeRate: 100, IsMain: true}
	err = service.Currency().Insert(s.initCurrency)
	require.NoError(s.T(), err, "occurred in SetupSuite")

	s.initAccount = &model.Account{Name: "Card1", Currency: s.initCurrency}
	err = service.Account().Insert(s.initAccount)
	require.NoError(s.T(), err, "occurred in SetupSuite")
}

func (s *TransactionPresenterTestSuite) TestToMap() {
	for _, tc := range []struct {
		name     string
		give     *model.Transaction
		expected map[string]string
	}{
		{
			name: "Zero",
			give: &model.Transaction{
				Date:     time.Date(2020, 5, 6, 11, 45, 04, 0, time.UTC),
				Account:  s.initAccount,
				Category: s.initCategory,
				Amount:   0,
				Note:     "Note1",
			},
			expected: map[string]string{
				"Date":     "2020-05-06",
				"Account":  s.initAccount.Name,
				"Category": s.initCategory.Title,
				"Currency": s.initCurrency.Abbreviation,
				"Amount":   "0.00",
				"Note":     "Note1",
			},
		},
		{
			name: "PositiveOneDigit",
			give: &model.Transaction{
				Date:     time.Date(2020, 5, 6, 11, 45, 04, 0, time.UTC),
				Account:  s.initAccount,
				Category: s.initCategory,
				Amount:   1,
				Note:     "Note1",
			},
			expected: map[string]string{
				"Date":     "2020-05-06",
				"Account":  s.initAccount.Name,
				"Category": s.initCategory.Title,
				"Currency": s.initCurrency.Abbreviation,
				"Amount":   "+0.01",
				"Note":     "Note1",
			},
		},
		{
			name: "PositiveTwoDigits",
			give: &model.Transaction{
				Date:     time.Date(2020, 5, 6, 11, 45, 04, 0, time.UTC),
				Account:  s.initAccount,
				Category: s.initCategory,
				Amount:   31,
				Note:     "Note1",
			},
			expected: map[string]string{
				"Date":     "2020-05-06",
				"Account":  s.initAccount.Name,
				"Category": s.initCategory.Title,
				"Currency": s.initCurrency.Abbreviation,
				"Amount":   "+0.31",
				"Note":     "Note1",
			},
		},
		{
			name: "NegativeLongNumber",
			give: &model.Transaction{
				Date:     time.Date(2020, 5, 6, 11, 45, 04, 0, time.UTC),
				Account:  s.initAccount,
				Category: s.initCategory,
				Amount:   -342342321231,
				Note:     "Note1",
			},
			expected: map[string]string{
				"Date":     "2020-05-06",
				"Account":  s.initAccount.Name,
				"Category": s.initCategory.Title,
				"Currency": s.initCurrency.Abbreviation,
				"Amount":   "-3423423212.31",
				"Note":     "Note1",
			},
		},
	} {
		s.Run(tc.name, func() {
			actual := s.presenter.ToMap(tc.give)
			assert.Equal(s.T(), tc.expected, actual)
		})
	}
}

func (s *TransactionPresenterTestSuite) TestFromMapPositive() {
	for _, tt := range []struct {
		name     string
		give     map[string]string
		expected *model.Transaction
	}{
		{
			name: "ExistingNestedModels",
			give: map[string]string{
				"Date":     "2020-05-06",
				"Account":  s.initAccount.Name,
				"Category": s.initCategory.Title,
				"Currency": s.initCurrency.Abbreviation,
				"Amount":   "0.00",
				"Note":     "Note1",
			},
			expected: &model.Transaction{
				Date:     time.Date(2020, 5, 6, 0, 0, 0, 0, time.UTC),
				Account:  s.initAccount,
				Category: s.initCategory,
				Amount:   0,
				Note:     "Note1",
			},
		},
		{
			name: "NotExistingNestedModels",
			give: map[string]string{
				"Date":     "2020-05-06",
				"Account":  "??",
				"Category": "???",
				"Currency": s.initCurrency.Abbreviation,
				"Amount":   "0.00",
				"Note":     "Note1",
			},
			expected: &model.Transaction{
				Date:     time.Date(2020, 5, 6, 0, 0, 0, 0, time.UTC),
				Account:  nil,
				Category: nil,
				Amount:   0,
				Note:     "Note1",
			},
		},
	} {
		s.Run(tt.name, func() {
			actual, err := s.presenter.FromMap(tt.give)
			require.NoError(s.T(), err)
			assert.Equal(s.T(), tt.expected, actual)
		})
	}
}



func (s *TransactionPresenterTestSuite) TestFromMapNegative() {
	for _, tt := range []struct {
		name     string
		give     map[string]string
		expected string
	}{
		{
			name: "InvalidDate",
			give: map[string]string{
				"Date":     "invalid",
				"Account":  s.initAccount.Name,
				"Category": s.initCategory.Title,
				"Currency": s.initCurrency.Abbreviation,
				"Amount":   "0.00",
				"Note":     "Note1",
			},
			expected: `time.Parse: parsing time "invalid" as "2006-01-02": cannot parse "invalid" as "2006"`,
		},
		{
			name: "InvalidAmount",
			give: map[string]string{
				"Date":     "2020-05-06",
				"Account":  s.initAccount.Name,
				"Category": s.initCategory.Title,
				"Currency": s.initCurrency.Abbreviation,
				"Amount":   "invalid",
				"Note":     "Note1",
			},
			expected: `parseMoney: strconv.Atoi: strconv.Atoi: parsing "invalid": invalid syntax`,
		},
		{
			name: "MissingDate",
			give: map[string]string{
				"Account":  s.initAccount.Name,
				"Category": s.initCategory.Title,
				"Currency": s.initCurrency.Abbreviation,
				"Amount":   "0.00",
				"Note":     "Note1",
			},
			expected: `checkKeys: key "Date" is missing`,
		},
		{
			name: "MissingAccount",
			give: map[string]string{
				"Date":     "2020-05-06",
				"Category": s.initCategory.Title,
				"Currency": s.initCurrency.Abbreviation,
				"Amount":   "0.00",
				"Note":     "Note1",
			},
			expected: `checkKeys: key "Account" is missing`,
		},
		{
			name: "MissingCategory",
			give: map[string]string{
				"Date":     "2020-05-06",
				"Account":  s.initAccount.Name,
				"Currency": s.initCurrency.Abbreviation,
				"Amount":   "0.00",
				"Note":     "Note1",
			},
			expected: `checkKeys: key "Category" is missing`,
		},
		{
			name: "MissingCurrency",
			give: map[string]string{
				"Date":     "2020-05-06",
				"Account":  s.initAccount.Name,
				"Category": s.initCategory.Title,
				"Amount":   "0.00",
				"Note":     "Note1",
			},
			expected: `checkKeys: key "Currency" is missing`,
		},
		{
			name: "MissingAmount",
			give: map[string]string{
				"Date":     "2020-05-06",
				"Account":  s.initAccount.Name,
				"Category": s.initCategory.Title,
				"Currency": s.initCurrency.Abbreviation,
				"Note":     "Note1",
			},
			expected: `checkKeys: key "Amount" is missing`,
		},
		{
			name: "MissingNote",
			give: map[string]string{
				"Date":     "2020-05-06",
				"Account":  s.initAccount.Name,
				"Category": s.initCategory.Title,
				"Currency": s.initCurrency.Abbreviation,
				"Amount":   "0.00",
			},
			expected: `checkKeys: key "Note" is missing`,
		},

	} {
		s.Run(tt.name, func() {
			_, err := s.presenter.FromMap(tt.give)
			assert.EqualError(s.T(), err, tt.expected)
		})
	}
}

func (s *TransactionPresenterTestSuite) TearDownSuite() {
	err := s.db.Close()
	require.NoError(s.T(), err, "occurred in TearDownSuite")
}

func TestTransactionPresenterTestSuite(t *testing.T) {
	suite.Run(t, new(TransactionPresenterTestSuite))
}
