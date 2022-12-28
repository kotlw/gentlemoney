package service_test

import (
	"database/sql"
	"time"
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

type transactionList []*model.Transaction

func (tt transactionList) Len() int           { return len(tt) }
func (tt transactionList) Less(i, j int) bool { return tt[i].Date.After(tt[j].Date) }
func (tt transactionList) Swap(i, j int)      { tt[i], tt[j] = tt[j], tt[i] }

type TransactionServiceStorageTestSuite struct {
	suite.Suite
	db                *sql.DB
	persistantStorage *sqlite.SqliteStorage
	inmemoryStorage   *inmemory.InmemoryStorage
	service           *service.Service
	InitCategories    []*model.Category
	InitCurrencies    []*model.Currency
	InitAccounts      []*model.Account
	InitTransactions  []*model.Transaction
}

func (s *TransactionServiceStorageTestSuite) SetupSuite() {
	db, err := sql.Open("sqlite3", "file::memory:?cache=shared")
	require.NoError(s.T(), err, "occured in SetupSuite")
	s.db = db

	s.persistantStorage, err = sqlite.New(db)
	require.NoError(s.T(), err, "occured in SetupSuite")
	s.inmemoryStorage = inmemory.New()
	s.service, err = service.New(s.persistantStorage, s.inmemoryStorage)
	require.NoError(s.T(), err, "occured in SetupSuite")

	// id's settled by sqlite on insert incrementally starting from 1,
	// so here they are initialized for match purpose
	s.InitCategories = []*model.Category{
		{ID: 1, Title: "Health"},
		{ID: 2, Title: "Grocery"},
	}
	s.InitCurrencies = []*model.Currency{
		{ID: 1, Abbreviation: "USD", ExchangeRate: 100, IsMain: true},
		{ID: 2, Abbreviation: "EUR", ExchangeRate: 124, IsMain: false},
	}
	s.InitAccounts = []*model.Account{
		{ID: 1, Name: "BCard1", Currency: s.InitCurrencies[0]},
		{ID: 2, Name: "ACard2", Currency: s.InitCurrencies[1]},
	}
	s.InitTransactions = []*model.Transaction{
		{
			ID:       1,
			Date:     time.Date(2022, time.Month(2), 21, 1, 10, 30, 0, time.UTC),
			Account:  s.InitAccounts[0],
			Category: s.InitCategories[0],
			Amount:   12345,
			Note:     "note1",
		},
		{
			ID:       2,
			Date:     time.Date(2022, time.Month(2), 22, 1, 10, 30, 0, time.UTC),
			Account:  s.InitAccounts[1],
			Category: s.InitCategories[1],
			Amount:   67890,
			Note:     "note2",
		},
	}

	// init categories
	for _, c := range s.InitCategories {
		_, err = s.persistantStorage.Category().Insert(c)
		require.NoError(s.T(), err, "occured in SetupTest")
	}

	err = s.service.Category().Init()
	require.NoError(s.T(), err, "occured in SetupTest")

	// init currencies
	for _, c := range s.InitCurrencies {
		_, err = s.persistantStorage.Currency().Insert(c)
		require.NoError(s.T(), err, "occured in SetupTest")
	}

	err = s.service.Currency().Init()
	require.NoError(s.T(), err, "occured in SetupTest")

	// init accounts
	for _, a := range s.InitAccounts {
		_, err = s.persistantStorage.Account().Insert(a)
		require.NoError(s.T(), err, "occured in SetupTest")
	}

	err = s.service.Account().Init(s.service.Currency())
	require.NoError(s.T(), err, "occured in SetupTest")

}

func (s *TransactionServiceStorageTestSuite) SetupTest() {
	// init transactions
	for _, t := range s.InitTransactions {
		_, err := s.persistantStorage.Transaction().Insert(t)
		require.NoError(s.T(), err, "occured in SetupTest")
	}

	err := s.service.Transaction().Init(s.service.Category(), s.service.Account())
	require.NoError(s.T(), err, "occured in SetupTest")
}

func (s *TransactionServiceStorageTestSuite) TestLinkage() {
	tt := s.service.Transaction().GetAll()

	assert.EqualValues(s.T(), s.InitCategories[0], tt[0].Category)
	assert.EqualValues(s.T(), s.InitCategories[1], tt[1].Category)
	assert.EqualValues(s.T(), s.InitAccounts[0], tt[0].Account)
	assert.EqualValues(s.T(), s.InitAccounts[1], tt[1].Account)
}

func (s *TransactionServiceStorageTestSuite) TestInsertPositive() {
	transaction := &model.Transaction{
		ID:       3,
		Date:     time.Date(2022, time.Month(2), 23, 1, 10, 30, 0, time.UTC),
		Account:  s.InitAccounts[0],
		Category: s.InitCategories[0],
		Amount:   4321,
		Note:     "note3",
	}
	expectedTransactions := append(s.InitTransactions, transaction)

	err := s.service.Transaction().Insert(transaction)
	require.NoError(s.T(), err)

	persistantTransactions := s.getLinkedPersistantTransactions()
	inmemoryTransactions := s.inmemoryStorage.Transaction().GetAll()
	assert.ElementsMatch(s.T(), persistantTransactions, expectedTransactions)
	assert.ElementsMatch(s.T(), inmemoryTransactions, expectedTransactions)
}

func (s *TransactionServiceStorageTestSuite) TestUpdatePositive() {
	tt := s.service.Transaction().GetAll()
	tt[0].Amount = 99102 
	tt[0].Note = "CHANGED"

	err := s.service.Transaction().Update(tt[0])
	require.NoError(s.T(), err)

	persistantTransactions := s.getLinkedPersistantTransactions()
	inmemoryTransactions := s.inmemoryStorage.Transaction().GetAll()
	require.NoError(s.T(), err)
	assert.ElementsMatch(s.T(), persistantTransactions, inmemoryTransactions)
}

func (s *TransactionServiceStorageTestSuite) TestUpdateNegative() {
	tt := s.service.Transaction().GetAll()
	tt[0].ID = 10
	tt[0].Note = "CHANGED"

	err := s.service.Transaction().Update(tt[0])
	assert.ErrorContains(s.T(), err, "s.persistantStorage.Update: total affected rows 0 while expected 1")
	tt[0].ID = 1 // return real id to proper teardown
}

func (s *TransactionServiceStorageTestSuite) TestDeletePositive() {
	tt := s.service.Transaction().GetAll()
	expectedTransactions := []*model.Transaction{tt[0]}

	err := s.service.Transaction().Delete(tt[1])
	require.NoError(s.T(), err)

	persistantTransactions := s.getLinkedPersistantTransactions()
	inmemoryTransactions := s.inmemoryStorage.Transaction().GetAll()
	assert.ElementsMatch(s.T(), persistantTransactions, expectedTransactions)
	assert.ElementsMatch(s.T(), inmemoryTransactions, expectedTransactions)
}

func (s *TransactionServiceStorageTestSuite) TestDeleteNegative() {
	tt := s.service.Transaction().GetAll()
	tt[0].ID = 10

	err := s.service.Transaction().Delete(tt[0])
	assert.ErrorContains(s.T(), err, "s.persistantStorage.Delete: total affected rows 0 while expected 1")
	tt[0].ID = 1 // return real id to proper teardown
}

func (s *TransactionServiceStorageTestSuite) TestGetByID() {
	t := s.service.Transaction().GetByID(2)
	assert.EqualValues(s.T(), s.InitTransactions[1], t)
}

func (s *TransactionServiceStorageTestSuite) TestGetAllSorted() {
	tt := s.service.Transaction().GetAllSorted()
	expectedTransactions := make([]*model.Transaction, len(s.InitTransactions))
	copy(expectedTransactions, s.InitTransactions)
	sort.Sort(transactionList(expectedTransactions))
	assert.ElementsMatch(s.T(), tt, expectedTransactions)
}

func (s *TransactionServiceStorageTestSuite) getLinkedPersistantTransactions() []*model.Transaction {
	persistantTransactions, err := s.persistantStorage.Transaction().GetAll()
	require.NoError(s.T(), err)

	for _, t := range persistantTransactions {
		t.Category = s.service.Category().GetByID(t.Category.ID)
		t.Account = s.service.Account().GetByID(t.Account.ID)
	}

	return persistantTransactions
}

func (s *TransactionServiceStorageTestSuite) TearDownTest() {
	for {
		tt := s.service.Transaction().GetAll()
		if len(tt) == 0 {
			break
		}

		err := s.persistantStorage.Transaction().Delete(tt[0].ID)
		require.NoError(s.T(), err, "occured in TearDownTest")
		s.inmemoryStorage.Transaction().Delete(tt[0])
	}
}

func (s *TransactionServiceStorageTestSuite) TearDownSuite() {
	err := s.db.Close()
	require.NoError(s.T(), err, "occured in TearDownSuite")
}

func TestTransactionServiceStorageTestSuite(t *testing.T) {
	suite.Run(t, new(TransactionServiceStorageTestSuite))
}
