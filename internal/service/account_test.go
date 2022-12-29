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

type accountList []*model.Account

func (cc accountList) Len() int           { return len(cc) }
func (cc accountList) Less(i, j int) bool { return cc[i].Name < cc[j].Name }
func (cc accountList) Swap(i, j int)      { cc[i], cc[j] = cc[j], cc[i] }

type AccountServiceTestSuite struct {
	suite.Suite
	db                *sql.DB
	persistentStorage *sqlite.SqliteStorage
	inmemoryStorage   *inmemory.InmemoryStorage
	service           *service.Service
	InitCurrencies    []*model.Currency
	InitAccounts      []*model.Account
}

func (s *AccountServiceTestSuite) SetupSuite() {
	db, err := sql.Open("sqlite3", "file::memory:?cache=shared")
	require.NoError(s.T(), err, "occurred in SetupSuite")
	s.db = db

	s.persistentStorage, err = sqlite.New(db)
	require.NoError(s.T(), err, "occurred in SetupSuite")
	s.inmemoryStorage = inmemory.New()
	s.service, err = service.New(s.persistentStorage, s.inmemoryStorage)
	require.NoError(s.T(), err, "occurred in SetupSuite")

	// id's settled by sqlite on insert incrementally starting from 1,
	// so here they are initialized for match purpose
	s.InitCurrencies = []*model.Currency{
		{ID: 1, Abbreviation: "USD", ExchangeRate: 100, IsMain: true},
		{ID: 2, Abbreviation: "EUR", ExchangeRate: 124, IsMain: false},
	}
	s.InitAccounts = []*model.Account{
		{ID: 1, Name: "BCard1", Currency: s.InitCurrencies[0]},
		{ID: 2, Name: "ACard2", Currency: s.InitCurrencies[1]},
	}

	// init currencies
	for _, c := range s.InitCurrencies {
		_, err = s.persistentStorage.Currency().Insert(c)
		require.NoError(s.T(), err, "occurred in SetupTest")
	}

	err = s.service.Currency().Init()
	require.NoError(s.T(), err, "occurred in SetupTest")
}

func (s *AccountServiceTestSuite) SetupTest() {
	for _, a := range s.InitAccounts {
		_, err := s.persistentStorage.Account().Insert(a)
		require.NoError(s.T(), err, "occurred in SetupTest")
	}

	err := s.service.Account().Init(s.service.Currency())
	require.NoError(s.T(), err, "occurred in SetupTest")
}

func (s *AccountServiceTestSuite) TestLinkage() {
	aa := s.service.Account().GetAll()

	assert.EqualValues(s.T(), s.InitCurrencies[0], aa[0].Currency)
	assert.EqualValues(s.T(), s.InitCurrencies[1], aa[1].Currency)
}

func (s *AccountServiceTestSuite) TestInsertPositive() {
	account := &model.Account{ID: 3, Name: "Card3", Currency: s.InitCurrencies[1]}
	expectedAccounts := append(s.InitAccounts, account)

	err := s.service.Account().Insert(account)
	require.NoError(s.T(), err)

	persistentAccounts := s.getLinkedPersistantAccounts()
	inmemoryAccounts := s.inmemoryStorage.Account().GetAll()
	assert.ElementsMatch(s.T(), persistentAccounts, expectedAccounts)
	assert.ElementsMatch(s.T(), inmemoryAccounts, expectedAccounts)
}

func (s *AccountServiceTestSuite) TestInsertNegative() {
	err := s.service.Account().Insert(s.InitAccounts[0])
	assert.ErrorContains(s.T(), err, "s.persistentStorage.Insert: e.db.Exec: UNIQUE constraint failed: account.name")
}

func (s *AccountServiceTestSuite) TestUpdatePositive() {
	aa := s.service.Account().GetAll()
	aa[0].Name = "Card10"

	err := s.service.Account().Update(aa[0])
	require.NoError(s.T(), err)

	persistentAccounts := s.getLinkedPersistantAccounts()
	inmemoryAccounts := s.inmemoryStorage.Account().GetAll()
	require.NoError(s.T(), err)
	assert.ElementsMatch(s.T(), persistentAccounts, inmemoryAccounts)
}

func (s *AccountServiceTestSuite) TestUpdateNegative() {
	aa := s.service.Account().GetAll()
	aa[0].Name = "Card10"
	aa[0].ID = 10

	err := s.service.Account().Update(aa[0])
	assert.ErrorContains(s.T(), err, "s.persistentStorage.Update: total affected rows 0 while expected 1")
	aa[0].ID = 1 // return real id to proper teardown
}

func (s *AccountServiceTestSuite) TestDeletePositive() {
	aa := s.service.Account().GetAll()
	expectedAccounts := []*model.Account{aa[0]}

	err := s.service.Account().Delete(aa[1])
	require.NoError(s.T(), err)

	persistentAccounts := s.getLinkedPersistantAccounts()
	inmemoryAccounts := s.inmemoryStorage.Account().GetAll()
	assert.ElementsMatch(s.T(), persistentAccounts, expectedAccounts)
	assert.ElementsMatch(s.T(), inmemoryAccounts, expectedAccounts)
}

func (s *AccountServiceTestSuite) TestDeleteNegative() {
	aa := s.service.Account().GetAll()
	aa[0].ID = 10

	err := s.service.Account().Delete(aa[0])
	assert.ErrorContains(s.T(), err, "s.persistentStorage.Delete: total affected rows 0 while expected 1")
	aa[0].ID = 1 // return real id to proper teardown
}

func (s *AccountServiceTestSuite) TestGetByID() {
	a := s.service.Account().GetByID(2)
	assert.Equal(s.T(), s.InitAccounts[1].Name, a.Name)
}

func (s *AccountServiceTestSuite) TestGetByName() {
	a := s.service.Account().GetByName(s.InitAccounts[1].Name)
	assert.Equal(s.T(), s.InitCurrencies[1].ID, a.ID)
}

func (s *AccountServiceTestSuite) TestGetAllSorted() {
	cc := s.service.Account().GetAllSorted()
	expectedAccounts := make([]*model.Account, len(s.InitAccounts))
	copy(expectedAccounts, s.InitAccounts)
	sort.Sort(accountList(expectedAccounts))
	assert.ElementsMatch(s.T(), cc, expectedAccounts)
}

func (s *AccountServiceTestSuite) getLinkedPersistantAccounts() []*model.Account {
	persistentAccounts, err := s.persistentStorage.Account().GetAll()
	require.NoError(s.T(), err)

	for _, a := range persistentAccounts {
		a.Currency = s.service.Currency().GetByID(a.Currency.ID)
	}

	return persistentAccounts
}

func (s *AccountServiceTestSuite) TearDownTest() {
	for {
		aa := s.service.Account().GetAll()
		if len(aa) == 0 {
			break
		}

		err := s.persistentStorage.Account().Delete(aa[0].ID)
		require.NoError(s.T(), err, "occurred in TearDownTest")
		s.inmemoryStorage.Account().Delete(aa[0])
	}
}

func (s *AccountServiceTestSuite) TearDownSuite() {
	err := s.db.Close()
	require.NoError(s.T(), err, "occurred in TearDownSuite")
}

func TestAccountServiceTestSuite(t *testing.T) {
	suite.Run(t, new(AccountServiceTestSuite))
}
